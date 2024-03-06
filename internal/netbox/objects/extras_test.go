package objects

import "testing"

func TestTag_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Tag
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Tag.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomField_String(t *testing.T) {
	tests := []struct {
		name string
		cf   CustomField
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cf.String(); got != tt.want {
				t.Errorf("CustomField.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
