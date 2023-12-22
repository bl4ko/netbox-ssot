package objects

type ClusterGroup struct {
	NetboxObject
	// Name is the name of the cluster group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the cluster group.
}

type ClusterType struct {
	NetboxObject
	// Name is the name of the cluster type. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the cluster type name. This field is read-only.
	Slug string `json:"slug,omitempty"`
}

type ClusterStatus struct {
	Choice
}

var (
	ClusterStatusActive  = ClusterStatus{Choice{Value: "active", Label: "Active"}}
	ClusterStatusOffline = ClusterStatus{Choice{Value: "offline", Label: "Offline"}}
)

type Cluster struct {
	NetboxObject
	// Name is the name of the cluster. This field is required.
	Name string `json:"name,omitempty"`
	// Type is the type of the cluster. This field is required.
	// e.g. oVirt,VMware...
	Type *ClusterType `json:"type,omitempty"`
	// ClusterGroup is the cluster group to which this cluster belongs.
	Group *ClusterGroup `json:"group,omitempty"`
	// Site is the site to which this cluster belongs.
	Site *Site `json:"site,omitempty"`
	// Status is the operational status of the cluster. This field is required.
	Status ClusterStatus `json:"status,omitempty"`
	// TenantGroup is the tenant group to which this cluster belongs.
	TenantGroup *TenantGroup `json:"tenant_group,omitempty"`
	// Tenant is the tenant to which this cluster belongs.
	Tenant *Tenant `json:"tenant,omitempty"`
}

type VMStatus struct {
	Choice
}

var (
	VMStatusActive  = VMStatus{Choice{Value: "active", Label: "Active"}}
	VMStatusOffline = VMStatus{Choice{Value: "offline", Label: "Offline"}}
)

// VM represents a virtual machine
type VM struct {
	NetboxObject
	// Name is the name of the virtual machine. This field is required.
	Name string `json:"name,omitempty"`
	// VMStatus is the status of the virtual machine. This field is required.
	VMStatus VMStatus `json:"status,omitempty"`
	// CPUs is the number of CPUs for the virtual machine.
	CPUs int `json:"vcpus,omitempty"`
	// RAM is the amount of RAM for the virtual machine in MB.
	RAM int `json:"memory,omitempty"`
	// Disk is the amount of disk space for the virtual machine in GB.
	Disk int `json:"disk,omitempty"`
	// Site is the site to which this virtual machine belongs.
	Site *Site `json:"site,omitempty"`
	// Cluster is the cluster to which this virtual machine belongs.
	Cluster *Cluster `json:"cluster,omitempty"`
	// Device is a specific host that this virtual machine is hosted on.
	Device *Device `json:"device,omitempty"`

	// Platform is the platform of the virtual machine.
	Platform *Platform `json:"platform,omitempty"`
}
