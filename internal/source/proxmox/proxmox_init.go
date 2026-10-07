package proxmox

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"

	"github.com/luthermonson/go-proxmox"
)

func (ps *ProxmoxSource) initCluster(ctx context.Context, c *proxmox.Client) error {
	cluster, err := c.Cluster(ctx)
	if err != nil {
		return fmt.Errorf("init cluster: %s", err)
	}
	ps.Cluster = cluster

	return nil
}

func (ps *ProxmoxSource) initNodes(ctx context.Context, c *proxmox.Client) error {
	nodes, err := c.Nodes(ctx)
	if err != nil {
		return fmt.Errorf("init nodes: %s", err)
	}

	// Nodes are fetched in parallel, each into its own nodeData, then merged in order.
	nodesData := make([]*nodeData, len(nodes))
	nodesErrs := make([]error, len(nodes))
	var wg sync.WaitGroup
	for i, nodeStatus := range nodes {
		wg.Add(1)
		go func(i int, nodeName string) {
			defer wg.Done()
			nodesData[i], nodesErrs[i] = ps.fetchNode(ctx, c, nodeName)
		}(i, nodeStatus.Node)
	}
	wg.Wait()
	if err := errors.Join(nodesErrs...); err != nil {
		return err
	}

	ps.Nodes = make([]*proxmox.Node, 0, len(nodes))
	ps.NodeIfaces = make(map[string][]*proxmox.NodeNetwork, len(nodes))
	ps.Vms = make(map[string][]*proxmox.VirtualMachine, len(nodes))
	ps.VMIfaces = make(map[uint64][]*proxmox.AgentNetworkIface, 0)
	ps.Containers = make(map[string][]*proxmox.Container, len(nodes))
	ps.ContainerIfaces = make(map[uint64][]*proxmox.ContainerInterface, 0)
	for _, data := range nodesData {
		ps.Nodes = append(ps.Nodes, data.node)
		ps.NodeIfaces[data.node.Name] = data.networks
		ps.Vms[data.node.Name] = data.vms
		maps.Copy(ps.VMIfaces, data.vmIfaces)
		ps.Containers[data.node.Name] = data.containers
		maps.Copy(ps.ContainerIfaces, data.containerIfaces)
	}
	return nil
}

// nodeData holds what is collected from the Proxmox API for one node.
type nodeData struct {
	node            *proxmox.Node
	networks        []*proxmox.NodeNetwork
	vms             []*proxmox.VirtualMachine
	vmIfaces        map[uint64][]*proxmox.AgentNetworkIface
	containers      []*proxmox.Container
	containerIfaces map[uint64][]*proxmox.ContainerInterface
}

// fetchNode collects the networks, VMs and containers of the node named nodeName.
func (ps *ProxmoxSource) fetchNode(ctx context.Context, c *proxmox.Client, nodeName string) (*nodeData, error) {
	node, err := c.Node(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("init node: %s", err)
	}
	data := &nodeData{node: node}

	data.networks, err = ps.initNodeNetworks(ctx, node)
	if err != nil {
		return nil, fmt.Errorf("init nodeNetworks: %s", err)
	}

	data.vms, data.vmIfaces, err = ps.initNodeVMs(ctx, node)
	if err != nil {
		return nil, fmt.Errorf("init nodeVMs: %s", err)
	}

	data.containers, data.containerIfaces, err = ps.initContainers(ctx, node)
	if err != nil {
		return nil, fmt.Errorf("init node containers: %s", err)
	}
	return data, nil
}

// Helper function for initNodes. It collects all nodeNetwork for given node.
func (ps *ProxmoxSource) initNodeNetworks(ctx context.Context, node *proxmox.Node) ([]*proxmox.NodeNetwork, error) {
	nodeNetworks, err := node.Networks(ctx)
	if err != nil {
		return nil, fmt.Errorf("init nodeNetworks: %s", err)
	}
	nodeIfaces := make([]*proxmox.NodeNetwork, 0, len(nodeNetworks))
	for _, nodeNetwork := range nodeNetworks {
		nodeIface, err := node.Network(ctx, nodeNetwork.Iface)
		if err != nil {
			return nil, fmt.Errorf("init nodeIface: %s", err)
		}
		nodeIfaces = append(nodeIfaces, nodeIface)
	}
	return nodeIfaces, nil
}

// Helper function for initNodes. It collects all vms for given node.
func (ps *ProxmoxSource) initNodeVMs(
	ctx context.Context,
	node *proxmox.Node,
) ([]*proxmox.VirtualMachine, map[uint64][]*proxmox.AgentNetworkIface, error) {
	// Fetch VMs list
	vms, err := node.VirtualMachines(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Fetch VM config for each VM
	nodeVMs := make([]*proxmox.VirtualMachine, 0, len(vms))
	vmIfaces := make(map[uint64][]*proxmox.AgentNetworkIface, len(vms))
	for _, vm := range vms {
		vmconfig, err := node.VirtualMachine(ctx, int(vm.VMID)) //nolint:gosec // VMID fits in int
		if err != nil {
			return nil, nil, fmt.Errorf("init nodeVms: %s", err)
		}

		// Load VM disks
		vmconfig.VirtualMachineConfig.MergeDisks()

		// Store VM info in our list
		nodeVMs = append(nodeVMs, vmconfig)

		// Only a running VM has a guest agent to ask: the interfaces of other VMs stay unknown.
		if vm.Status != "running" {
			continue
		}

		// Load VM interfaces. When the guest agent does not answer, the VM is left out
		// of VMIfaces: its interfaces are unknown, which is not the same as having none.
		ifaces, err := vm.AgentGetNetworkIFaces(ctx)
		if err != nil {
			ps.Logger.Debugf(ps.Ctx, "vm %s: guest agent network data unavailable: %s", vm.Name, err)
			continue
		}
		vmIfaces[uint64(vm.VMID)] = ifaces
	}
	return nodeVMs, vmIfaces, nil
}

// Helper function for initNodes. It collects all containers for given node.
func (ps *ProxmoxSource) initContainers(
	ctx context.Context,
	node *proxmox.Node,
) ([]*proxmox.Container, map[uint64][]*proxmox.ContainerInterface, error) {
	containers, err := node.Containers(ctx)
	if err != nil {
		return nil, nil, err
	}

	nodeContainers := make([]*proxmox.Container, 0, len(containers))
	containerIfaces := make(map[uint64][]*proxmox.ContainerInterface, len(containers))
	for _, container := range containers {
		nodeContainers = append(nodeContainers, container)
		// When the interfaces cannot be read, the container is left out of ContainerIfaces.
		ifaces, err := container.Interfaces(ctx)
		if err != nil {
			ps.Logger.Debugf(ps.Ctx, "container %s: network data unavailable: %s", container.Name, err)
			continue
		}
		containerIfaces[uint64(container.VMID)] = ifaces
	}
	return nodeContainers, containerIfaces, nil
}
