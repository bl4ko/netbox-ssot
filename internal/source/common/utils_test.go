package common

import (
	"context"
	"testing"

	"github.com/bl4ko/netbox-ssot/internal/constants"
	"github.com/bl4ko/netbox-ssot/internal/netbox/inventory"
	"github.com/bl4ko/netbox-ssot/internal/netbox/objects"
	"github.com/bl4ko/netbox-ssot/internal/netbox/service"
)

func setupMockServer(t *testing.T) {
	t.Helper()
	mockServer := service.CreateMockServer()
	t.Cleanup(mockServer.Close)
	service.MockNetboxClient.BaseURL = mockServer.URL
}

func testCtx() context.Context {
	return context.WithValue(context.Background(), constants.CtxSourceKey, "test")
}

// --- Nil relations tests ---

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

// --- No match tests ---

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

func TestMatchHostToTenant_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Tenant1",
	}
	tenant, err := MatchHostToTenant(context.Background(), nil, "dev-host", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant, got %v", tenant)
	}
}

func TestMatchHostToRole_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Role1",
	}
	role, err := MatchHostToRole(context.Background(), nil, "dev-host", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != nil {
		t.Errorf("expected nil role, got %v", role)
	}
}

func TestMatchVMToTenant_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Tenant1",
	}
	tenant, err := MatchVMToTenant(context.Background(), nil, "dev-vm", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant, got %v", tenant)
	}
}

func TestMatchVMToRole_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Role1",
	}
	role, err := MatchVMToRole(context.Background(), nil, "dev-vm", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role != nil {
		t.Errorf("expected nil role, got %v", role)
	}
}

func TestMatchVlanToTenant_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Tenant1",
	}
	tenant, err := MatchVlanToTenant(context.Background(), nil, "dev-vlan", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant != nil {
		t.Errorf("expected nil tenant, got %v", tenant)
	}
}

func TestMatchVlanToSite_NoMatch(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^prod-.*$": "Site1",
	}
	site, err := MatchVlanToSite(context.Background(), nil, "dev-vlan", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site != nil {
		t.Errorf("expected nil site, got %v", site)
	}
}

// --- Invalid regex tests ---

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

func TestMatchClusterToSite_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Site1",
	}
	_, err := MatchClusterToSite(context.Background(), nil, "test-cluster", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchHostToSite_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Site1",
	}
	_, err := MatchHostToSite(context.Background(), nil, "test-host", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchHostToTenant_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Tenant1",
	}
	_, err := MatchHostToTenant(context.Background(), nil, "test-host", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchHostToRole_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Role1",
	}
	_, err := MatchHostToRole(context.Background(), nil, "test-host", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchVMToTenant_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Tenant1",
	}
	_, err := MatchVMToTenant(context.Background(), nil, "test-vm", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchVMToRole_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Role1",
	}
	_, err := MatchVMToRole(context.Background(), nil, "test-vm", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchVlanToTenant_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Tenant1",
	}
	_, err := MatchVlanToTenant(context.Background(), nil, "test-vlan", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchVlanToSite_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "Site1",
	}
	_, err := MatchVlanToSite(context.Background(), nil, "test-vlan", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestMatchIPToVRF_InvalidRegex(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"[invalid": "vrf1",
	}
	_, err := MatchIPToVRF(context.Background(), nil, "10.0.0.1", relations)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

// --- Match with existing object (found in inventory) ---

func TestMatchClusterToTenant_MatchExisting(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^test-.*$": "existing_tenant1",
	}
	tenant, err := MatchClusterToTenant(testCtx(), nbi, "test-cluster", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
	if tenant.Name != "existing_tenant1" {
		t.Errorf("expected tenant name 'existing_tenant1', got %q", tenant.Name)
	}
}

func TestMatchClusterToSite_MatchExisting(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^test-.*$": "existing_site1",
	}
	site, err := MatchClusterToSite(testCtx(), nbi, "test-cluster", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site == nil {
		t.Fatal("expected non-nil site")
	}
	if site.Name != "existing_site1" {
		t.Errorf("expected site name 'existing_site1', got %q", site.Name)
	}
}

func TestMatchVlanToTenant_MatchExisting(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vlan-.*$": "existing_tenant2",
	}
	tenant, err := MatchVlanToTenant(testCtx(), nbi, "vlan-100", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
	if tenant.Name != "existing_tenant2" {
		t.Errorf("expected tenant name 'existing_tenant2', got %q", tenant.Name)
	}
}

func TestMatchVlanToSite_MatchExisting(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vlan-.*$": "existing_site2",
	}
	site, err := MatchVlanToSite(testCtx(), nbi, "vlan-100", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site == nil {
		t.Fatal("expected non-nil site")
	}
	if site.Name != "existing_site2" {
		t.Errorf("expected site name 'existing_site2', got %q", site.Name)
	}
}

func TestMatchVMToTenant_MatchExisting(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vm-.*$": "existing_tenant1",
	}
	tenant, err := MatchVMToTenant(testCtx(), nbi, "vm-web01", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
	if tenant.Name != "existing_tenant1" {
		t.Errorf("expected tenant name 'existing_tenant1', got %q", tenant.Name)
	}
}

// --- Match with new object (not in inventory, triggers Create) ---

func TestMatchClusterToTenant_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^prod-.*$": "brand-new-tenant",
	}
	tenant, err := MatchClusterToTenant(testCtx(), nbi, "prod-cluster", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
}

func TestMatchClusterToSite_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^prod-.*$": "brand-new-site",
	}
	site, err := MatchClusterToSite(testCtx(), nbi, "prod-cluster", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site == nil {
		t.Fatal("expected non-nil site")
	}
}

func TestMatchHostToSite_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^host-.*$": "new-dc-site",
	}
	site, err := MatchHostToSite(testCtx(), nbi, "host-01", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site == nil {
		t.Fatal("expected non-nil site")
	}
}

func TestMatchHostToTenant_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^host-.*$": "new-host-tenant",
	}
	tenant, err := MatchHostToTenant(testCtx(), nbi, "host-01", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
}

func TestMatchHostToRole_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^host-.*$": "new-server-role",
	}
	role, err := MatchHostToRole(testCtx(), nbi, "host-01", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role == nil {
		t.Fatal("expected non-nil role")
	}
}

func TestMatchVMToTenant_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vm-new-.*$": "brand-new-vm-tenant",
	}
	tenant, err := MatchVMToTenant(testCtx(), nbi, "vm-new-01", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
}

func TestMatchVMToRole_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vm-.*$": "new-vm-role",
	}
	role, err := MatchVMToRole(testCtx(), nbi, "vm-app01", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if role == nil {
		t.Fatal("expected non-nil role")
	}
}

func TestMatchVlanToSite_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vlan-new-.*$": "brand-new-vlan-site",
	}
	site, err := MatchVlanToSite(testCtx(), nbi, "vlan-new-100", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if site == nil {
		t.Fatal("expected non-nil site")
	}
}

func TestMatchVlanToTenant_MatchNew(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^vlan-new-.*$": "brand-new-vlan-tenant",
	}
	tenant, err := MatchVlanToTenant(testCtx(), nbi, "vlan-new-100", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tenant == nil {
		t.Fatal("expected non-nil tenant")
	}
}

// --- MatchHostToSite with no match returns default site ---

func TestMatchHostToSite_NoMatchReturnsDefault(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^prod-.*$": "prod-site",
	}
	site, err := MatchHostToSite(testCtx(), nbi, "dev-host", relations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No match and no default site in inventory → returns nil
	if site != nil {
		t.Errorf("expected nil site when no match and no default site, got %v", site)
	}
}

// --- MatchIPToVRF with CIDR stripping ---

func TestMatchIPToVRF_CIDRStripped(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^10\\.0\\..*$": "internal-vrf",
	}
	// VRF matches but doesn't exist in inventory → error
	_, err := MatchIPToVRF(context.Background(), inventory.MockInventory, "10.0.0.1/24", relations)
	if err == nil {
		t.Fatal("expected error for VRF not found in inventory, got nil")
	}
}

func TestMatchIPToVRF_PlainIP(t *testing.T) {
	t.Parallel()
	relations := map[string]string{
		"^10\\.0\\..*$": "internal-vrf",
	}
	_, err := MatchIPToVRF(context.Background(), inventory.MockInventory, "10.0.0.1", relations)
	if err == nil {
		t.Fatal("expected error for VRF not found in inventory, got nil")
	}
}

func TestMatchIPToVRF_VRFNotInInventory(t *testing.T) {
	setupMockServer(t)
	nbi := inventory.MockInventory
	relations := map[string]string{
		"^172\\.16\\..*$": "nonexistent-vrf",
	}
	_, err := MatchIPToVRF(testCtx(), nbi, "172.16.0.1/24", relations)
	if err == nil {
		t.Fatal("expected error for VRF not found in inventory")
	}

	_ = &objects.VRF{
		NetboxObject: objects.NetboxObject{ID: 1},
		Name:         "test-vrf",
	}
}
