package objects

import (
	"fmt"
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
