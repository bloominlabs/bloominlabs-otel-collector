// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package billingscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper/internal/metadata"
)

// Config relating to Billing Metric Scraper.
type Config struct {
	// MetricsBuilderConfig allows customizing scraped metrics/attributes representation.
	metadata.MetricsBuilderConfig `mapstructure:",squash"`
	internal.ScraperConfig
}
