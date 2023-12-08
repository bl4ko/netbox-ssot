package tenancy

import "github.com/bl4ko/netbox-ssot/pkg/netbox/extras"

type TenantGroup struct {
	ID int `json:"id,omitempty"`
	// Name is the name of the tenant group. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the tenant group name. This field is read-only.
	Slug string `json:"slug,omitempty"`
	// Description is a description of the tenant group.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the tenant group.
	Tags []*extras.Tag `json:"tags,omitempty"`
}

type Tenant struct {
	ID int `json:"id,omitempty"`
	// Name is the name of the tenant. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slugified version of the tenant name. This field is read-only.
	Slug string `json:"slug,omitempty"`

	// Group is the tenant group to which this tenant belongs.
	Group *TenantGroup `json:"group,omitempty"`
	// Description is a description of the tenant.
	Description string `json:"description,omitempty"`
	// Tags is a list of tags for the tenant.
	Tags []*extras.Tag `json:"tags,omitempty"`
}
