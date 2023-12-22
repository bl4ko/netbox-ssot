package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/objects"
	"github.com/bl4ko/netbox-ssot/pkg/utils"
)

type ClusterTypeResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []objects.ClusterType `json:"results"`
}

// GET /api/virtualization/cluster-types/?limit=0
func (api *NetboxAPI) GetAllClusterTypes() ([]*objects.ClusterType, error) {
	api.Logger.Debug("Getting all cluster types from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/cluster-types/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var clusterTypeResponse ClusterTypeResponse
	err = json.Unmarshal(response.Body, &clusterTypeResponse)
	if err != nil {
		return nil, err
	}

	clusterTypes := make([]*objects.ClusterType, len(clusterTypeResponse.Results))
	for i := range clusterTypeResponse.Results {
		clusterTypes[i] = &clusterTypeResponse.Results[i]
	}
	api.Logger.Debug("Cluster types: ", clusterTypeResponse.Results)

	return clusterTypes, nil
}

// PATCH /api/virtualization/cluster-types/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchClusterType(diffMap map[string]interface{}, clusterTypeId int) (*objects.ClusterType, error) {
	api.Logger.Debug("Patching cluster type in NetBox, with data: ", diffMap)

	patchBodyJson, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/virtualization/cluster-types/%d/", clusterTypeId), bytes.NewBuffer(patchBodyJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var patchedClusterType objects.ClusterType
	err = json.Unmarshal(response.Body, &patchedClusterType)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched ClusterType: ", patchedClusterType.Name, " with patchData: ", diffMap)

	return &patchedClusterType, nil

}

// POST /api/virtualization/cluster-types/
func (api *NetboxAPI) CreateClusterType(clusterType *objects.ClusterType) (*objects.ClusterType, error) {
	api.Logger.Debug("Creating cluster type in NetBox")

	requestBody, err := utils.NetboxJsonMarshal(clusterType)
	if err != nil {
		return nil, err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBody)

	response, err := api.doRequest(MethodPost, "/api/virtualization/cluster-types/", requestBodyBuffer)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var createdClusterType objects.ClusterType
	err = json.Unmarshal(response.Body, &createdClusterType)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Created cluster type: ", createdClusterType)
	return &createdClusterType, nil
}

type ClusterGroupResponse struct {
	Count    int                    `json:"count"`
	Next     *string                `json:"next"`
	Previous *string                `json:"previous"`
	Results  []objects.ClusterGroup `json:"results"`
}

// GET /api/virtualization/cluster-groups/?limit=0
func (api *NetboxAPI) GetAllClusterGroups() ([]*objects.ClusterGroup, error) {
	api.Logger.Debug("Getting all cluster groups from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/cluster-groups/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var clusterGroupResponse ClusterGroupResponse
	err = json.Unmarshal(response.Body, &clusterGroupResponse)
	if err != nil {
		return nil, err
	}

	clusterGroups := make([]*objects.ClusterGroup, len(clusterGroupResponse.Results))
	for i := range clusterGroupResponse.Results {
		clusterGroups[i] = &clusterGroupResponse.Results[i]
	}
	api.Logger.Debug("Cluster groups: ", clusterGroupResponse.Results)
	return clusterGroups, nil
}

// POST /api/virtualization/cluster-groups/
func (api *NetboxAPI) CreateClusterGroup(clusterGroup *objects.ClusterGroup) (*objects.ClusterGroup, error) {
	api.Logger.Debug("Creating cluster group in NetBox")

	clusterGroupJson, err := utils.NetboxJsonMarshal(clusterGroup)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPost, "/api/virtualization/cluster-groups/", bytes.NewBuffer(clusterGroupJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var createdClusterGroup objects.ClusterGroup
	err = json.Unmarshal(response.Body, &createdClusterGroup)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Created cluster group: ", createdClusterGroup)

	return &createdClusterGroup, nil
}

// PATCH /api/virtualization/cluster-groups/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchClusterGroup(diffMap map[string]interface{}, clusterGroupId int) (*objects.ClusterGroup, error) {
	api.Logger.Debug("Patching cluster group in NetBox, with data: ", diffMap)

	clusterGroupJson, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/virtualization/cluster-groups/%d/", clusterGroupId), bytes.NewBuffer(clusterGroupJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var patchedClusterGroup objects.ClusterGroup
	err = json.Unmarshal(response.Body, &patchedClusterGroup)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched ClusterGroup: ", patchedClusterGroup.Name, " with patchData: ", diffMap)

	return &patchedClusterGroup, nil
}

type ClustersResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []objects.Cluster `json:"results"`
}

// GET /api/vritualization/clusters/?limit=0
func (api *NetboxAPI) GetAllClusters() ([]*objects.Cluster, error) {
	api.Logger.Debug("Getting all clusters from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/clusters/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var clustersResponse ClustersResponse
	err = json.Unmarshal(response.Body, &clustersResponse)
	if err != nil {
		return nil, err
	}

	clusters := make([]*objects.Cluster, len(clustersResponse.Results))
	for i := range clustersResponse.Results {
		clusters[i] = &clustersResponse.Results[i]
	}

	api.Logger.Debug("Clusters: ", clustersResponse.Results)

	return clusters, nil
}

// PATCH /api/virtualization/clusters/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchCluster(diffMap map[string]interface{}, clusterId int) (*objects.Cluster, error) {
	api.Logger.Debug("Patching cluster in NetBox, with data: ", diffMap)

	clusterJson, err := json.Marshal(diffMap)
	fmt.Println(string(clusterJson))
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/virtualization/clusters/%d/", clusterId), bytes.NewBuffer(clusterJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var patchedCluster objects.Cluster
	err = json.Unmarshal(response.Body, &patchedCluster)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched cluster: ", patchedCluster)

	return &patchedCluster, nil
}

// POST /api/virtualization/clusters/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) CreateCluster(cluster *objects.Cluster) (*objects.Cluster, error) {
	api.Logger.Debug("Creating cluster in NetBox")

	clusterJson, err := utils.NetboxJsonMarshal(cluster)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Cluster JSON: ", string(clusterJson))

	response, err := api.doRequest(MethodPost, "/api/virtualization/clusters/", bytes.NewBuffer(clusterJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var createdCluster objects.Cluster
	err = json.Unmarshal(response.Body, &createdCluster)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Created cluster: ", createdCluster)

	return &createdCluster, nil
}

type VMResponse struct {
	Count    int          `json:"count"`
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Results  []objects.VM `json:"results"`
}

// GET /api/virtualization/virtual-machines/?limit=0
func (api *NetboxAPI) GetAllVMs() ([]*objects.VM, error) {
	api.Logger.Debug("Getting all virtual machines from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/virtual-machines/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d. Error %s", response.StatusCode, response.Body)
	}

	var vmResponse VMResponse
	err = json.Unmarshal(response.Body, &vmResponse)
	if err != nil {
		return nil, err
	}

	vms := make([]*objects.VM, len(vmResponse.Results))
	for i := range vmResponse.Results {
		vms[i] = &vmResponse.Results[i]
	}

	api.Logger.Debug("Successfully received virtual machines: ", vmResponse.Results)

	return vms, nil
}

// PATH /api/virtualization/virtual-machines/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchVM(diffMap map[string]interface{}, vmId int) (*objects.VM, error) {
	api.Logger.Debug("Patching VM with id", vmId, " with data: ", diffMap)

	vmJson, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/virtualization/virtual-machines/%d/", vmId), bytes.NewBuffer(vmJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var patchedVM objects.VM
	err = json.Unmarshal(response.Body, &patchedVM)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched VM: ", patchedVM)

	return &patchedVM, nil
}

// POST /api/virtualization/virtual-machines/
func (api *NetboxAPI) CreateVM(vm *objects.VM) (*objects.VM, error) {
	api.Logger.Debug("Creating VM in NetBox with data: ", vm)

	vmJson, err := utils.NetboxJsonMarshal(vm)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPost, "/api/virtualization/virtual-machines/", bytes.NewBuffer(vmJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d: %s", response.StatusCode, response.Body)
	}

	var createdVM objects.VM
	err = json.Unmarshal(response.Body, &createdVM)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created VM: ", createdVM)

	return &createdVM, nil
}

type VMInterfaceResponse struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"`
	Previous *string               `json:"previous"`
	Results  []objects.VMInterface `json:"results"`
}

// GET /api/virtualization/interfaces/?limit=0
func (api *NetboxAPI) GetAllVMInterfaces() ([]*objects.VMInterface, error) {
	api.Logger.Debug("Getting all VM interfaces from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/interfaces/?limit=0", nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, response.Body)
	}

	var vmInterfaceResponse VMInterfaceResponse
	err = json.Unmarshal(response.Body, &vmInterfaceResponse)
	if err != nil {
		return nil, err
	}

	vmInterfaces := make([]*objects.VMInterface, len(vmInterfaceResponse.Results))
	for i := range vmInterfaceResponse.Results {
		vmInterfaces[i] = &vmInterfaceResponse.Results[i]
	}

	api.Logger.Debug("Successfully received VM interfaces: ", vmInterfaceResponse.Results)

	return vmInterfaces, nil
}

// PATCH /api/virtualization/interfaces/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchVMInterface(diffMap map[string]interface{}, vmInterfaceId int) (*objects.VMInterface, error) {
	api.Logger.Debug("Patching VM interface with id", vmInterfaceId, " with data: ", diffMap)

	vmInterfaceJson, err := json.Marshal(diffMap)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPatch, fmt.Sprintf("/api/virtualization/interfaces/%d/", vmInterfaceId), bytes.NewBuffer(vmInterfaceJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, response.Body)
	}

	var patchedVMInterface objects.VMInterface
	err = json.Unmarshal(response.Body, &patchedVMInterface)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully patched VM interface: ", patchedVMInterface)

	return &patchedVMInterface, nil
}

// POST /api/virtualization/interfaces/
func (api *NetboxAPI) CreateVMInterface(vmInterface *objects.VMInterface) (*objects.VMInterface, error) {
	api.Logger.Debug("Creating VM interface in NetBox with data: ", vmInterface)

	vmInterfaceJson, err := utils.NetboxJsonMarshal(vmInterface)
	if err != nil {
		return nil, err
	}

	response, err := api.doRequest(MethodPost, "/api/virtualization/interfaces/", bytes.NewBuffer(vmInterfaceJson))
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, response.Body)
	}

	var createdVMInterface objects.VMInterface
	err = json.Unmarshal(response.Body, &createdVMInterface)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Successfully created VM interface: ", createdVMInterface)

	return &createdVMInterface, nil
}
