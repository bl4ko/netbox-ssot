package ovirt

import (
	"fmt"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// Fetches networks from ovirt api and stores them to local object.
func (o *OVirtSource) initNetworks(conn *ovirtsdk4.Connection) error {
	networksResponse, err := conn.SystemService().
		NetworksService().
		List().
		Follow("vnicprofiles").
		Send()
	if err != nil {
		return fmt.Errorf("init oVirt networks: %v", err)
	}
	o.Networks = &NetworkData{
		OVirtNetworks:       make(map[string]*ovirtsdk4.Network),
		Vid2Name:            make(map[int]string),
		VnicProfile2Network: make(map[string]string),
	}
	if networks, ok := networksResponse.Networks(); ok {
		for _, network := range networks.Slice() {
			if networkID, ok := network.Id(); ok {
				o.Networks.OVirtNetworks[networkID] = network
				if vlan, exists := network.Vlan(); exists {
					if vlanID, exists := vlan.Id(); exists {
						o.Networks.Vid2Name[int(vlanID)] = network.MustName()
					}
				}
				if vnicProfiles, ok := network.VnicProfiles(); ok {
					for _, vnicProfile := range vnicProfiles.Slice() {
						if vnicProfileID, ok := vnicProfile.Id(); ok {
							o.Networks.VnicProfile2Network[vnicProfileID] = networkID
						}
					}
				}
			}
		}
		o.Logger.Debug(o.Ctx, "Successfully initialized oVirt networks: ", o.Networks)
	} else {
		o.Logger.Warning(o.Ctx, "Error initializing oVirt networks")
	}
	return nil
}

func (o *OVirtSource) initDisks(conn *ovirtsdk4.Connection) error {
	// Get the disks
	disksResponse, err := conn.SystemService().DisksService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt disks: %v", err)
	}
	o.Disks = make(map[string]*ovirtsdk4.Disk)
	if disks, ok := disksResponse.Disks(); ok {
		for _, disk := range disks.Slice() {
			o.Disks[disk.MustId()] = disk
		}
		o.Logger.Debug(o.Ctx, "Successfully initialized oVirt disks: ", o.Disks)
	} else {
		o.Logger.Warning(o.Ctx, "Error initializing oVirt disks")
	}
	return nil
}

func (o *OVirtSource) initDataCenters(conn *ovirtsdk4.Connection) error {
	dataCentersResponse, err := conn.SystemService().DataCentersService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt data centers: %v", err)
	}
	o.DataCenters = make(map[string]*ovirtsdk4.DataCenter)
	if dataCenters, ok := dataCentersResponse.DataCenters(); ok {
		for _, dataCenter := range dataCenters.Slice() {
			o.DataCenters[dataCenter.MustId()] = dataCenter
		}
		o.Logger.Debug(o.Ctx, "Successfully initialized oVirt data centers: ", o.DataCenters)
	} else {
		o.Logger.Warning(o.Ctx, "Error initializing oVirt data centers")
	}
	return nil
}

// Function that queries ovirt api for clusters and stores them locally.
func (o *OVirtSource) initClusters(conn *ovirtsdk4.Connection) error {
	clustersResponse, err := conn.SystemService().ClustersService().List().Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt clusters: %v", err)
	}
	o.Clusters = make(map[string]*ovirtsdk4.Cluster)
	if clusters, ok := clustersResponse.Clusters(); ok {
		for _, cluster := range clusters.Slice() {
			o.Clusters[cluster.MustId()] = cluster
		}
		o.Logger.Debug(o.Ctx, "Successfully initialized oVirt clusters: ", o.Clusters)
	} else {
		o.Logger.Warning(o.Ctx, "Error initializing oVirt clusters")
	}
	return nil
}

// Function that queries ovirt api for hosts and stores them locally.
func (o *OVirtSource) initHosts(conn *ovirtsdk4.Connection) error {
	hostsResponse, err := conn.SystemService().HostsService().List().Follow("nics").Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt hosts: %+v", err)
	}
	o.Hosts = make(map[string]*ovirtsdk4.Host)
	if hosts, ok := hostsResponse.Hosts(); ok {
		for _, host := range hosts.Slice() {
			o.Hosts[host.MustId()] = host
		}
		o.Logger.Debug(o.Ctx, "Successfully initialized oVirt hosts: ", hosts)
	} else {
		o.Logger.Warning(o.Ctx, "Error initializing oVirt hosts")
	}
	return nil
}

// Function that queries the ovirt api for vms and stores them locally.
func (o *OVirtSource) initVms(conn *ovirtsdk4.Connection) error {
	vmsResponse, err := conn.SystemService().
		VmsService().
		List().
		Follow("nics,diskattachments,reporteddevices").
		Send()
	if err != nil {
		return fmt.Errorf("failed to get oVirt vms: %+v", err)
	}
	o.Vms = make(map[string]*ovirtsdk4.Vm)
	if vms, ok := vmsResponse.Vms(); ok {
		for _, vm := range vms.Slice() {
			o.Vms[vm.MustId()] = vm
		}
		o.Logger.Debug(o.Ctx, "Successfully initialized oVirt vms: ", vms)
	} else {
		o.Logger.Warning(o.Ctx, "Error initializing oVirt vms")
	}
	return nil
}
