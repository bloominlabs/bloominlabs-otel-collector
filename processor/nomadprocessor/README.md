# Nomad Processor

Supported pipeline types: logs

Enriches log attributes with the 'nomad.allocation.id', 'nomad_allocation_id', 'allocation.id', or 'allocation_id' attributes by querying the nomad API with the provided allocation ID.
**NOTE**: this should only process logs pulled from nomad jobs. Please refer to [config.go](./config.go) for the config spec.
