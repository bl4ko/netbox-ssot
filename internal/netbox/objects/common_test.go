// This file contains all objects that are common to all Netbox objects.
package objects

import "testing"

func TestChoice_String(t *testing.T) {
	tests := []struct {
		name string
		c    Choice
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Choice.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetboxObject_String(t *testing.T) {
	tests := []struct {
		name string
		n    NetboxObject
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("NetboxObject.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
