package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/virtualization"
)

type ClusterTypeResponse struct {
	Count    int                          `json:"count"`
	Next     int                          `json:"next"`
	Previous int                          `json:"previous"`
	Results  []virtualization.ClusterType `json:"results"`
}

// GET /api/virtualization/cluster-types/{id}
func (api *NetboxAPI) GetAllClusterTypes() ([]*virtualization.ClusterType, error) {
	api.Logger.Debug("Getting all cluster types from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/cluster-types/", nil)
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

	clusterTypes := make([]*virtualization.ClusterType, len(clusterTypeResponse.Results))
	for i := range clusterTypeResponse.Results {
		clusterTypes[i] = &clusterTypeResponse.Results[i]
	}
	api.Logger.Debug("Cluster types: ", clusterTypeResponse.Results)

	return clusterTypes, nil
}

// PATCH /api/virtualization/cluster-types/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchClusterType(diffMap map[string]interface{}, clusterTypeId int) (*virtualization.ClusterType, error) {
	api.Logger.Debug("Patching cluster group in NetBox")

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

	var patchedClusterType virtualization.ClusterType
	err = json.Unmarshal(response.Body, &patchedClusterType)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched cluster group: ", patchedClusterType, " with patchData: ", diffMap)

	return &patchedClusterType, nil

}

// POST /api/virtualization/cluster-types/
func (api *NetboxAPI) CreateClusterType(clusterType *virtualization.ClusterType) (*virtualization.ClusterType, error) {
	api.Logger.Debug("Creating cluster type in NetBox")

	requestBody, err := json.Marshal(clusterType)
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

	var createdClusterType virtualization.ClusterType
	err = json.Unmarshal(response.Body, &createdClusterType)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Created cluster type: ", createdClusterType)
	return &createdClusterType, nil
}

type ClusterGroupResponse struct {
	Count    int                           `json:"count"`
	Next     int                           `json:"next"`
	Previous int                           `json:"previous"`
	Results  []virtualization.ClusterGroup `json:"results"`
}

// GET /api/virtualization/cluster-groups/
func (api *NetboxAPI) GetAllClusterGroups() ([]*virtualization.ClusterGroup, error) {
	api.Logger.Debug("Getting all cluster groups from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/cluster-groups/", nil)
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

	clusterGroups := make([]*virtualization.ClusterGroup, len(clusterGroupResponse.Results))
	for i := range clusterGroupResponse.Results {
		clusterGroups[i] = &clusterGroupResponse.Results[i]
	}
	api.Logger.Debug("Cluster groups: ", clusterGroupResponse.Results)
	return clusterGroups, nil
}

// POST /api/virtualization/cluster-groups/
func (api *NetboxAPI) CreateClusterGroup(clusterGroup *virtualization.ClusterGroup) (*virtualization.ClusterGroup, error) {
	api.Logger.Debug("Creating cluster group in NetBox")

	clusterGroupJson, err := json.Marshal(clusterGroup)
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

	var createdClusterGroup virtualization.ClusterGroup
	err = json.Unmarshal(response.Body, &createdClusterGroup)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Created cluster group: ", createdClusterGroup)

	return &createdClusterGroup, nil
}

// PATCH /api/virtualization/cluster-groups/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchClusterGroup(diffMap map[string]interface{}, clusterGroupId int) (*virtualization.ClusterGroup, error) {
	api.Logger.Debug("Patching cluster group in NetBox")

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

	var patchedClusterGroup virtualization.ClusterGroup
	err = json.Unmarshal(response.Body, &patchedClusterGroup)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched cluster group: ", patchedClusterGroup)

	return &patchedClusterGroup, nil
}

type ClustersResponse struct {
	Count    int                      `json:"count"`
	Next     int                      `json:"next"`
	Previous int                      `json:"previous"`
	Results  []virtualization.Cluster `json:"results"`
}

// GET /api/vritualization/clusters
func (api *NetboxAPI) GetAllClusters() ([]*virtualization.Cluster, error) {
	api.Logger.Debug("Getting all clusters from NetBox")

	response, err := api.doRequest(MethodGet, "/api/virtualization/clusters/", nil)
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

	clusters := make([]*virtualization.Cluster, len(clustersResponse.Results))
	for i := range clustersResponse.Results {
		clusters[i] = &clustersResponse.Results[i]
	}

	api.Logger.Debug("Clusters: ", clustersResponse.Results)

	return clusters, nil
}

// PATCH /api/virtualization/clusters/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) PatchCluster(diffMap map[string]interface{}, clusterId int) (*virtualization.Cluster, error) {
	api.Logger.Debug("Patching cluster in NetBox")

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

	var patchedCluster virtualization.Cluster
	err = json.Unmarshal(response.Body, &patchedCluster)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Patched cluster: ", patchedCluster)

	return &patchedCluster, nil
}

// POST /api/virtualization/clusters/{id}/ -d '{"name": "new_name", ...}'
func (api *NetboxAPI) CreateCluster(cluster *virtualization.Cluster) (*virtualization.Cluster, error) {
	api.Logger.Debug("Creating cluster in NetBox")

	clusterJson, err := json.Marshal(cluster)
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

	var createdCluster virtualization.Cluster
	err = json.Unmarshal(response.Body, &createdCluster)
	if err != nil {
		return nil, err
	}

	api.Logger.Debug("Created cluster: ", createdCluster)

	return &createdCluster, nil
}
