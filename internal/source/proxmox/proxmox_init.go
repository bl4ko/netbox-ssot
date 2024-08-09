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
	ps.VMIfaces = make(map[string][]*proxmox.AgentNetworkIface, 0)
	ps.Containers = make(map[string][]*proxmox.Container, len(nodes))
	ps.ContainerIfaces = make(map[string][]*proxmox.ContainerInterface, 0)

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
	vms, err := node.VirtualMachines(ctx)
	if err != nil {
		return err
	}
	ps.Vms[node.Name] = make([]*proxmox.VirtualMachine, 0, len(vms))
	for _, vm := range vms {
		ps.Vms[node.Name] = append(ps.Vms[node.Name], vm)
		ifaces, _ := vm.AgentGetNetworkIFaces(ctx)
		ps.VMIfaces[vm.Name] = make([]*proxmox.AgentNetworkIface, 0, len(ifaces))
		ps.VMIfaces[vm.Name] = append(ps.VMIfaces[vm.Name], ifaces...)
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
		ifaces, _ := container.Interfaces(ctx)
		ps.ContainerIfaces[container.Name] = make([]*proxmox.ContainerInterface, 0, len(ifaces))
		ps.ContainerIfaces[container.Name] = append(ps.ContainerIfaces[container.Name], ifaces...)
	}
	return nil
}
