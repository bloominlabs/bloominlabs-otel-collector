// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package billingscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper"

import (
	"context"

	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/scraper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper/internal/metadata"
)

// This file implements Factory for Processes scraper.

var TypeStr = "certificates"

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
	s := newBillingScraper(ctx, settings, cfg)
	return scraper.NewMetrics(s.scrape)
}
