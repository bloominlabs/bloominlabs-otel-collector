// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package userstatsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver/internal/metadata"
)

const (
	typeStr = "userstats"

	stability = component.StabilityLevelDevelopment
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, stability))
}

func createDefaultConfig() component.Config {
	return &Config{
		Endpoint:                  "digitaloceanspaces.com",
		Region:                    "sfo3",
		AccessKeyID:               "",
		SecretAccessKey:           "",
		Bucket:                    "dev-stratos-server-backups",
		ScraperControllerSettings: scraperhelper.NewDefaultScraperControllerSettings(typeStr),
		MetricsBuilderConfig:      metadata.DefaultMetricsBuilderConfig(),
	}
}

func createMetricsReceiver(
	ctx context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	rCfg := cfg.(*Config)

	ns := newBackupsUtilizationScraper(set, rCfg, &defaultClientFactory{})
	scraper, err := scraperhelper.NewScraper(typeStr, ns.scrape)
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(
		&rCfg.ScraperControllerSettings, set, consumer,
		scraperhelper.AddScraper(scraper),
	)
}