package vmware

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/source/common"
	"github.com/bl4ko/netbox-ssot/internal/utils"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

// VmwareSource represents an vsphere source.
//
//nolint:revive
type VmwareSource struct {
	common.Config
	// Vmware API data initialized in init functions
	Disks       map[string]mo.Datastore
	DataCenters map[string]mo.Datacenter
	Clusters    map[string]mo.ClusterComputeResource
	Hosts       map[string]mo.HostSystem
	Vms         map[string]mo.VirtualMachine
	Networks    NetworkData

	// Relations between objects "object_id": "object_id"
	Cluster2Datacenter map[string]string // ClusterKey -> DatacenterKey
	Host2Cluster       map[string]string // HostKey -> ClusterKey
	VM2Host            map[string]string // VmKey ->  HostKey

	// CustomField2Name is a map of custom field ids to their names
	CustomFieldID2Name map[int32]string
	// Object2Tags is a map of object ids to their tags
	Object2Tags   map[string][]*tags.Tag
	Object2NBTags map[string][]*objects.Tag // Created in sync function
}

type NetworkData struct {
	DistributedVirtualPortgroups map[string]*DistributedPortgroupData         // Portgroup.key -> PortgroupData
	Vid2Name                     map[int]string                               // Helper map, for quickly obtaining name of the vid
	HostVirtualSwitches          map[string]map[string]*HostVirtualSwitchData // hostName -> VSwitchName-> VSwitchData
	HostProxySwitches            map[string]map[string]*HostProxySwitchData   // hostName -> PSwitchName ->
	HostPortgroups               map[string]map[string]*HostPortgroupData     // hostname -> Portgroup.Spec.Name -> HostPortgroupData
}

type DistributedPortgroupData struct {
	Name         string
	VlanIDs      []int
	VlanIDRanges []string
	Private      bool
	Tenant       *objects.Tenant
}

type HostVirtualSwitchData struct {
	mtu   int
	pnics []string
}

type HostProxySwitchData struct {
	name  string
	mtu   int
	pnics []string
}

type HostPortgroupData struct {
	vlanID  int
	vswitch string
	nics    []string
}

func (vc *VmwareSource) Init() error {
	// Initialize the connection
	vc.Logger.Debug(vc.Ctx, "vmware source ", vc.SourceConfig.Name)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Correctly handle backslashes in username and password
	escapedUsername := url.PathEscape(vc.SourceConfig.Username)
	escapedPassword := url.PathEscape(vc.SourceConfig.Password)

	vcURL := fmt.Sprintf("%s://%s:%s@%s:%d/sdk", vc.SourceConfig.HTTPScheme, escapedUsername, escapedPassword, vc.SourceConfig.Hostname, vc.SourceConfig.Port)

	url, err := url.Parse(vcURL)
	if err != nil {
		return fmt.Errorf("failed parsing url for %s with error %s", vc.SourceConfig.Hostname, err)
	}

	// How to set custom ca certificates for govmomi: https://github.com/vmware/govmomi/issues/1200#issuecomment-412950179
	soapClient := soap.NewClient(url, !vc.SourceConfig.ValidateCert)
	if vc.Config.CAFile != "" {
		err = soapClient.SetRootCAs(vc.Config.CAFile)
		if err != nil {
			return fmt.Errorf("set root CAs: %s", err)
		}
	}
	vim25Client, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return fmt.Errorf("failed creating a govmomi client with an error: %s", err)
	}

	// Create a SessionManager and login to authenticate the session
	sessionManager := session.NewManager(vim25Client)

	// Perform login
	if err = sessionManager.Login(ctx, url.User); err != nil {
		return fmt.Errorf("login failed: %s", err)
	}

	// View manager is used to create and manage views. Views are a mechanism in vSphere
	// to group and manage objects in the inventory.
	viewManager := view.NewManager(vim25Client)

	// viewType specifies the types of objects to be included in our container view.
	// Each string in this slice represents a different vSphere Managed Object type.
	viewType := []string{
		"Datastore", "Datacenter", "ClusterComputeResource", "HostSystem", "VirtualMachine", "Network",
	}

	// A container view is a subset of the vSphere inventory, focusing on the specified
	// object types, making it easier to manage and retrieve data for these objects.
	containerView, err := viewManager.CreateContainerView(ctx, vim25Client.ServiceContent.RootFolder, viewType, true)
	if err != nil {
		return fmt.Errorf("failed creating containerView: %s", err)
	}

	vc.Logger.Debug(vc.Ctx, "Connection to vmware source ", vc.SourceConfig.Hostname, " established successfully")

	// Create CustomFieldManager to map custom field ids to their names
	// This is required to determine which custom field key is used for
	// which custom field name (e.g.g 202 -> vm owner, 203 -> vm description...)
	err = vc.CreateCustomFieldRelation(ctx, vim25Client)
	if err != nil {
		return fmt.Errorf("create custom field relation failed: %s", err)
	}

	// Find relation between data centers and clusters. Currently we have to manually traverse
	// the tree to get this relation.
	err = vc.CreateClusterDataCenterRelation(ctx, vim25Client)
	if err != nil {
		return fmt.Errorf("create cluster datacenter relation failed: %s", err)
	}

	err = vc.CreateObjectTagsRelation(ctx, vim25Client, url.User)
	if err != nil {
		return fmt.Errorf("create object tags relation failed: %s", err)
	}

	// Initialize items from vsphere API to local storage
	initFunctions := []func(context.Context, *view.ContainerView) error{
		vc.initNetworks,
		vc.initDisks,
		vc.initDataCenters,
		vc.initClusters,
		vc.initHosts,
		vc.initVms,
	}

	for _, initFunc := range initFunctions {
		startTime := time.Now()
		if err := initFunc(ctx, containerView); err != nil {
			return fmt.Errorf("vmware initialization failure: %v", err)
		}
		duration := time.Since(startTime)
		vc.Logger.Infof(vc.Ctx, "Successfully initialized %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(initFunc, "init"), duration.Seconds())
	}

	// Ensure the containerView is destroyed after we are done with it
	err = containerView.Destroy(ctx)
	if err != nil {
		vc.Logger.Errorf(vc.Ctx, "failed destroying containerView: %s", err)
	}

	err = sessionManager.Logout(ctx)
	if err != nil {
		return fmt.Errorf("error occurred when ending vmware connection to host %s: %s", vc.SourceConfig.Hostname, err)
	}

	vc.Logger.Debug(vc.Ctx, "Successfully closed connection to vmware host: ", vc.SourceConfig.Hostname)

	return nil
}

// Function that syncs all data from oVirt to Netbox.
func (vc *VmwareSource) Sync(nbi *inventory.NetboxInventory) error {
	syncFunctions := []func(*inventory.NetboxInventory) error{
		vc.syncTags,
		vc.syncNetworks,
		vc.syncDatacenters,
		vc.syncClusters,
		vc.syncHosts,
		vc.syncVMs,
	}
	for _, syncFunc := range syncFunctions {
		startTime := time.Now()
		err := syncFunc(nbi)
		if err != nil {
			return err
		}
		duration := time.Since(startTime)
		vc.Logger.Infof(vc.Ctx, "Successfully synced %s in %f seconds", utils.ExtractFunctionNameWithTrimPrefix(syncFunc, "sync"), duration.Seconds())
	}
	return nil
}

// Currently we have to traverse the vsphere tree to get datacenter to cluster relation
// For other objects relations are available in with containerView.
func (vc *VmwareSource) CreateClusterDataCenterRelation(ctx context.Context, client *vim25.Client) error {
	finder := find.NewFinder(client, true)
	datacenters, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		return fmt.Errorf("finder failed creating datacenter list: %s", err)
	}
	vc.Cluster2Datacenter = make(map[string]string)
	for _, dc := range datacenters {
		finder.SetDatacenter(dc)
		clusters, err := finder.ClusterComputeResourceList(ctx, "*")
		if err != nil {
			return fmt.Errorf("finder failed finding clusters for datacenter %s: %s", dc.InventoryPath, err)
		}
		for _, cluster := range clusters {
			vc.Cluster2Datacenter[cluster.Reference().Value] = dc.Reference().Value
		}
	}
	return nil
}

// Creates a map of custom field ids to their names.
func (vc *VmwareSource) CreateCustomFieldRelation(ctx context.Context, client *vim25.Client) error {
	cfm, err := object.GetCustomFieldsManager(client)
	if err != nil {
		return fmt.Errorf("createCustomFieldRelation: %s", err)
	}
	fieldDefs, err := cfm.Field(ctx)
	if err != nil {
		return fmt.Errorf("createCustomFieldRelation fieldDefs: %s", err)
	}

	vc.CustomFieldID2Name = make(map[int32]string)
	for _, field := range fieldDefs {
		vc.CustomFieldID2Name[field.Key] = field.Name
	}

	return nil
}

// Creates a map of object ids to their tags. This uses the rest api since tags
// are not available in the containerView (see https://github.com/vmware/govmomi/issues/1825).
func (vc *VmwareSource) CreateObjectTagsRelation(ctx context.Context, vim25client *vim25.Client, userInfo *url.Userinfo) error {
	restClient := rest.NewClient(vim25client)

	_, err := restClient.Session(ctx)
	if err != nil {
		return fmt.Errorf("failed creating rest client session: %s", err)
	}

	if err = restClient.Login(ctx, userInfo); err != nil {
		return fmt.Errorf("rest client login failed: %s", err)
	}
	defer func() {
		err := restClient.Logout(ctx)
		if err != nil {
			vc.Logger.Errorf(vc.Ctx, "failed logging out from rest client: %s", err)
		}
	}()

	tagManager := tags.NewManager(restClient)
	tagList, err := tagManager.ListTags(ctx)
	if err != nil {
		return fmt.Errorf("failed listing tags: %s", err)
	}

	objectNames2tags := make(map[string][]*tags.Tag)
	for _, tag := range tagList {
		objs, err := tagManager.GetAttachedObjectsOnTags(ctx, []string{tag})
		if err != nil {
			return fmt.Errorf("failed getting attached objects on tag %+v: %s", tag, err)
		}
		if len(objs) > 0 {
			tagInfo, err := tagManager.GetTag(ctx, tag)
			if err != nil {
				return fmt.Errorf("failed getting tag %+v info: %s", tag, err)
			}
			for _, elem := range objs[0].ObjectIDs {
				objectNames2tags[elem.Reference().Value] = append(objectNames2tags[elem.Reference().Value], tagInfo)
			}
		}
	}

	vc.Object2Tags = objectNames2tags
	return nil
}
