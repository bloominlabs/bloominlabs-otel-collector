// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package digitaloceanreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver"

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/accountscraper"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/dropletscraper"
)

// This file implements Factory for Digitalocean receiver.
var (
	scraperFactories = map[string]internal.ScraperFactory{
		accountscraper.TypeStr: &accountscraper.Factory{},
		billingscraper.TypeStr: &billingscraper.Factory{},
		dropletscraper.TypeStr: &dropletscraper.Factory{},
	}
)

// NewFactory creates a new factory for host metrics receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}

func getScraperFactory(key string) (internal.ScraperFactory, bool) {
	if factory, ok := scraperFactories[key]; ok {
		return factory, true
	}

	return nil, false
}

// createDefaultConfig creates the default configuration for receiver.
func createDefaultConfig() component.Config {
	return &Config{ScraperControllerSettings: scraperhelper.NewDefaultScraperControllerSettings(metadata.Type)}
}

// createMetricsReceiver creates a metrics receiver based on provided config.
func createMetricsReceiver(
	ctx context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	oCfg := cfg.(*Config)

	addScraperOptions, err := createAddScraperOptions(ctx, set, oCfg, scraperFactories)
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(
		&oCfg.ScraperControllerSettings,
		set,
		consumer,
		addScraperOptions...,
	)
}

func createAddScraperOptions(
	ctx context.Context,
	set receiver.CreateSettings,
	config *Config,
	factories map[string]internal.ScraperFactory,
) ([]scraperhelper.ScraperControllerOption, error) {
	scraperControllerOptions := make([]scraperhelper.ScraperControllerOption, 0, len(config.Scrapers))

	for key, cfg := range config.Scrapers {
		digitaloceanScraper, ok, err := createDigitaloceanScraper(ctx, set, key, cfg, factories)
		if err != nil {
			return nil, fmt.Errorf("failed to create scraper for key %q: %w", key, err)
		}

		if ok {
			scraperControllerOptions = append(scraperControllerOptions, scraperhelper.AddScraper(digitaloceanScraper))
			continue
		}

		return nil, fmt.Errorf("host metrics scraper factory not found for key: %q", key)
	}

	return scraperControllerOptions, nil
}

func createDigitaloceanScraper(ctx context.Context, set receiver.CreateSettings, key string, cfg internal.Config, factories map[string]internal.ScraperFactory) (scraper scraperhelper.Scraper, ok bool, err error) {
	factory := factories[key]
	if factory == nil {
		ok = false
		return
	}

	ok = true
	scraper, err = factory.CreateMetricsScraper(ctx, set, cfg)
	return
}

type environment interface {
	Lookup(k string) (string, bool)
	Set(k, v string) error
}

type osEnv struct{}

var _ environment = (*osEnv)(nil)

func (e *osEnv) Set(k, v string) error {
	return os.Setenv(k, v)
}

func (e *osEnv) Lookup(k string) (string, bool) {
	return os.LookupEnv(k)
}
