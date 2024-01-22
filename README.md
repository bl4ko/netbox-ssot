# Netbox-SSOT

Netbox-ssot is a tool to keep Netbox in sync with external data sources. 

Currently the supported external data sources types are:

- `ovirt`
- `vmware`


> This plugin operates only with objects that have `netbox-ssot` Tag assigned to them.
> So if you want your existing objects to be managed by it, you have to assign them `netbox-ssot`  Tag.


## Configuration

Example config can be found in [example section](#example-config)

### Logger

| Parameter      | Description                                            | Type | Possible values | Default | Required |
| -------------- | ------------------------------------------------------ | ---- | --------------- | ------- | -------- |
| `logger.level` | Log level                                              | int  | 0-3             | 1       | Yes      |
| `logger.dest`  | Log output filename. Default `""` representing stdout. | str  | Any valid path  | ""      | No       |


### Netbox

| Parameter              | Description                                                                                                | Type   | Possible values | Default       | Required |
| ---------------------- | ---------------------------------------------------------------------------------------------------------- | ------ | --------------- | ------------- | -------- |
| `netbox.apiToken`      | apiToken to access netbox                                                                                  | str    | Any valid token | ""            | Yes      |
| `netbox.hostname`      | Netbox hostname (e.g `netbox.example.com`)                                                                 | str    | Valid hostname  | ""            | Yes      |
| `netbox.port`          | Netbox port                                                                                                | int    | 0-65536         | 443           | No       |
| `netbox.HTTPScheme`    | Netbox API HTTP scheme                                                                                     | str    | http, https     | https         | No       |
| `netbox.validateCert`  | Validate Netbox's TLS certificate                                                                          | bool   | true, false     | false         | No       |
| `netbox.timeout`       | Max netbox API call length in seconds                                                                      | int    | >=0             | 30            | No       |
| `netbox.removeOrphans` | Remove all objects tagged with **netbox-ssot** which, were not found on the sources, during this iteration | bool   | true, false     | true          | No       |
| `netbox.tag`           | Tag to be applied to all objects managed by netbox-ssot                                                    | string | any             | "netbox-ssot" | No       |
| `netbox.tagColor`      | TagColor for the netbox-ssot tag.                                                                          | string | any             | "07426b"      | No       |

### Source

| Parameter                       | Description                                                                                          | Type         | Possible values | Default                 | Required |
| ------------------------------- | ---------------------------------------------------------------------------------------------------- | ------------ | --------------- | ----------------------- | -------- |
| `source`                        | Array of data sources. Each data source requires its own config                                      | []SourceType | SourceType      | []                      | No       |
| `source.name`                   | Name of the data source.                                                                             | str          | any             | ""                      | Yes      |
| `source.type`                   | Data source type                                                                                     | str          | ovirt, vmware   | ""                      | Yes      |
| `source.hostname`               | Hostname of the data source                                                                          | str          | any             | ""                      | Yes      |
| `source.port`                   | Port of the data source                                                                              | int          | 0-65536         | 443                     | No       |
| `source.username`               | Username of the data source account.                                                                 | str          | any             | ""                      | Yes      |
| `source.password`               | Password of the data source account.                                                                 | str          | any             | ""                      | Yes      |
| `source.validateCert`           | Enforce TLS certificate validation.                                                                  | bool         | true, false     | false                   | No       |
| `source.permittedSubnets`       | Array of subnets permitted for the osurce. Format: CIDR notation                                     | []string     | any             | []                      | No       |
| `source.tag`                    | Tag to be applied to all objects created by this source.                                             | string       | any             | "source-" + source.name | No       |
| `source.tagColor`               | TagColor for the source tag.                                                                         | string       | any             | ovirt: "07426b"         | No       |
| `source.hostSiteRelations`      | Regex relations in format `regex = siteName`, that map each host that satisifies regex to site       | []string     | any             | []                      | No       |
| `source.clusterSiteRelations`   | Regex relations in format `regex = siteName`, that map each cluster that satisifes regex to site     | []string     | any             | []                      | No       |
| `source.clusterTenantRelations` | Regex relations in format `regex = tenantName`, that map each cluster that satisifes regex to tenant | []string     | any             | []                      | no       |
| `source.hostTenantRelations`    | Regex relations in format `regex = tenantName`, that map each host that satisifes regex to tenant    | []string     | any             | []                      | no       |
| `source.vmTenantRelations`      | Regex relations in format `regex = tenantName`, that map each vm that satisifes regex to tenant      | []string     | any             | []                      | no       |
| `source.vlanGroupRelations`     | Regex relations in format `regex = vlanGroup`, that map each vlan that satisifes regex to vlanGroup  | []string     | any             | []                      | no       |
| `source.vlanTenantRelations`    | Regex relations in format `regex = tenantName`, that map each vlan that satisfies regex to tenant    | []string     | any             | []                      | no       |  |


### Example config

```yaml
logger:
  level: 1 # 0=Debug, 1=Info, 2=Warn, 3=Error
  dest: "" # Leave blank for stdout, or specify a file path

netbox:
  apiToken: "" # Netbox API Token
  hostname: "netbox.example.com" # Netbox FQDN
  port: 443
  timeout: 30 # API call timeout in seconds

source:
  - name: "Test oVirt"
    type: "ovirt"
    hostname: "ovirt.example.com"
    port: 443
    username: "admin" # Username of the source account "admin"
    password: "topsecret" # Password of the source account "secretpass"
```

## Deployment

### Via docker

```bash
docker run -v /path/to/config.yaml:/app/config.yaml ghcr.io/bl4ko/netbox-ssot
```

### Via k8s

Create k8s secret from self defined config.yaml:

```yaml
kubectl create secret generic netbox-ssot-secret --from-file=config.yaml
```

Apply [cronjob](cronjob.yaml) with custom settings:

```yaml
kubectl apply -f cronjob.yaml
```