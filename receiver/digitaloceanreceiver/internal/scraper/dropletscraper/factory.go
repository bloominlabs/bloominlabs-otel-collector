// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package dropletscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/dropletscraper"

import (
	"context"

	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/scraper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/dropletscraper/internal/metadata"
)

// This file implements Factory for Processes scraper.

const (
	// TypeStr the value of "type" key in configuration.
	TypeStr = "droplet"
)

// Factory is the Factory for scraper.
type Factory struct{}

// CreateDefaultConfig creates the default configuration for the Scraper.
func (f *Factory) CreateDefaultConfig() internal.Config {
	return &Config{
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

// CreateMetricsScraper creates a scraper based on provided config.
func (f *Factory) CreateMetricsScraper(
	ctx context.Context,
	settings receiver.Settings,
	config internal.Config,
) (scraper.Metrics, error) {
	cfg := config.(*Config)
	s := newDropletScraper(ctx, settings, cfg)

	return scraper.NewMetrics(
		s.scrape,
		scraper.WithStart(s.start),
	)
}
