package objects

import "testing"

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
			want: "Cluster{Name: Test cluster, Type: ovirt}",
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
			want: "VM{ID: 1, Name: Test vm}",
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
