package vmware

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// In vsphere we get vlans from DistributedVirtualPortgroups
func (vc *VmwareSource) InitNetworks(ctx context.Context, containerView *view.ContainerView) error {
	var dvpgs []mo.DistributedVirtualPortgroup
	err := containerView.Retrieve(ctx, []string{"DistributedVirtualPortgroup"}, []string{"config"}, &dvpgs)
	if err != nil {
		return fmt.Errorf("failed retrieving DistributedVirtualPortgroups: %s", err)
	}
	vc.Networks = NetworkData{
		DistributedVirtualPortgroups: make(map[string]*DistributedPortgroupData),
		HostVirtualSwitches:          make(map[string]map[string]*HostVirtualSwitchData),
		HostProxySwitches:            make(map[string]map[string]*HostProxySwitchData),
		HostPortgroups:               make(map[string]map[string]*HostPortgroupData),
	}
	for _, dvpg := range dvpgs {
		if dvpg.Config.Key == "" || dvpg.Config.Name == "" {
			continue
		}

		vlanInfo := dvpg.Config.DefaultPortConfig.(*types.VMwareDVSPortSetting)
		var vlanIds []int
		var vlanIdRanges []string
		private := false

		switch v := vlanInfo.Vlan.(type) {
		case *types.VmwareDistributedVirtualSwitchTrunkVlanSpec:
			for _, item := range v.VlanId {
				if item.Start == item.End {
					vlanIds = append(vlanIds, int(item.Start))
					vlanIdRanges = append(vlanIdRanges, fmt.Sprintf("%d", item.Start))
				} else if item.Start == 0 && item.End == 4094 {
					vlanIds = append(vlanIds, 4095)
					vlanIdRanges = append(vlanIdRanges, fmt.Sprintf("%d-%d", item.Start, item.End))
				} else {
					for vlan := item.Start; vlan <= item.End; vlan++ {
						vlanIds = append(vlanIds, int(vlan))
					}
					vlanIdRanges = append(vlanIdRanges, fmt.Sprintf("%d-%d", item.Start, item.End))
				}
			}
		case *types.VmwareDistributedVirtualSwitchPvlanSpec:
			vlanIds = append(vlanIds, int(v.PvlanId))
			private = true
		case *types.VmwareDistributedVirtualSwitchVlanIdSpec:
			vlanIds = append(vlanIds, int(v.VlanId))
		default:
			return fmt.Errorf("uknown vlan info spec %T", v)
		}

		vc.Networks.DistributedVirtualPortgroups[dvpg.Config.Key] = &DistributedPortgroupData{
			Name:         dvpg.Config.Name,
			VlanIds:      vlanIds,
			VlanIdRanges: vlanIdRanges,
			Private:      private,
		}

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

		// Add network data which is received from hosts
		// Iterate over hosts virtual switches, needed to enrich data on physical interfaces
		for _, vswitch := range host.Config.Network.Vswitch {
			if vswitch.Name != "" {
				vc.Networks.HostVirtualSwitches[host.Name][vswitch.Name] = &HostVirtualSwitchData{
					mtu:   int(vswitch.Mtu),
					pnics: vswitch.Pnic,
				}
			}
		}
		// Iterate over hosts proxy switches, needed to enrich data on physical interfaces
		// Also stores mtu data which is used for VM interfaces
		hostProxyswitchData := make(map[string]map[string]*HostProxySwitchData)
		for _, pswitch := range host.Config.Network.ProxySwitch {
			if pswitch.DvsUuid != "" {
				hostProxyswitchData[host.Name][pswitch.DvsUuid] = &HostProxySwitchData{
					mtu:   int(pswitch.Mtu),
					pnics: pswitch.Pnic,
					name:  pswitch.DvsName,
				}
			}
		}
		// Iterate over hosts port groups, needed to enrich data on physical interfaces
		for _, pgroup := range host.Config.Network.Portgroup {
			if pgroup.Spec.Name != "" {
				nic_order := pgroup.ComputedPolicy.NicTeaming.NicOrder
				pgroup_nics := []string{}
				if len(nic_order.ActiveNic) > 0 {
					pgroup_nics = append(pgroup_nics, nic_order.ActiveNic...)
				}
				if len(nic_order.StandbyNic) > 0 {
					pgroup_nics = append(pgroup_nics, nic_order.StandbyNic...)
				}
				vc.Networks.HostPortgroups[host.Name][pgroup.Spec.Name] = &HostPortgroupData{
					vlanId:  int(pgroup.Spec.VlanId),
					vswitch: pgroup.Spec.VswitchName,
					nics:    pgroup_nics,
				}
			}
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
