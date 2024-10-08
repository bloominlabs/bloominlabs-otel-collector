// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/confmap"
)

// MetricConfig provides common config for a particular metric.
type MetricConfig struct {
	Enabled bool `mapstructure:"enabled"`

	enabledSetByUser bool
}

func (ms *MetricConfig) Unmarshal(parser *confmap.Conf) error {
	if parser == nil {
		return nil
	}
	err := parser.Unmarshal(ms)
	if err != nil {
		return err
	}
	ms.enabledSetByUser = parser.IsSet("enabled")
	return nil
}

// MetricsConfig provides config for vaultkv metrics.
type MetricsConfig struct {
	VaultkvCreatedOn     MetricConfig `mapstructure:"vaultkv.created_on"`
	VaultkvMetadata      MetricConfig `mapstructure:"vaultkv.metadata"`
	VaultkvMetadataError MetricConfig `mapstructure:"vaultkv.metadata.error"`
}

func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		VaultkvCreatedOn: MetricConfig{
			Enabled: true,
		},
		VaultkvMetadata: MetricConfig{
			Enabled: true,
		},
		VaultkvMetadataError: MetricConfig{
			Enabled: true,
		},
	}
}

// MetricsBuilderConfig is a configuration for vaultkv metrics builder.
type MetricsBuilderConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
}

func DefaultMetricsBuilderConfig() MetricsBuilderConfig {
	return MetricsBuilderConfig{
		Metrics: DefaultMetricsConfig(),
	}
}
