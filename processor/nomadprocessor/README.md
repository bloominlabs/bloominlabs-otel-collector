# Loki Processor

Supported pipeline types: logs

Custom nomad processor that converts systemd journal logs into pdata.Logs that
are consumable by the Loki exporter. **NOTE**: this should only process systemd
data. Please refer to [config.go](./config.go) for the config spec.

`attributes` represents actions that can be applied on resource attributes.
See processor/attributesprocessor/README.md for more details on supported attributes actions.

Examples:

```yaml
processors:
  resource:
    attributes:
      - key: cloud.availability_zone
        value: "zone-1"
        action: upsert
      - key: k8s.cluster.name
        from_attribute: k8s-cluster
        action: insert
      - key: redundant-attribute
        action: delete
```

Refer to [config.yaml](./testdata/config.yaml) for detailed
examples on using the processor.
