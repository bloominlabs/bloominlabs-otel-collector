# VaultKV Receiver

This receiver queries the [Vault KV V2 metadata](https://www.vaultproject.io/api-docs/secret/kv/kv-v2#list-secrets).

Supported pipeline types: `metrics`

> :construction: This receiver is in **BETA**. Configuration fields and metric data model are subject to change.

## Prerequisites

This receiver supports vault API version v1.

The `VAULT_TOKEN` provided to the receiver must have one of the following policies

```hcl
# list all metadata for all mounts

# list metadata for a mount
path "+/metadata*" {
  capabilities = ["read", "list"]
}

# list metadata for all keys under a mount
path "secret/metadata/*" {
  capabilities = ["read", "list"]
}
```

## Configuration

The following settings are optional:

- `addr` (default = `https://localhost:8200`): The endpoint of the vaultserver.

- `mount` (default = `infra/`): The mount to query key metadata for. The receiver will automatically recurse through the mount and find all keys.

- `collection_interval` (default = `10s`): This receiver collects metrics on an interval. This value must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration). Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

### Example Configuration

```yaml
receivers:
  vaultkv:
    addr: https://localhost:8200
    mount: "infra/"
    collection_interval: 10s
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml)
