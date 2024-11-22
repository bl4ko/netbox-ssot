package objects

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
)

func TestTenant_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Tenant
		want string
	}{
		{
			name: "Test tenant correct string",
			tr: Tenant{
				Name: "Test tenant",
			},
			want: "Tenant{Name: Test tenant}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Tenant.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactRole_String(t *testing.T) {
	tests := []struct {
		name string
		cr   ContactRole
		want string
	}{
		{
			name: "Test contact role correct string",
			cr: ContactRole{
				Name: "Test contact role",
			},
			want: "ContactRole{Name: Test contact role}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cr.String(); got != tt.want {
				t.Errorf("ContactRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContact_String(t *testing.T) {
	tests := []struct {
		name string
		c    Contact
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Contact.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactAssignment_String(t *testing.T) {
	tests := []struct {
		name string
		ca   ContactAssignment
		want string
	}{
		{
			name: "Test contact assignment correct string",
			ca: ContactAssignment{
				NetboxObject: NetboxObject{
					ID: 1,
				},
				ModelType: constants.ContentTypeVirtualizationVirtualMachine,
				Contact: &Contact{
					Name: "Test contact",
				},
				Role: &ContactRole{
					Name: "Test contact role",
				},
				ObjectID: 5,
			},
			want: fmt.Sprintf("ContactAssignment{ObjectType: %s, ObjectID: %d, %v, %v}", constants.ContentTypeVirtualizationVirtualMachine, 5, Contact{Name: "Test contact"}, ContactRole{Name: "Test contact role"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ca.String(); got != tt.want {
				t.Errorf("ContactAssignment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantGroup_GetID(t *testing.T) {
	tests := []struct {
		name string
		tg   *TenantGroup
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tg.GetID(); got != tt.want {
				t.Errorf("TenantGroup.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenantGroup_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		tg   *TenantGroup
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tg.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TenantGroup.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenant_GetID(t *testing.T) {
	tests := []struct {
		name string
		tr   *Tenant
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetID(); got != tt.want {
				t.Errorf("Tenant.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTenant_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		tr   *Tenant
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tenant.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactGroup_String(t *testing.T) {
	tests := []struct {
		name string
		cg   ContactGroup
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cg.String(); got != tt.want {
				t.Errorf("ContactGroup.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactGroup_GetID(t *testing.T) {
	tests := []struct {
		name string
		cg   *ContactGroup
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cg.GetID(); got != tt.want {
				t.Errorf("ContactGroup.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactGroup_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		cg   *ContactGroup
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cg.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContactGroup.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactRole_GetID(t *testing.T) {
	tests := []struct {
		name string
		cr   *ContactRole
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cr.GetID(); got != tt.want {
				t.Errorf("ContactRole.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactRole_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		cr   *ContactRole
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cr.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContactRole.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContact_GetID(t *testing.T) {
	tests := []struct {
		name string
		c    *Contact
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetID(); got != tt.want {
				t.Errorf("Contact.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContact_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		c    *Contact
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Contact.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactAssignment_GetID(t *testing.T) {
	tests := []struct {
		name string
		ca   *ContactAssignment
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ca.GetID(); got != tt.want {
				t.Errorf("ContactAssignment.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactAssignment_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		ca   *ContactAssignment
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ca.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContactAssignment.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
