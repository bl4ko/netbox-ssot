package objects

import (
	"fmt"
	"testing"
)

func TestSite_String(t *testing.T) {
	tests := []struct {
		name string
		s    Site
		want string
	}{
		{
			name: "Correct string output for site",
			s: Site{
				Name:      "Test site",
				Slug:      "test_site",
				Latitude:  68.034,
				Longitude: 69.324,
			},
			want: "Site{Name: Test site}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("Site.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlatform_String(t *testing.T) {
	tests := []struct {
		name string
		p    Platform
		want string
	}{
		{
			name: "Correct string output for platform",
			p: Platform{
				Name: "TestPlatform",
				Slug: "testplatform",
				Manufacturer: &Manufacturer{
					Name: "TestManufacturer",
				},
			},
			want: fmt.Sprintf("Platform{Name: %s, Manufacturer: %s}", "TestPlatform", Manufacturer{Name: "TestManufacturer"}),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Platform.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManufacturer_String(t *testing.T) {
	tests := []struct {
		name string
		m    Manufacturer
		want string
	}{
		{
			name: "Correct string representation of manufacturer",
			m: Manufacturer{
				Name: "Test manufacturer",
				Slug: "test_manufacturer",
			},
			want: "Manufacturer{Name: Test manufacturer}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Manufacturer.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceType_String(t *testing.T) {
	tests := []struct {
		name string
		dt   DeviceType
		want string
	}{
		{
			name: "Correct string representation of DeviceType",
			dt: DeviceType{
				Manufacturer: &Manufacturer{
					Name: "Test manufacturer",
				},
				Model: "test model",
			},
			want: fmt.Sprintf("DeviceType{Manufacturer: %s, Model: %s}", "Test manufacturer", "test model"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.String(); got != tt.want {
				t.Errorf("DeviceType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceRole_String(t *testing.T) {
	tests := []struct {
		name string
		dr   DeviceRole
		want string
	}{
		{
			name: "Correct string representation of Device Role",
			dr: DeviceRole{
				Name: "Test device-role",
			},
			want: "DeviceRole{Name: Test device-role}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dr.String(); got != tt.want {
				t.Errorf("DeviceRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_String(t *testing.T) {
	tests := []struct {
		name string
		d    Device
		want string
	}{
		{
			name: "Correct string representation of Device",
			d: Device{
				Name: "Test device",
				DeviceType: &DeviceType{
					Manufacturer: &Manufacturer{
						Name: "Test manufacturer",
					},
					Model: "test model",
				},
				DeviceRole: &DeviceRole{
					Name: "Test device-role",
				},
				Site: &Site{
					Name: "Test site",
				},
			},
			want: fmt.Sprintf("Device{Name: %s, Type: %s, Role: %s, Site: %s}", "Test device", "DeviceType{Manufacturer: Test manufacturer, Model: test model}", "DeviceRole{Name: Test device-role}", "Site{Name: Test site}"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Device.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_String(t *testing.T) {
	tests := []struct {
		name string
		i    Interface
		want string
	}{
		{
			name: "Correct string representation of Interface",
			i: Interface{
				Name: "Test interface",
				Type: &OtherInterfaceType,
				Device: &Device{
					Name: "Test device",
					DeviceType: &DeviceType{
						Manufacturer: &Manufacturer{
							Name: "Test manufacturer",
						},
						Model: "test model",
					},
					DeviceRole: &DeviceRole{
						Name: "Test device-role",
					},
					Site: &Site{
						Name: "Test site",
					},
				},
			},
			want: fmt.Sprintf("Interface{Name: %s, Device: %s, Type: %s}", "Test interface", "Test device", OtherInterfaceType.Label),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("Interface.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualDeviceContext_String(t *testing.T) {
	tests := []struct {
		name string
		vdc  VirtualDeviceContext
		want string
	}{
		{
			name: "Correct string representation of virtual device context",
			vdc: VirtualDeviceContext{
				Name:   "test",
				Device: &Device{Name: "testdevice"},
				Status: &VDCStatusActive,
			},
			want: fmt.Sprintf("VirtualDeviceContext{Name: %s, Device: %s, Status: %s}", "test", &Device{Name: "testdevice"}, &VDCStatusActive),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vdc.String(); got != tt.want {
				t.Errorf("VirtualDeviceContext.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
