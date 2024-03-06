package objects

import "testing"

func TestClusterType_String(t *testing.T) {
	tests := []struct {
		name string
		ct   ClusterType
		want string
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vmi.String(); got != tt.want {
				t.Errorf("VMInterface.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
