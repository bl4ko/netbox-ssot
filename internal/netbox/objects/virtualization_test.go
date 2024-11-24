package objects

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClusterType_String(t *testing.T) {
	tests := []struct {
		name string
		ct   ClusterType
		want string
	}{
		{
			name: "Test cluster type correct string",
			ct: ClusterType{
				Name: "Test cluster type",
			},
			want: "ClusterType{Name: Test cluster type}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("ClusterType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_String(t *testing.T) {
	tests := []struct {
		name string
		c    Cluster
		want string
	}{
		{
			name: "Test cluster correct string",
			c: Cluster{
				Name: "Test cluster",
				Type: &ClusterType{Name: "ovirt"},
			},
			want: fmt.Sprintf("Cluster{Name: Test cluster, Type: %s}", &ClusterType{Name: "ovirt"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Cluster.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVM_String(t *testing.T) {
	tests := []struct {
		name string
		vm   VM
		want string
	}{
		{
			name: "Test vm correct string",
			vm: VM{
				NetboxObject: NetboxObject{
					ID: 1,
				},
				Name: "Test vm",
			},
			want: "VM{Name: Test vm, Cluster: <nil>}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vm.String(); got != tt.want {
				t.Errorf("VM.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVMInterface_String(t *testing.T) {
	tests := []struct {
		name string
		vmi  VMInterface
		want string
	}{
		{
			name: "Test vm interface correct string",
			vmi: VMInterface{
				Name: "Test vm interface",
				VM:   &VM{NetboxObject: NetboxObject{ID: 1}, Name: "Test vm"},
			},
			want: "VMInterface{Name: Test vm interface, VM: Test vm}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vmi.String(); got != tt.want {
				t.Errorf("VMInterface.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterGroup_String(t *testing.T) {
	tests := []struct {
		name string
		cg   ClusterGroup
		want string
	}{
		{
			name: "String representation of cluster group",
			cg: ClusterGroup{
				Name: "Test cluster group",
			},
			want: "ClusterGroup{Name: Test cluster group}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cg.String(); got != tt.want {
				t.Errorf("ClusterGroup.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterGroup_GetID(t *testing.T) {
	tests := []struct {
		name string
		cg   *ClusterGroup
		want int
	}{
		{
			name: "Test cluster group get id",
			cg: &ClusterGroup{
				NetboxObject: NetboxObject{
					ID: 1,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cg.GetID(); got != tt.want {
				t.Errorf("ClusterGroup.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterGroup_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		cg   *ClusterGroup
		want *NetboxObject
	}{
		{
			name: "Test cluster group get netbox object",
			cg: &ClusterGroup{
				NetboxObject: NetboxObject{
					ID: 1,
					CustomFields: map[string]interface{}{
						"x": "y",
					},
				},
			},
			want: &NetboxObject{
				ID: 1,
				CustomFields: map[string]interface{}{
					"x": "y",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cg.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterGroup.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterType_GetID(t *testing.T) {
	tests := []struct {
		name string
		ct   *ClusterType
		want int
	}{
		{
			name: "Test cluster type get id",
			ct: &ClusterType{
				NetboxObject: NetboxObject{
					ID: 1,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.GetID(); got != tt.want {
				t.Errorf("ClusterType.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterType_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		ct   *ClusterType
		want *NetboxObject
	}{
		{
			name: "Test cluster type get netbox object",
			ct: &ClusterType{
				NetboxObject: NetboxObject{
					ID: 1,
					CustomFields: map[string]interface{}{
						"x": "y",
					},
				},
			},
			want: &NetboxObject{
				ID: 1,
				CustomFields: map[string]interface{}{
					"x": "y",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClusterType.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_GetID(t *testing.T) {
	tests := []struct {
		name string
		c    *Cluster
		want int
	}{
		{
			name: "Test cluster get id",
			c: &Cluster{
				NetboxObject: NetboxObject{
					ID: 1,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetID(); got != tt.want {
				t.Errorf("Cluster.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		c    *Cluster
		want *NetboxObject
	}{
		{
			name: "Test cluster get netbox object",
			c: &Cluster{
				NetboxObject: NetboxObject{
					ID: 1,
					CustomFields: map[string]interface{}{
						"x": "y",
					},
				},
			},
			want: &NetboxObject{
				ID: 1,
				CustomFields: map[string]interface{}{
					"x": "y",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cluster.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVM_GetID(t *testing.T) {
	tests := []struct {
		name string
		vm   *VM
		want int
	}{
		{
			name: "Test vm get id",
			vm: &VM{
				NetboxObject: NetboxObject{
					ID: 1,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vm.GetID(); got != tt.want {
				t.Errorf("VM.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVM_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		vm   *VM
		want *NetboxObject
	}{
		{
			name: "Test vm get netbox object",
			vm: &VM{
				NetboxObject: NetboxObject{
					ID: 1,
					CustomFields: map[string]interface{}{
						"x": "y",
					},
				},
			},
			want: &NetboxObject{
				ID: 1,
				CustomFields: map[string]interface{}{
					"x": "y",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vm.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VM.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVMInterface_GetID(t *testing.T) {
	tests := []struct {
		name string
		vmi  *VMInterface
		want int
	}{
		{
			name: "Test vmi get id",
			vmi: &VMInterface{
				NetboxObject: NetboxObject{
					ID: 1,
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vmi.GetID(); got != tt.want {
				t.Errorf("VMInterface.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVMInterface_GetNetboxObject(t *testing.T) {
	tests := []struct {
		name string
		vmi  *VMInterface
		want *NetboxObject
	}{
		{
			name: "Test vmi get netbox object",
			vmi: &VMInterface{
				NetboxObject: NetboxObject{
					ID: 1,
					CustomFields: map[string]interface{}{
						"x": "y",
					},
				},
			},
			want: &NetboxObject{
				ID: 1,
				CustomFields: map[string]interface{}{
					"x": "y",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vmi.GetNetboxObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VMInterface.GetNetboxObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
