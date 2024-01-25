package objects

import "fmt"

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

type ContactGroup struct {
	NetboxObject
	// Name is the name of the ContactGroup. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slug for the ContactGroup. This field is required.
	Slug string `json:"slug,omitempty"`
	// Parent contact group.
	Parent *ContactGroup `json:"parent,omitempty"`
}

// Default role name for admins of vms
const (
	AdminContactRoleName = "Admin"
)

// Contacts can be organized by functional roles.
// For example, we might create roles for administrative, emergency, operational contacts.
type ContactRole struct {
	NetboxObject
	// Name is the name of the role. This field is required.
	Name string `json:"name,omitempty"`
	// Slug is the slug of the role. This field is required.
	Slug string `json:"slug,omitempty"`
}

func (cr ContactRole) String() string {
	return fmt.Sprintf("ContactRole{Name: %s}", cr.Name)
}

type Contact struct {
	NetboxObject
	// Name is the name of the Contact. This field is required.
	Name string `json:"name,omitempty"`
	// Title is the title of the Contact.]
	Title string `json:"title,omitempty"`
	// Phone is the phone number of the contact.
	Phone string `json:"phone,omitempty"`
	// Email is the email of the contact.
	Email string `json:"email,omitempty"`
	// Address is the address of the contact.
	Address string `json:"address,omitempty"`
	// Link is the web link of the contact.
	Link string `json:"link,omitempty"`
	// Comments for the contact.
	Comments string `json:"comments,omitempty"`
}

func (c Contact) String() string {
	return fmt.Sprintf("Contact{Name: %s}", c.Name)
}
