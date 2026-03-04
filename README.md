# Netbox-SSOT

[![Go](https://github.com/src-doo/netbox-ssot/actions/workflows/ci.yml/badge.svg)](https://github.com/src-doo/netbox-ssot/actions/workflows/ci.yml)
![GitHub last commit](https://img.shields.io/github/last-commit/src-doo/netbox-ssot)
![GitHub Tag](https://img.shields.io/github/v/tag/src-doo/netbox-ssot)
![GitHub License](https://img.shields.io/github/license/src-doo/netbox-ssot)

Netbox-ssot is a small but powerful microservice designed to
keep your Netbox instance in sync with external data sources.

It is designed to be run as a cronjob, and will periodically update Netbox
with the latest data from the external sources. It syncs each source in parallel
to speed up the process of syncing.

Currently, the supported external data sources types are:

- [`ovirt`](https://www.ovirt.org/)
- [`vmware`](https://www.vmware.com/products/vcenter.html)
- [`dnac`](https://www.cisco.com/site/us/en/products/networking/catalyst-center/index.html)
- [`proxmox`](https://www.proxmox.com/en/)
- [`paloalto`](https://www.paloaltonetworks.com/network-security/next-generation-firewall)
  - PAN-OS firewall
- [`fortigate`](https://www.fortinet.com/products/next-generation-firewall)
- [`fmc`](https://www.cisco.com/site/us/en/products/security/firewalls/firewall-management-center/index.html)
- [`ios-xe`](https://www.cisco.com/c/en/us/products/ios-nx-os-software/ios-xe/index.html)
  - All devices with ios-xe supporting netconf

## Compatability Matrix

> [!WARNING]
> Since netbox introduces breaking changes in minor releases, netbox-ssot also introduces breaking changes in minor releases.
> See the table below for compatibility between netbox-ssot and netbox.

| Version       | Supported Netbox Version |
| ------------- | ------------------------ |
| v1.9.x        | >= 4.2.0                 |
| v1.0.0-v1.8.x | >=4.0.0, < 4.2.0         |
| v0.x.x        | >=3.7.0, < 4.0.0         |

## Configuration

Netbox-ssot is configured via a single yaml file.
The configuration file is divided into three sections:

- [`logger`](#logger): Logger configuration
- [`netbox`](#netbox): Netbox configuration
- [`source`](#source): Array of configuration for each data source

Example configuration can be found [here](#example-config).

### Logger

| Parameter      | Description                                            | Type       | Possible values                  | Default | Required |
| -------------- | ------------------------------------------------------ | ---------- | -------------------------------- | ------- | -------- |
| `logger.level` | Log level                                              | int/string | [0-3] or [debug,info,warn,error] | 1,info  | No       |
| `logger.dest`  | Log output filename. Default `""` representing stdout. | str        | Any valid path                   | ""      | No       |

### Netbox

| Parameter                       | Description                                                                                                                                                                                                                                                                                                                                       | Type     | Possible values | Default       | Required |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------------- | ------------- | -------- |
| `netbox.apiToken`               | Netbox API token | str      | Any valid token | ""            | Yes      |
| `netbox.hostname`               | Hostname of your netbox instance (e.g `netbox.example.com`).                                                                                                                                                                                                                                                                                      | str      | Valid hostname  | ""            | Yes      |
| `netbox.port`                   | Port of your netbox instance.                                                                                                                                                                                                                                                                                                                     | int      | 0-65536         | 443           | No       |
| `netbox.httpScheme`             | HTTP scheme of your netbox instance.                                                                                                                                                                                                                                                                                                              | str      | [http, https]   | https         | No       |
| `netbox.validateCert`           | Validate the TLS certificate of your netbox instance.                                                                                                                                                                                                                                                                                             | bool     | [true, false]   | false         | No       |
| `netbox.timeout`                | Max timeout for api call of your netbox instance.                                                                                                                                                                                                                                                                                                 | int      | >=0             | 30            | No       |
| `netbox.removeOrphans`          | If set to **true** all objects, marked with netbox-ssot tag that were not found during this iteration are automatically deleted. If set to **false**, objects that were not found are marked with an **Orphan** tag. We can then use **netbox.removeOrphansAfterDays** to remove the orphans after n days that they were not seen on the sources. | bool     | [true, false]   | true          | No       |
| `netbox.removeOrphansAfterDays` | Specifies the number of days to wait before automatically deleting objects marked as Orphan. This setting is only applicable if netbox.removeOrphans is set to false. A value of 5 means objects are deleted in five days after being marked as Orphan and not found since.                                                                       | int      | >0              | MaxInt        | No       |
| `netbox.tag`                    | Tag to be applied to all objects managed by netbox-ssot.                                                                                                                                                                                                                                                                                          | string   | any             | "netbox-ssot" | No       |
| `netbox.tagColor`               | TagColor for the netbox-ssot tag.                                                                                                                                                                                                                                                                                                                 | string   | any             | "07426b"      | No       |
| `netbox.sourcePriority`         | Array of source names in order of priority. If an object (e.g. Vlan) is found in multiple sources, the first source in the list will be used.                                                                                                                                                                                                     | []string | any             | []            | No       |
| `netbox.caFile`                 | Path to a self signed certificate for netbox.                                                                                                                                                                                                                                                                                                     | string   | Valid path      | ""            | No       |

### Source

| Parameter                                | Description                                                                                                              | Source Type                | Type     | Possible values                          | Default    | Required |
|------------------------------------------|--------------------------------------------------------------------------------------------------------------------------|----------------------------| -------- | ---------------------------------------- |------------| -------- |
| `source.name`                            | Name of the data source.                                                                                                 | all                        | str      | any                                      | ""         | Yes      |
| `source.type`                            | Type of the data source.                                                                                                 | all                        | str      | [ovirt, vmware, dnac, proxmox, paloalto] | ""         | Yes      |
| `source.httpScheme`                      | Http scheme for the source                                                                                               | all                        | str      | [ http,https]                            | https      | No       |
| `source.hostname`                        | Hostname of the data source.                                                                                             | all                        | str      | any                                      | ""         | Yes      |
| `source.port`                            | Port of the data source.                                                                                                 | all                        | int      | 0-65536                                  | 443        | No       |
| `source.username`                        | Username of the data source account.                                                                                     | all                        | str      | any                                      | ""         | Yes      |
| `source.password`                        | Password of the data source account.                                                                                     | all                        | str      | any                                      | ""         | Yes      |
| `source.apiToken`                        | API token of the data source account.                                                                                    | [**fortigate**]            | str      | any                                      | ""         | Yes      |
| `source.validateCert`                    | Enforce TLS certificate validation.                                                                                      | all                        | bool     | [true, false]                            | false      | No       |
| `source.tagColor`                        | TagColor for the source tag.                                                                                             | all                        | string   | any                                      | Predefined | No       |
| `source.ignoredSubnets`                  | List of subnets, which will be ignored (e.g. IPs won't be synced).                                                       | all                        | []string | any                                      | []         | No       |
| `source.permittedSubnets`                | List of subnets, which will be permitted (e.g. only IPs in these subnets will be synced).                                | all                        | []string | any                                      | []         | No       |
| `source.interfaceFilter`                 | Regex representation of interface names to be ignored (e.g. `(cali\|vxlan\|flannel\|[a-f0-9]{15})`)                      | all                        | string   | any                                      | []         | No       |
| `source.collectArpData`                  | Collect data from the arp table of the device.                                                                           | [**paloalto**, **ios-xe**] | bool     | [true, false]                            | false      | No       |
| `source.ignoreAssetTags`                 | Don't sync asset tags of devices.                                                                                        | all                        | bool     | [true, false]                            | false      | No       |
| `source.ignoreSerialNumbers`             | Don't sync serial numbers of devices.                                                                                    | all                        | bool     | [true, false]                            | false      | No       |
| `source.ignoreVMTemplates`               | Don't sync vm templates.                                                                                                 | [**vmware**,**Proxmox**]   | bool     | [true, false]                            | false      | No       |
| `source.AssignDomainName`                | Suffix node name with `AssignDomainName`.                                                                                | [**proxmox**]              | str      | any                                      | ""         | No       |
| `source.vlanPrefix`                      | Prefix vlan name with `vlanPrefix`.                                                                                      | [**vmware**]               | str      | any                                      | ""         | No       |
| `source.datacenterClusterGroupRelations` | Regex relations in format `regex = clusterGroupName`, that map each datacenter that satisfies regex to clusterGroupname. | [**vmware**, **ovirt**]    | []string | any                                      | []         | No       |
| `source.hostSiteRelations`               | Regex relations in format `regex = siteName`, that map each host that satisfies regex to site.                           | all                        | []string | any                                      | []         | No       |
| `source.clusterSiteRelations`            | Regex relations in format `regex = siteName`, that map each cluster that satisfies regex to site.                        | all                        | []string | any                                      | []         | No       |
| `source.clusterTenantRelations`          | Regex relations in format `regex = tenantName`, that map each cluster that satisfies regex to tenant.                    | all                        | []string | any                                      | []         | No       |
| `source.hostTenantRelations`             | Regex relations in format `regex = tenantName`, that map each host that satisfies regex to tenant.                       | all                        | []string | any                                      | []         | No       |
| `source.hostRoleRelations`               | Regex relations in format `regex = roleName`, that map each host that satisfies regex to device role.                    | all                        | []string | any                                      | []         | No       |
| `source.hostTenantRelations`             | Regex relations in format `regex = tenantName`, that map each host that satisfies regex to tenant.                       | all                        | []string | any                                      | []         | No       |
| `source.vmTenantRelations`               | Regex relations in format `regex = tenantName`, that map each vm that satisfies regex to tenant.                         | all                        | []string | any                                      | []         | No       |
| `source.vmRoleRelations`                 | Regex relations in format `regex = roleName`, that map each vm that satisfies regex to device role.                      | all                        | []string | any                                      | []         | No       |
| `source.ipVrfRelations`                  | Regex relations in format `regex = vrfName`, that map each ip that satisfies regex to vrf.                               | all                        | []string | any                                      | []         | No       |
| `source.vlanGroupRelations`              | Regex relations in format `regex = vlanGroup`, that map each vlan that satisfies regex to vlanGroup.                     | all                        | []string | any                                      | []         | No       |
| `source.vlanGroupSiteRelations`          | Regex relations in format `regex = vlanGroup`, that map each vlanGroup that satisfies regex to site.                     | all                        | []string | any                                      | []         | No       |
| `source.vlanSiteRelations`               | Regex relations in format `regex = vlan`, that map each vlan that satisfies regex to site.                               | all                        | []string | any                                      | []         | No       |
| `source.wlanTenantRelations`             | Regex relations in format `regex = tenantName`, that map each wlan that satisfies regex to tenant.                       | [dnac]                     | []string | any                                      | []         | No       |
| `source.customFieldMappings`             | Mappings of format `customFieldName = option`. Currently, supported options are `contact`, `owner`, `description`.       | [**vmware**]               | []string | any                                      | []         | No       |
| `source.caFile`                          | Path to a self signed certificate for the source.                                                                        | any                        | string   | Valid path                               | ""         | No       |

### Example config

```yaml
logger:
  level: 1
  dest: ""

netbox:
  apiToken: "el1aof2azu6n50ks5zcenp3..."
  hostname: "netbox.example.com"
  httpScheme: http
  port: 443
  timeout: 30
  sourcePriority: ["olvm", "prodvmware", "prodprox", "dnacenter", "testvmware", "pa-uk", "fmc-lab"] # Not required, but recommended

source:
  - name: olvm
    type: ovirt
    hostname: ovirt.example.com
    port: 443
    username: "admin"
    password: "topsecret"
    interfaceFilter: (cali|vxlan|flannel|docker|[a-f0-9]{15})

  - name: prodvmware
    type: vmware
    hostname: vcenter.example.com
    username: user
    password: "top_secret"
    clusterSiteRelations:
      - .* = ExampleSite
    hostSiteRelations:
      - .*_NYC = New York
      - nyc.* = New York
    customFieldMappings: # Here we define map of our custom field names, to 3 option [email, owner, description]
      - Mail = email
      - Creator = owner
      - Description = description

  - name: prodprox
    type: proxmox
    username: svc@pve
    password: changeme
    hostname: 192.168.1.254
    port: 8006
    validateCert: false
    clusterSiteRelations:
     - .* = Site

  - name: forti
    type: fortigate
    hostname: forti.example.com
    apiToken: "apitokenhere"
    validateCert: False
    hostTenantRelations:
      - .* = MyTenant
    hostSiteRelations:
      - .* = MyTenant
    vlanTenantRelations:
      - .* = MyTenant

  - name: pa-uk
    type: paloalto
    hostname: 192.168.1.52
    username: user
    password: passw0rd
    hostTenantRelations:
      - .* = MyTenant
    hostSiteRelations:
      - .* = MySite
    vlanTenantRelations:
      - .* = MyTenant
    collectArpData: true

  - name: dnacenter
    type: dnac
    hostname: dnac.example.com
    username: user
    password: "pa$$w0rd"
    vlanTenantRelations:
      - .* = MyTenant

  - name: fmc-lab
    type: fmc
    hostname: 172.16.1.30
    username: user
    password: password
    validateCert: False
    hostTenantRelations:
      - .* = MyTenant
    hostSiteRelations:
      - .* = MySite
    vlanTenantRelations:
      - .* = MyTenant

  - name: cs1
    type: ios-xe
    hostname: 10.10.1.1
    username: user
    password: password
    port: 830
    validateCert: False
    hostTenantRelations:
      - .* = MyTenant
    hostSiteRelations:
      - .* = MySite
    vlanTenantRelations:
      - .* = MyTenant
    collectArpData:
      true

```

## Deployment

### Via docker

```bash
docker run -v /path/to/config.yaml:/app/config.yaml ghcr.io/src-doo/netbox-ssot
```

### Via k8s

Create k8s secret from self defined config.yaml:

```yaml
kubectl create secret generic netbox-ssot-secret --from-file=config.yaml
```

Apply [cronjob](./k8s/cronjob.yaml) with custom settings:

```yaml
kubectl apply -f cronjob.yaml
```

#### Using self signed certificate

Create self signed certificate e.g.:

```yaml
kubectl create secret generic netbox-ssot-cert --from-file=sub.pem=./sub.pem
```

Use [cronjob with cert mounted](./k8s/cronjob_with_cert.yaml):

```yaml
kubectl apply -f cronjob_with_cert.yaml
```
