package objects

import (
	"fmt"
	"reflect"
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

func TestSite_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		s    *Site
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Site.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSite_GetID(t *testing.T) {
	tests := []struct {
		name string
		s    *Site
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GetID(); got != tt.want {
				t.Errorf("Site.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlatform_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		p    *Platform
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Platform.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlatform_GetID(t *testing.T) {
	tests := []struct {
		name string
		p    *Platform
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetID(); got != tt.want {
				t.Errorf("Platform.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegion_String(t *testing.T) {
	tests := []struct {
		name string
		r    Region
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("Region.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegion_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		r    *Region
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Region.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegion_GetID(t *testing.T) {
	tests := []struct {
		name string
		r    *Region
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetID(); got != tt.want {
				t.Errorf("Region.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_String(t *testing.T) {
	tests := []struct {
		name string
		l    Location
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("Location.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_GetID(t *testing.T) {
	tests := []struct {
		name string
		l    *Location
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetID(); got != tt.want {
				t.Errorf("Location.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		l    *Location
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Location.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManufacturer_GetID(t *testing.T) {
	tests := []struct {
		name string
		m    *Manufacturer
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetID(); got != tt.want {
				t.Errorf("Manufacturer.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManufacturer_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		m    *Manufacturer
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manufacturer.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceType_GetID(t *testing.T) {
	tests := []struct {
		name string
		dt   *DeviceType
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.GetID(); got != tt.want {
				t.Errorf("DeviceType.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceType_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		dt   *DeviceType
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dt.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeviceType.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceRole_GetID(t *testing.T) {
	tests := []struct {
		name string
		dr   *DeviceRole
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dr.GetID(); got != tt.want {
				t.Errorf("DeviceRole.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceRole_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		dr   *DeviceRole
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dr.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeviceRole.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_GetID(t *testing.T) {
	tests := []struct {
		name string
		d    *Device
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetID(); got != tt.want {
				t.Errorf("Device.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		d    *Device
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Device.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_GetID(t *testing.T) {
	tests := []struct {
		name string
		i    *Interface
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetID(); got != tt.want {
				t.Errorf("Interface.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterface_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		i    *Interface
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interface.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualDeviceContext_GetID(t *testing.T) {
	tests := []struct {
		name string
		vdc  *VirtualDeviceContext
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vdc.GetID(); got != tt.want {
				t.Errorf("VirtualDeviceContext.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVirtualDeviceContext_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		vdc  *VirtualDeviceContext
		want *NetboxObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vdc.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VirtualDeviceContext.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
