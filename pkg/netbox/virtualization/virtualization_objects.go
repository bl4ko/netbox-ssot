package virtualization

import (
	"encoding/json"

	"github.com/bl4ko/netbox-ssot/pkg/netbox/dcim"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/extras"
	"github.com/bl4ko/netbox-ssot/pkg/netbox/tenancy"
)

type ClusterGroup struct {
	// ID is the unique identifier of the cluster group.
	ID int `json:"id,omitempty"`
	// Name is the name of the cluster group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the cluster group.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the cluster group.
	Tags []*extras.Tag `json:"tags,omitempty"`
}

type ClusterType struct {
	ID int `json:"id,omitempty"`
	// Name is the name of the cluster type. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster type name. This field is read-only.
	Slug string `json:"slug,omitempty"`

	// Description is a description of the cluster type.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the cluster type.
	Tags []*extras.Tag `json:"tags,omitempty"`
}

type Cluster struct {
	ID int `json:"id,omitempty"`
	// Name is the name of the cluster. This field is required.
	Name string `json:"name,omitempty"`
	// Type is the type of the cluster. This field is required.
	// e.g. oVirt,VMware...
	Type *ClusterType `json:"type,omitempty"`
	// ClusterGroup is the cluster group to which this cluster belongs.
	Group *ClusterGroup `json:"group,omitempty"`
	// Site is the site to which this cluster belongs.
	Site *dcim.Site `json:"site,omitempty"`
	// Status is the operational status of the cluster. This field is required.
	Status *dcim.Status `json:"status,omitempty"`
	// TenantGroup is the tenant group to which this cluster belongs.
	TenantGroup *tenancy.TenantGroup `json:"tenant_group,omitempty"`
	// Tenant is the tenant to which this cluster belongs.
	Tenant *tenancy.Tenant `json:"tenant,omitempty"`
	// Description is a description of the cluster.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the cluster.
	Tags []*extras.Tag `json:"tags,omitempty"`
}

// Custom marshal because we shouldn't pass the status as object but as string
func (c *Cluster) MarshalJSON() ([]byte, error) {
	type Alias Cluster
	return json.Marshal(&struct {
		Status string `json:"status,omitempty"`
		*Alias
	}{
		Status: c.Status.Value,
		Alias:  (*Alias)(c),
	})
}

// VM represents a virtual machine
type VM struct {
	// Name is the name of the virtual machine. This field is required.
	Name string `json:"name,omitempty"`
	// VMStatus is the status of the virtual machine. This field is required.
	VMStatus *dcim.Status `json:"status,omitempty"`

	// CPUs is the number of CPUs for the virtual machine.
	CPUs int `json:"vcpus,omitempty"`
	// RAM is the amount of RAM for the virtual machine in MB.
	RAM int `json:"memory,omitempty"`
	// Disk is the amount of disk space for the virtual machine in GB.
	Disk int `json:"disk,omitempty"`

	// Description is a description of the virtual machine.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the virtual machine.
	Tags []*extras.Tag `json:"tags,omitempty"`

	// Platform is the platform of the virtual machine.
	Platform *dcim.Platform `json:"platform,omitempty"`
}
