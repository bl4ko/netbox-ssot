package vmware

import (
	"context"
	"fmt"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// In vsphere we get vlans from DistributedVirtualPortgroups.
func (vc *VmwareSource) InitNetworks(ctx context.Context, containerView *view.ContainerView) error {
	var dvpgs []mo.DistributedVirtualPortgroup
	err := containerView.Retrieve(ctx, []string{"DistributedVirtualPortgroup"}, []string{"config"}, &dvpgs)
	if err != nil {
		return fmt.Errorf("failed retrieving DistributedVirtualPortgroups: %s", err)
	}
	vc.Networks = NetworkData{
		DistributedVirtualPortgroups: make(map[string]*DistributedPortgroupData),
		Vid2Name:                     make(map[int]string),
		HostVirtualSwitches:          make(map[string]map[string]*HostVirtualSwitchData),
		HostProxySwitches:            make(map[string]map[string]*HostProxySwitchData),
		HostPortgroups:               make(map[string]map[string]*HostPortgroupData),
	}
	for _, dvpg := range dvpgs {
		if dvpg.Config.Key == "" || dvpg.Config.Name == "" {
			continue
		}

		if vlanInfo, ok := dvpg.Config.DefaultPortConfig.(*types.VMwareDVSPortSetting); ok {
			var vlanIDs []int
			var vlanIDRanges []string
			private := false
			switch v := vlanInfo.Vlan.(type) {
			case *types.VmwareDistributedVirtualSwitchTrunkVlanSpec:
				for _, item := range v.VlanId {
					switch {
					case item.Start == item.End:
						vlanIDs = append(vlanIDs, int(item.Start))
						vlanIDRanges = append(vlanIDRanges, fmt.Sprintf("%d", item.Start))
					case item.Start == constants.UntaggedVID && item.End == constants.MaxVID:
						vlanIDs = append(vlanIDs, constants.TaggedVID)
						vlanIDRanges = append(vlanIDRanges, fmt.Sprintf("%d-%d", item.Start, item.End))
					default:
						for vlan := item.Start; vlan <= item.End; vlan++ {
							vlanIDs = append(vlanIDs, int(vlan))
							vlanIDRanges = append(vlanIDRanges, fmt.Sprintf("%d-%d", item.Start, item.End))
						}
					}
				}
			case *types.VmwareDistributedVirtualSwitchPvlanSpec:
				vlanIDs = append(vlanIDs, int(v.PvlanId))
				private = true
			case *types.VmwareDistributedVirtualSwitchVlanIdSpec:
				vlanIDs = append(vlanIDs, int(v.VlanId))
			default:
				return fmt.Errorf("unknown vlan info spec %T", v)
			}

			for _, vid := range vlanIDs {
				if vid == constants.UntaggedVID || vid == constants.TaggedVID {
					continue
				}
				vc.Networks.Vid2Name[vid] = dvpg.Config.Name
			}

			vc.Networks.DistributedVirtualPortgroups[dvpg.Config.Key] = &DistributedPortgroupData{
				Name:         dvpg.Config.Name,
				VlanIDs:      vlanIDs,
				VlanIDRanges: vlanIDRanges,
				Private:      private,
			}
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
	vc.Disks = make(map[string]mo.Datastore, len(disks))
	for _, disk := range disks {
		vc.Disks[disk.Self.Value] = disk
	}
	return nil
}

func (vc *VmwareSource) InitDataCenters(ctx context.Context, containerView *view.ContainerView) error {
	var datacenters []mo.Datacenter
	err := containerView.Retrieve(ctx, []string{"Datacenter"}, []string{"name"}, &datacenters)
	if err != nil {
		return fmt.Errorf("failed retrieving datacenters: %s", err)
	}
	vc.DataCenters = make(map[string]mo.Datacenter, len(datacenters))
	for _, datacenter := range datacenters {
		vc.DataCenters[datacenter.Self.Value] = datacenter
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
	vc.Clusters = make(map[string]mo.ClusterComputeResource, len(clusters))
	for _, cluster := range clusters {
		vc.Clusters[cluster.Self.Value] = cluster
		for _, host := range cluster.Host {
			vc.Host2Cluster[host.Value] = cluster.Self.Value
		}
	}
	return nil
}

func (vc *VmwareSource) InitHosts(ctx context.Context, containerView *view.ContainerView) error {
	var hosts []mo.HostSystem
	err := containerView.Retrieve(ctx, []string{"HostSystem"}, []string{"name", "summary.host", "summary.hardware", "summary.runtime", "summary.config", "vm", "config.network"}, &hosts)
	if err != nil {
		return fmt.Errorf("failed retrieving hosts: %s", err)
	}
	vc.VM2Host = make(map[string]string)
	vc.Hosts = make(map[string]mo.HostSystem, len(hosts))
	for _, host := range hosts {
		vc.Hosts[host.Self.Value] = host
		for _, vm := range host.Vm {
			vc.VM2Host[vm.Value] = host.Self.Value
		}

		// Add network data which is received from hosts
		// Iterate over hosts virtual switches, needed to enrich data on physical interfaces
		vc.Networks.HostVirtualSwitches[host.Name] = make(map[string]*HostVirtualSwitchData)
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
		vc.Networks.HostProxySwitches[host.Name] = make(map[string]*HostProxySwitchData)
		for _, pswitch := range host.Config.Network.ProxySwitch {
			if pswitch.DvsUuid != "" {
				vc.Networks.HostProxySwitches[host.Name][pswitch.DvsUuid] = &HostProxySwitchData{
					mtu:   int(pswitch.Mtu),
					pnics: pswitch.Pnic,
					name:  pswitch.DvsName,
				}
			}
		}
		// Iterate over hosts port groups, needed to enrich data on physical interfaces
		vc.Networks.HostPortgroups[host.Name] = make(map[string]*HostPortgroupData)
		for _, pgroup := range host.Config.Network.Portgroup {
			if pgroup.Spec.Name != "" {
				nicOrder := pgroup.ComputedPolicy.NicTeaming.NicOrder
				pgroupNics := []string{}
				if len(nicOrder.ActiveNic) > 0 {
					pgroupNics = append(pgroupNics, nicOrder.ActiveNic...)
				}
				if len(nicOrder.StandbyNic) > 0 {
					pgroupNics = append(pgroupNics, nicOrder.StandbyNic...)
				}
				vc.Networks.HostPortgroups[host.Name][pgroup.Spec.Name] = &HostPortgroupData{
					vlanID:  int(pgroup.Spec.VlanId),
					vswitch: pgroup.Spec.VswitchName,
					nics:    pgroupNics,
				}
			}
		}
	}
	return nil
}

func (vc *VmwareSource) InitVms(ctx context.Context, containerView *view.ContainerView) error {
	var vms []mo.VirtualMachine
	err := containerView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary", "name", "runtime", "guest", "config.hardware", "config.guestFullName"}, &vms)
	if err != nil {
		return fmt.Errorf("failed retrieving vms: %s", err)
	}
	vc.Vms = make(map[string]mo.VirtualMachine, len(vms))
	for _, vm := range vms {
		vc.Vms[vm.Self.Value] = vm
	}
	return nil
}
