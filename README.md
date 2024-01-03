# Netbox-SSOT


## Configuration

| Parameter                 | Description                                                      | Type         | Possible values | Default                 | Required |
| ------------------------- | ---------------------------------------------------------------- | ------------ | --------------- | ----------------------- | -------- |
| `logger.level`            | Log level                                                        | int          | 0-3             | 1                       | Yes      |
| `logger.dest`             | Log output filename. Default `""` representing stdout.           | str          | Any valid path  | ""                      | No       |
| `netbox.apiToken`         | apiToken to access netbox                                        | str          | Any valid token | ""                      | Yes      |
| `netbox.hostname`         | Netbox hostname (e.g `netbox.example.com`)                       | str          | Valid hostname  | ""                      | Yes      |
| `netbox.port`             | Netbox port                                                      | int          | 0-65536         | 443                     | No       |
| `netbox.HTTPScheme`       | Netbox API HTTP scheme                                           | str          | http, https     | https                   | No       |
| `netbox.validateCert`     | Validate Netbox's TLS certificate                                | bool         | true, false     | false                   | No       |
| `netbox.clientCert`       | Path to client certificate                                       | str          | any             | ""                      | No       |
| `netbox.clientCertKey`    | Path to client certificate key                                   | str          | any             | ""                      | No       |
| `netbox.timeout`          | API call timeout in seconds                                      | int          | >=0             | 30                      | No       |
| `netbox.maxRetries`       | Number of retries for failed API calls                           | int          | >=0             | 3                       | No       |
| `netbox.tag`              | Tag to be applied to all objects managed by netbox-ssot          | string       | any             | "netbox-ssot"           | No       |
| `netbox.tagColor`         | TagColor for the netbox-ssot tag.                                | string       | any             | "07426b"                | No       |
| `source`                  | Array of data sources. Each data source requires its own config  | []SourceType | SourceType      | []                      | No       |
| `source.name`             | Name of the data source.                                         | str          | any             | ""                      | Yes      |
| `source.type`             | Data source type                                                 | str          | ovirt, vmware   | ""                      | Yes      |
| `source.hostname`         | Hostname of the data source                                      | str          | any             | ""                      | Yes      |
| `source.port`             | Port of the data source                                          | int          | 0-65536         | 443                     | No       |
| `source.username`         | Username of the data source account.                             | str          | any             | ""                      | Yes      |
| `source.password`         | Password of the data source account.                             | str          | any             | ""                      | Yes      |
| `source.validateCert`     | Enforce TLS certificate validation.                              | bool         | true, false     | false                   | No       |
| `source.permittedSubnets` | Array of subnets permitted for the osurce. Format: CIDR notation | []string     | any             | []                      | No       |
| `source.tag`              | Tag to be applied to all objects created by this source.         | string       | any             | "source-" + source.name | No       |
| `source.tagColor`         | TagColor for the source tag.                                     | string       | any             | ovirt: "07426b"         | No       |


Example config

```yaml
logger:
  level: 1 # 0=Debug, 1=Info, 2=Warn, 3=Error
  dest: "" # Leave blank for stdout, or specify a file path

netbox:
  apiToken: "" # Netbox API Token
  hostname: "netbox.example.com" # Netbox FQDN
  port: 443 # Netbox Port
  # Proxy: TODO
  # proxy: "" # Defines a proxy which will be used to connect to Netbox. Proxy setting needs to include ; the schema. Proxy basic auth example: http://user:pass@10.10.1.10:312
  # proxyPort: ""
  validateCert: true # Validate Netbox TLS certificate
  clientCert: "" # Path to client certificate
  clientCertKey: "" # Path to client certificate key
  timeout: 30 # API call timeout in seconds
  maxRetries: 3 # Number of retries for failed API calls
  # TODO:
  # useCaching: false # Enable caching of Netbox data
  # cacheDir: "/tmp" # Path to cache directory

source:  # Array of sources
  - name: "" # Name of the source
    type: "" # Source type: [ovirt,  vmware]
    hostname: "" # Hostname of the source vcenter.example.com
    port: # Port of the source
    username: "" # Username of the source account "admin"
    password: "" # Password of the source account "secretpass"
    validateCert: false # Enforce TLS certificate validation bool
    # Proxy TODO
    # proxyHost: ""
    # proxyPort: ""
    permittedSubnets: # Array of permitted subnets
      - "172.16.0.0/12"
      - "10.0.0.0/8"
      - "192.168.0.0/16"
      - "fd00::/8"
    # filters:
    #   clusterExcludeFilter: # array string filters to exclude clusters
    #     - ".*-template"
    #   hostExcludeFilter: # array string filters to exclude hosts
    #     - ".*-template"
    #   vmExcludeFilter: # array string filters to exclude VMs
    #     - ".*-template"
    # relations:
    #   cluster_site_relation: # array of cluster to site relations
    #     - cluster: ".*"
    #       site: ".*"
    #   host_site_relation:
    #     - host: ".*"
    #       site: ".*"
    #   host_cluster_relation: # array of host to cluster relations
    #     - host: ".*"
    #       cluster: ".*"
    #   vm_host_relation: # array of vm to host relations
    #     - vm: ".*"
    #       host: ".*"
```