package common

import (
	"context"
	"testing"
)

func TestMatchClusterToTenant_NilRelations(t *testing.T) {
	t.Parallel()
	tenant, err := MatchClusterToTenant(context.Background(), nil, "test-cluster", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant for nil relations, got %v", tenant)
	}
}

func TestMatchClusterToSite_NilRelations(t *testing.T) {
	t.Parallel()
	site, err := MatchClusterToSite(context.Background(), nil, "test-cluster", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site != nil {
		t.Errorf("expected nil site for nil relations, got %v", site)
	}
}

func TestMatchVlanToTenant_NilRelations(t *testing.T) {
	t.Parallel()
	tenant, err := MatchVlanToTenant(context.Background(), nil, "test-vlan", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant for nil relations, got %v", tenant)
	}
}

func TestMatchVlanToSite_NilRelations(t *testing.T) {
	t.Parallel()
	site, err := MatchVlanToSite(context.Background(), nil, "test-vlan", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site != nil {
		t.Errorf("expected nil site for nil relations, got %v", site)
	}
}

func TestMatchHostToSite_NilRelations(t *testing.T) {
	t.Parallel()
	site, err := MatchHostToSite(context.Background(), nil, "test-host", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site != nil {
		t.Errorf("expected nil site for nil relations, got %v", site)
	}
}

func TestMatchHostToTenant_NilRelations(t *testing.T) {
	t.Parallel()
	tenant, err := MatchHostToTenant(context.Background(), nil, "test-host", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant for nil relations, got %v", tenant)
	}
}

func TestMatchHostToRole_NilRelations(t *testing.T) {
	t.Parallel()
	role, err := MatchHostToRole(context.Background(), nil, "test-host", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != nil {
		t.Errorf("expected nil role for nil relations, got %v", role)
	}
}

func TestMatchVMToTenant_NilRelations(t *testing.T) {
	t.Parallel()
	tenant, err := MatchVMToTenant(context.Background(), nil, "test-vm", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant for nil relations, got %v", tenant)
	}
}

func TestMatchVMToRole_NilRelations(t *testing.T) {
	t.Parallel()
	role, err := MatchVMToRole(context.Background(), nil, "test-vm", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != nil {
		t.Errorf("expected nil role for nil relations, got %v", role)
	}
}

func TestMatchIPToVRF_NilRelations(t *testing.T) {
	t.Parallel()
	vrf, err := MatchIPToVRF(context.Background(), nil, "10.0.0.1/24", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vrf != nil {
		t.Errorf("expected nil VRF for nil relations, got %v", vrf)
	}
}

func TestMatchClusterToTenant_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Production",
	}
	tenant, err := MatchClusterToTenant(context.Background(), nil, "dev-cluster", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant for non-matching cluster, got %v", tenant)
	}
}

func TestMatchClusterToSite_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Datacenter1",
	}
	site, err := MatchClusterToSite(context.Background(), nil, "dev-cluster", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site != nil {
		t.Errorf("expected nil site for non-matching cluster, got %v", site)
	}
}

func TestMatchIPToVRF_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^192\\.168\\..*$": "internal-vrf",
	}
	vrf, err := MatchIPToVRF(context.Background(), nil, "10.0.0.1/24", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vrf != nil {
		t.Errorf("expected nil VRF for non-matching IP, got %v", vrf)
	}
}

func TestMatchClusterToTenant_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Tenant1",
	}
	_, err := MatchClusterToTenant(context.Background(), nil, "test-cluster", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}
