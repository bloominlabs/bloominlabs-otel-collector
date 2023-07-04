// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package accountscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/accountscraper"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/accountscraper/internal/metadata"
)

// Config relating to Account Metric Scraper.
type Config struct {
	// MetricsBuilderConfig allows customizing scraped metrics/attributes representation.
	metadata.MetricsBuilderConfig `mapstructure:",squash"`
	internal.ScraperConfig
}
