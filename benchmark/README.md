# Benchmarks

## Resource usage

### Netbox

Tested on netbox deployment, with netbox using the following resource limits:

| CPU req | CPU limit | MEM req | MEM limit |
| ------- | --------- | ------- | --------- |
| 500m    | 2         | 128Mi   | 2Gi       |

Netbox resource usage (grafana)
- single

| MEM (single) | MEM (parallel) | CPU (single) | CPU (parallel) |
| ------------ | -------------- | ------------ | -------------- |
| 750MiB       |                | 0.355        |                |

- multiple

### Netbox-ssot

Tested on netbox deployment, with netbox using the following resource limits:

| CPU req | CPU limit | MEM req | MEM limit |
| ------- | --------- | ------- | --------- |
| 50m     | 100m      | 50Mi    | 100Mi     |

Netbox resource usage (grafana)

| MEM (single) | MEM (parallel) | CPU (single) | CPU(parallel) |
| ------------ | -------------- | ------------ | ------------- |
| 40MiB        |                | 0.004        |               |


## Init Run

- `v0.1.5` init run vs parallel run
  - around 6000 objects total
  - 4 external sources

| Single goroutine | Parallel |
| ---------------- | -------- |
| 24m 30s          |          |

## Sync Run

- `v0.1.5` sync run (around 6000 objects total)
  - around 5k objects total
  - 4 external sources

| Single goroutine | Parallel |
| ---------------- | -------- |
| 2m 27s           |          |
| 2m 17s           |          |
| 2m 19s           |          |
