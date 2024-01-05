package objects

type TenantGroup struct {
	NetboxObject
	// Name is the name of the tenant group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the tenant group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the tenant group.
}

type Tenant struct {
	NetboxObject
	// Name is the name of the tenant. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the tenant name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Group is the tenant group to which this tenant belongs.
	Group *TenantGroup `json:"group,omitempty"`
}
