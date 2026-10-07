package proxmox

import (
	"context"
	"fmt"

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

	ps.Nodes = make([]*proxmox.Node, 0, len(nodes))
	ps.NodeIfaces = make(map[string][]*proxmox.NodeNetwork, len(nodes))
	ps.Vms = make(map[string][]*proxmox.VirtualMachine, len(nodes))
	ps.VMIfaces = make(map[uint64][]*proxmox.AgentNetworkIface, 0)
	ps.Containers = make(map[string][]*proxmox.Container, len(nodes))
	ps.ContainerIfaces = make(map[uint64][]*proxmox.ContainerInterface, 0)

	for _, node := range nodes {
		node, err := c.Node(ctx, node.Node)
		if err != nil {
			return fmt.Errorf("init node: %s", err)
		}
		ps.Nodes = append(ps.Nodes, node)

		err = ps.initNodeNetworks(ctx, node)
		if err != nil {
			return fmt.Errorf("init nodeNetworks: %s", err)
		}

		err = ps.initNodeVMs(ctx, node)
		if err != nil {
			return fmt.Errorf("init nodeVMs: %s", err)
		}

		err = ps.initContainers(ctx, node)
		if err != nil {
			return fmt.Errorf("init node containers: %s", err)
		}
	}
	return nil
}

// Helper function for initNodes. It collects all nodeNetwork for given node.
func (ps *ProxmoxSource) initNodeNetworks(ctx context.Context, node *proxmox.Node) error {
	nodeNetworks, err := node.Networks(ctx)
	if err != nil {
		return fmt.Errorf("init nodeNetworks: %s", err)
	}
	ps.NodeIfaces[node.Name] = make([]*proxmox.NodeNetwork, 0, len(nodeNetworks))
	for _, nodeNetwork := range nodeNetworks {
		nodeIface, err := node.Network(ctx, nodeNetwork.Iface)
		if err != nil {
			return fmt.Errorf("init nodeIface: %s", err)
		}
		ps.NodeIfaces[node.Name] = append(ps.NodeIfaces[node.Name], nodeIface)
	}
	return nil
}

// Helper function for initNodes. It collects all vms for given node.
func (ps *ProxmoxSource) initNodeVMs(ctx context.Context, node *proxmox.Node) error {
	// Fetch VMs list
	vms, err := node.VirtualMachines(ctx)
	if err != nil {
		return err
	}

	// Fetch VM config for each VM
	ps.Vms[node.Name] = make([]*proxmox.VirtualMachine, 0, len(vms))
	for _, vm := range vms {
		vmconfig, err := node.VirtualMachine(ctx, int(vm.VMID)) //nolint:gosec // VMID fits in int
		if err != nil {
			return fmt.Errorf("init nodeVms: %s", err)
		}

		// Load VM disks
		vmconfig.VirtualMachineConfig.MergeDisks()

		// Store VM info in our list
		ps.Vms[node.Name] = append(ps.Vms[node.Name], vmconfig)

		// Load VM interfaces. When the guest agent does not answer, the VM is left out
		// of VMIfaces: its interfaces are unknown, which is not the same as having none.
		ifaces, err := vm.AgentGetNetworkIFaces(ctx)
		if err != nil {
			ps.Logger.Debugf(ps.Ctx, "vm %s: guest agent network data unavailable: %s", vm.Name, err)
			continue
		}
		ps.VMIfaces[uint64(vm.VMID)] = ifaces
	}
	return nil
}

// Helper function for initNodes. It collects all containers for given node.
func (ps *ProxmoxSource) initContainers(ctx context.Context, node *proxmox.Node) error {
	containers, err := node.Containers(ctx)
	if err != nil {
		return err
	}

	ps.Containers[node.Name] = make([]*proxmox.Container, 0, len(containers))
	for _, container := range containers {
		ps.Containers[node.Name] = append(ps.Containers[node.Name], container)
		// When the interfaces cannot be read, the container is left out of ContainerIfaces.
		ifaces, err := container.Interfaces(ctx)
		if err != nil {
			ps.Logger.Debugf(ps.Ctx, "container %s: network data unavailable: %s", container.Name, err)
			continue
		}
		ps.ContainerIfaces[uint64(container.VMID)] = ifaces
	}
	return nil
}
