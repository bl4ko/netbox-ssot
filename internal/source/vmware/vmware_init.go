package vmware

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

// In vsphere we get vlans from DistributedVirtualPortgroups
func (vc *VmwareSource) InitVlans(ctx context.Context, containerView *view.ContainerView) error {
	var dvpgs []mo.DistributedVirtualPortgroup
	err := containerView.Retrieve(ctx, []string{"DistributedVirtualPortgroup"}, []string{"config"}, &dvpgs)
	if err != nil {
		return fmt.Errorf("failed retrieving DistributedVirtualPortgroups: %s", err)
	}
	for _, dvpg := range dvpgs {
		vc.DistributedVirtualPortgrups[dvpg.Self.Value] = &dvpg
	}
	return nil
}

func (vc *VmwareSource) InitDisks(ctx context.Context, containerView *view.ContainerView) error {
	var disks []mo.Datastore
	err := containerView.Retrieve(ctx, []string{"Datastore"}, []string{"summary", "host", "vm"}, &disks)
	if err != nil {
		return fmt.Errorf("failed retrieving disks: %s", err)
	}
	vc.Disks = make(map[string]*mo.Datastore, len(disks))
	for _, disk := range disks {
		vc.Disks[disk.Self.Value] = &disk
	}
	return nil
}

func (vc *VmwareSource) InitDataCenters(ctx context.Context, containerView *view.ContainerView) error {
	var datacenters []mo.Datacenter
	err := containerView.Retrieve(ctx, []string{"Datacenter"}, []string{"name"}, &datacenters)
	if err != nil {
		return fmt.Errorf("failed retrieving datacenters: %s", err)
	}
	vc.DataCenters = make(map[string]*mo.Datacenter, len(datacenters))
	for _, datacenter := range datacenters {
		vc.DataCenters[datacenter.Self.Value] = &datacenter
	}
	return nil
}

func (vc *VmwareSource) InitClusters(ctx context.Context, containerView *view.ContainerView) error {
	var clusters []mo.ClusterComputeResource
	err := containerView.Retrieve(ctx, []string{"ClusterComputeResource"}, []string{"summary", "host", "name"}, &clusters)
	if err != nil {
		return fmt.Errorf("failed retrieving clusters: %s", err)
	}
	vc.Host2Cluster = make(map[string]string)
	vc.Clusters = make(map[string]*mo.ClusterComputeResource, len(clusters))
	for _, cluster := range clusters {
		vc.Clusters[cluster.Self.Value] = &cluster
		for _, host := range cluster.Host {
			vc.Host2Cluster[host.Value] = cluster.Self.Value
		}
	}
	return nil
}

func (vc *VmwareSource) InitHosts(ctx context.Context, containerView *view.ContainerView) error {
	var hosts []mo.HostSystem
	err := containerView.Retrieve(ctx, []string{"HostSystem"}, []string{"name", "summary.host", "summary.hardware", "summary.runtime", "vm", "config.network"}, &hosts)
	if err != nil {
		return fmt.Errorf("failed retrieving hosts: %s", err)
	}
	vc.Vm2Host = make(map[string]string)
	vc.Hosts = make(map[string]*mo.HostSystem, len(hosts))
	for _, host := range hosts {
		vc.Hosts[host.Self.Value] = &host
		for _, vm := range host.Vm {
			vc.Vm2Host[vm.Value] = host.Self.Value
		}
	}
	return nil
}

func (vc *VmwareSource) InitVms(ctx context.Context, containerView *view.ContainerView) error {
	var vms []mo.VirtualMachine
	err := containerView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary", "name", "guest.net"}, &vms)
	if err != nil {
		return fmt.Errorf("failed retrieving vms: %s", err)
	}
	vc.Vms = make(map[string]*mo.VirtualMachine, len(vms))
	for _, vm := range vms {
		vc.Vms[vm.Self.Value] = &vm
	}
	return nil
}
