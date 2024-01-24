package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/utils"
)

// PATCH /api/virtualization/cluster-types/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchClusterType(diffMap map[string]interface{}, clusterTypeId int) (*objects.ClusterType, error) {
	api.Logger.Debug("Patching cluster type with id ", clusterTypeId, "with data: ", diffMap)

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
	api.Logger.Debug("Creating cluster type in Netbox")

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

// POST /api/virtualization/cluster-groups/
func (api *NetboxAPI) CreateClusterGroup(clusterGroup *objects.ClusterGroup) (*objects.ClusterGroup, error) {
	api.Logger.Debug("Creating cluster group in Netbox")

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
	api.Logger.Debug("Patching cluster group with id ", clusterGroupId, ", with data: ", diffMap)

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

// PATCH /api/virtualization/clusters/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchCluster(diffMap map[string]interface{}, clusterId int) (*objects.Cluster, error) {
	api.Logger.Debug("Patching cluster with id ", clusterId, ", with data: ", diffMap)

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
	api.Logger.Debug("Creating cluster in Netbox")

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

// PATH /api/virtualization/virtual-machines/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchVM(diffMap map[string]interface{}, vmId int) (*objects.VM, error) {
	api.Logger.Debug("Patching VM with id ", vmId, " with data: ", diffMap)

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
	api.Logger.Debug("Creating VM in Netbox with data: ", vm)

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

// PATCH /api/virtualization/interfaces/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchVMInterface(diffMap map[string]interface{}, vmInterfaceId int) (*objects.VMInterface, error) {
	api.Logger.Debug("Patching VM interface with id ", vmInterfaceId, " with data: ", diffMap)

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
	api.Logger.Debug("Creating VM interface in Netbox with data: ", vmInterface)

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
