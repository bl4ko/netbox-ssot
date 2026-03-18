# CLAUDE.md - Netbox-SSOT

## Project Overview

Go microservice that syncs NetBox (IPAM/DCIM) with external data sources. Runs as a Kubernetes CronJob, syncing all configured sources in parallel via goroutines. Supports 8 source types: oVirt, VMware, DNAC, Proxmox, Palo Alto, Fortigate, FMC, iOS-XE.

## Build & Run

```bash
# Run tests (CI uses -race -covermode=atomic)
go test -race ./...

# Run linter (v1.64.5, 120-char line limit, 40+ linters)
golangci-lint run

# Run with config
go run cmd/netbox-ssot/main.go -config config.yaml

# Build multi-arch Docker image
make build_and_push
```

## Project Structure

```
cmd/netbox-ssot/main.go          # Entry point
internal/
├── constants/                    # Source types, API paths, content types, custom fields, device roles
├── logger/                       # Context-aware logger (source name from ctx)
├── parser/                       # YAML config parsing with validation (regex, subnets, hex colors)
├── utils/                        # Slugify, HTTP client, JSON diff, networking, regex matching
├── netbox/
│   ├── objects/                  # 40+ Netbox data models (dcim, ipam, virtualization, tenancy, wireless, extras)
│   ├── service/                  # HTTP API client (Create, Patch, GetAll, Delete with pagination)
│   ├── inventory/                # Thread-safe singleton: indexed in-memory cache of all Netbox state
│   └── mapper/                   # Reflection-based type↔API-path mapping for generic CRUD
├── source/
│   ├── source.go                 # Factory: NewSource() routes config.Type to implementation
│   ├── common/                   # Source interface + shared relationship matching helpers
│   ├── ovirt/                    # oVirt (SDK: go-ovirt)
│   ├── vmware/                   # VMware vSphere (SDK: govmomi + REST for tags/custom fields)
│   ├── dnac/                     # Cisco DNA Center (SDK: dnacenter-go-sdk)
│   ├── proxmox/                  # Proxmox VE (SDK: go-proxmox)
│   ├── paloalto/                 # Palo Alto (SDK: pango)
│   ├── fortigate/                # FortiGate (custom REST client, bearer token)
│   ├── fmc/                      # Cisco FMC (custom REST client, domain-scoped)
│   └── ios-xe/                   # Cisco IOS-XE (NETCONF/SSH via scrapligo, XML parsing)
k8s/                              # Kubernetes CronJob manifests
testdata/                         # Test fixtures
```

## Architecture

### Execution Flow

1. Parse YAML config (`-config` flag, default `config.yaml`)
2. Initialize logger
3. Create `NetboxInventory` singleton → fetches all existing Netbox objects into indexed in-memory maps
4. Initialize inventory with predefined items (custom fields, tags, device roles)
5. Launch all sources in parallel goroutines (`sync.WaitGroup`)
6. Each source: `Init()` (fetch from external API) → `Sync()` (diff & push to Netbox)
7. If **all** sources succeed → clean up orphans. If any fail → skip deletion, exit code 1

### Source Interface

All sources implement `common.Source`:
```go
type Source interface {
    Init() error                              // Connect & fetch all data from external source
    Sync(*inventory.NetboxInventory) error   // Transform & create/update in Netbox
}
```

**File layout per source**: `{type}.go` (struct + dispatch), `{type}_init.go` (API fetch), `{type}_sync.go` (Netbox sync).

**Factory** (`source.go`): Creates per-source tags (SourceNameTag + SourceTypeTag), routes to implementation.

### NetboxInventory

Thread-safe singleton with per-object-type `sync.Mutex` locks (23+). Objects indexed in multi-dimensional maps for O(1) lookups:
- `devicesIndexByNameAndSiteID[name][siteID]`
- `interfacesIndexByDeviceIDAndName[deviceID][name]`
- `vmsIndexByNameAndClusterID[name][clusterID]`
- `ipAddressesIndex[ifaceType][ifaceName][parentName][address]`
- `vlansIndexByVlanGroupIDAndVID[groupID][vid]`

**Add pattern**: Lock → check index → if exists: diff & PATCH changed fields → if not: POST create → update index → unlock.

**Diff strategy** (`utils.JSONDiffMapExceptID`): Reflection-based field comparison. Source priority (lower int = higher priority) determines which source wins conflicts. Non-ARP IPs always override ARP entries.

### Orphan Management

Objects tagged with `netbox-ssot` are tracked. After sync, unvisited objects are orphans.

- **Hard delete** (`RemoveOrphans: true`): Direct API deletion in dependency order (VLANs before Devices before Clusters)
- **Soft delete** (`RemoveOrphansAfterDays: N`): Tags with `netbox-ssot-orphan` + sets `orphan_last_seen` custom field. Hard-deletes after N days
- Orphans only removed when **all** sources succeed (fail-safe)

### Netbox API Client

- Paginated fetch: 250 items/page, offset-based
- Bulk delete: batches of 50
- Auth: `Authorization: Token <API_TOKEN>`
- Configurable timeout (default 15s), optional TLS cert validation, custom CA support
- Requires Netbox >= 4.2.0

## Config Structure

```yaml
logger:
  level: 0-3           # DEBUG=0, INFO=1, WARNING=2, ERROR=3
  dest: ""              # filepath or empty for stdout

netbox:
  apiToken: ""
  hostname: ""
  port: 443
  httpScheme: https
  validateCert: true
  timeout: 30
  tag: netbox-ssot
  tagColor: "00add8"
  removeOrphans: false
  removeOrphansAfterDays: 0   # 0 = disabled
  sourcePriority: [source1, source2]  # lower index = higher priority
  caFile: ""

sources:
  - name: ""
    type: ovirt|vmware|dnac|proxmox|paloalto|fortigate|fmc|ios-xe
    hostname: ""
    port: 443
    username: ""
    password: ""
    apiToken: ""
    httpScheme: https
    validateCert: true
    caFile: ""
    tag: ""
    tagColor: ""
    ignoredSubnets: []
    permittedSubnets: []
    interfaceFilter: ""           # regex
    collectArpData: false
    ignoreAssetTags: false
    ignoreSerialNumbers: false
    ignoreVMTemplates: false
    assignDomainName: ""
    # Regex relations (map regex patterns to Netbox object names):
    hostSiteRelations: []
    hostTenantRelations: []
    hostRoleRelations: []
    clusterSiteRelations: []
    clusterTenantRelations: []
    vmTenantRelations: []
    vmRoleRelations: []
    vlanGroupRelations: []
    vlanGroupSiteRelations: []
    vlanTenantRelations: []
    vlanSiteRelations: []
    wlanTenantRelations: []
    datacenterClusterGroupRelations: []
    customFieldMappings: []
```

## Custom Fields & Tags

**SSOT-managed custom fields**: `source`, `source_id`, `orphan_last_seen`, `host_cpu_cores`, `host_memory`, `uuid`, `arp_entry`

**Tags**: `netbox-ssot` (marks managed objects), `netbox-ssot-orphan` (orphaned), per-source name tag, per-source type tag

**Device roles**: Server, VM, VM Template, Container, Switch, Firewall (auto-created with predefined colors)

## Code Conventions

- Line length: 120 chars (golangci-lint + golines)
- Table-driven tests with `t.Run` subtests, no mock framework (concrete objects)
- Context propagation via `constants.CtxSourceKey` for source-aware logging
- Explicit error returns with `fmt.Errorf("context: %s", err)` wrapping, no panics
- `sync.WaitGroup` for parallel source syncing, per-type mutex in inventory
- Self-documenting code preferred over comments
- Indent: 2 spaces for YAML/JSON, tabs for Go (gofmt)
- Max field lengths enforced (e.g., `MaxDeviceNameLength=64`, `MaxSerialNumberLength=50`)

## Adding a New Source

1. Create `internal/source/{type}/` with `{type}.go`, `{type}_init.go`, `{type}_sync.go`
2. Implement `common.Source` interface (Init + Sync)
3. Embed `common.Config` struct
4. Add constant in `internal/constants/constants.go` with default tag color
5. Add case to factory switch in `internal/source/source.go`
6. Init: fetch all data into local maps. Sync: use `inventory.Add*()` methods with relationship matching helpers from `common/utils.go`
7. Always apply source tags to every created object

## Tech Stack

- **Go 1.24.0**, CGO_ENABLED=0 (static binary, trimpath, stripped)
- **Docker**: multi-arch (amd64, arm64, arm/v7, arm/v6, 386), Alpine 3.21.3 base, non-root (UID 10001)
- **CI**: GitHub Actions — golangci-lint, `go test -race`, Codecov, Trivy (CRITICAL/HIGH), runs only on Go file changes
- **CD**: Semantic Release → GoReleaser → `ghcr.io/src-doo/netbox-ssot`
- **K8s**: CronJob (every 20min), non-root, read-only fs, seccomp, resource limits 100-200m CPU / 128-256Mi mem
- **Pre-commit**: YAML/JSON validation, golangci-lint, golines, go-unit-tests, go-mod-tidy, gitleaks, trivy, osv-scanner

## Deployment

```bash
kubectl create secret generic netbox-ssot-secret --from-file=config.yaml
kubectl apply -f k8s/cronjob.yaml
# With self-signed certs: k8s/cronjob_with_cert.yaml
```
