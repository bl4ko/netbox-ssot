# Netbox-SSOT

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

> [!WARNING]
> **This project is still under heavy development, use with caution.**
> Works with `netbox>=3.7.x`.

## Configuration

Netbox-ssot is configured via a single yaml file.
The configuration file is divided into three sections:

- [`logger`](#logger): Logger configuration
- [`netbox`](#netbox): Netbox configuration
- [`source`](#source): Array of configuration for each data source

Example configuration can be found [here](#example-config).

### Logger

| Parameter      | Description                                            | Type | Possible values | Default | Required |
| -------------- | ------------------------------------------------------ | ---- | --------------- | ------- | -------- |
| `logger.level` | Log level                                              | int  | 0-3             | 1       | Yes      |
| `logger.dest`  | Log output filename. Default `""` representing stdout. | str  | Any valid path  | ""      | No       |

### Netbox

| Parameter               | Description                                                                                                                                   | Type     | Possible values | Default       | Required |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- | -------- | --------------- | ------------- | -------- |
| `netbox.apiToken`       | Netbox [API token](https://demo.netbox.dev/static/docs/rest-api/authentication/).                                                             | str      | Any valid token | ""            | Yes      |
| `netbox.hostname`       | Hostname of your netbox instance (e.g `netbox.example.com`).                                                                                  | str      | Valid hostname  | ""            | Yes      |
| `netbox.port`           | Port of your netbox instance.                                                                                                                 | int      | 0-65536         | 443           | No       |
| `netbox.httpScheme`     | HTTP scheme of your netbox instance.                                                                                                          | str      | [http, https]   | https         | No       |
| `netbox.validateCert`   | Validate the TLS certificate of your netbox instance.                                                                                         | bool     | [true, false]   | false         | No       |
| `netbox.timeout`        | Max timeout for api call of your netbox instance.                                                                                             | int      | >=0             | 30            | No       |
| `netbox.removeOrphans`  | Automatically remove all objects tagged with **netbox-ssot** which, were not found on the sources, during this iteration.                     | bool     | [true, false]   | true          | No       |
| `netbox.tag`            | Tag to be applied to all objects managed by netbox-ssot.                                                                                      | string   | any             | "netbox-ssot" | No       |
| `netbox.tagColor`       | TagColor for the netbox-ssot tag.                                                                                                             | string   | any             | "07426b"      | No       |
| `netbox.sourcePriority` | Array of source names in order of priority. If an object (e.g. Vlan) is found in multiple sources, the first source in the list will be used. | []string | any             | []            | No       |

### Source

| Parameter                       | Description                                                                                                        | Source Type     | Type     | Possible values       | Default    | Required |
| ------------------------------- | ------------------------------------------------------------------------------------------------------------------ | --------------- | -------- | --------------------- | ---------- | -------- |
| `source.name`                   | Name of the data source.                                                                                           | all             | str      | any                   | ""         | Yes      |
| `source.type`                   | Type of the data source.                                                                                           | all             | str      | [ovirt, vmware, dnac] | ""         | Yes      |
| `source.httpScheme`             | Http scheme for the source                                                                                         | all             | str      | [ http,https]         | https      | no       |
| `source.hostname`               | Hostname of the data source.                                                                                       | all             | str      | any                   | ""         | Yes      |
| `source.port`                   | Port of the data source.                                                                                           | all             | int      | 0-65536               | 443        | No       |
| `source.username`               | Username of the data source account.                                                                               | all             | str      | any                   | ""         | Yes      |
| `source.password`               | Password of the data source account.                                                                               | all             | str      | any                   | ""         | Yes      |
| `source.validateCert`           | Enforce TLS certificate validation.                                                                                | all             | bool     | [true, false]         | false      | No       |
| `source.tagColor`               | TagColor for the source tag.                                                                                       | all             | string   | any                   | Predefined | No       |
| `source.hostSiteRelations`      | Regex relations in format `regex = siteName`, that map each host that satisfies regex to site.                     | [vmware, ovirt] | []string | any                   | []         | No       |
| `source.clusterSiteRelations`   | Regex relations in format `regex = siteName`, that map each cluster that satisfies regex to site.                  | [vmware, ovirt] | []string | any                   | []         | No       |
| `source.clusterTenantRelations` | Regex relations in format `regex = tenantName`, that map each cluster that satisfies regex to tenant.              | [vmware, ovirt] | []string | any                   | []         | No       |
| `source.hostTenantRelations`    | Regex relations in format `regex = tenantName`, that map each host that satisfies regex to tenant.                 | all             | []string | any                   | []         | No       |
| `source.vmTenantRelations`      | Regex relations in format `regex = tenantName`, that map each vm that satisfies regex to tenant.                   | [vmware, ovirt] | []string | any                   | []         | No       |
| `source.vlanGroupRelations`     | Regex relations in format `regex = vlanGroup`, that map each vlan that satisfies regex to vlanGroup.               | all             | []string | any                   | []         | No       |
| `source.vlanTenantRelations`    | Regex relations in format `regex = tenantName`, that map each vlan that satisfies regex to tenant.                 | all             | []string | any                   | []         | No       |
| `source.customFieldMappings`    | Mappings of format `customFieldName = option`. Currently, supported options are `contact`, `owner`, `description`. | [ vmware ]      | []string | any                   | []         | No       |

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
  sourcePriority: ["olvm", "prodvmware", "prodprox", "dnacenter", "testvmware"] # Not required, but recommended

source:
  - name: olvm
    type: ovirt
    hostname: ovirt.example.com
    port: 443
    username: "admin"
    password: "topsecret"

  - name: prodprox
    type: proxmox
    username: svc@pve
    password: changeme
    hostname: 192.168.1.254
    port: 8006
    validateCert: false
    clusterSiteRelations:
     - .* = Site

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

  - name: testvmare
    type: vmware
    hostname: vcenter-test.example.com
    username: user
    password: passw0rd
    customFieldMappings: # Here we define map of our custom field names, to 3 option [email, owner, description]
      - Email = email
      - Maintainer = owner
      - Notes = description


  - name: dnacenter
    type: dnac
    hostname: dnac.example.com
    username: user
    password: "pa$$w0rd"
    vlanTenantRelations:
      - .* = MyTenant
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

Apply [cronjob](./k8s/cronjob.yaml) with custom settings:

```yaml
kubectl apply -f cronjob.yaml
```
