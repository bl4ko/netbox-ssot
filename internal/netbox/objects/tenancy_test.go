package objects

import "testing"

func TestTenant_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Tenant
		want string
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ca.String(); got != tt.want {
				t.Errorf("ContactAssignment.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
