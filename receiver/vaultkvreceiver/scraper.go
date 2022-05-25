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

package vaultkvreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver"

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scrapererror"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver/internal/metadata"
)

type vaultKVScraper struct {
	logger        *zap.Logger
	config        *Config
	clientFactory vaultKVClientFactory
	mb            *metadata.MetricsBuilder
}

type vaultKVClientFactory interface {
	getClient(c *Config, mount string) (client, error)
}

type defaultClientFactory struct{}

func (d *defaultClientFactory) getClient(c *Config, database string) (client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return newVaultKVClient(client)
}

func newVaultKVScraper(
	settings component.ReceiverCreateSettings,
	config *Config,
	clientFactory vaultKVClientFactory,
) *vaultKVScraper {
	return &vaultKVScraper{
		logger:        settings.Logger,
		config:        config,
		clientFactory: clientFactory,
		mb:            metadata.NewMetricsBuilder(config.Metrics, settings.BuildInfo),
	}
}

// scrape scrapes the metric stats, transforms them and attributes them into a metric slices.
func (p *vaultKVScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	client, err := p.clientFactory.getClient(p.config, "infra/")
	if err != nil {
		p.logger.Error("Failed to initialize vault client", zap.Error(err))
		return pmetric.NewMetrics(), err
	}

	now := pcommon.NewTimestampFromTime(time.Now())

	client.listKeys(ctx)
	keys, err := client.listKeys(ctx)
	if err != nil {
		return pmetric.NewMetrics(), fmt.Errorf("failed to list keys: %w", err)
	}

	var errors scrapererror.ScrapeErrors

	fmt.Println(keys)
	for _, key := range keys {
		p.collectCreatedTime(ctx, now, client, key, errors)
	}

	return p.mb.Emit(), errors.Combine()
}

func (p *vaultKVScraper) collectCreatedTime(
	ctx context.Context,
	now pcommon.Timestamp,
	client client,
	path string,
	errors scrapererror.ScrapeErrors,
) {
	return
	// metadata, err := client.getPathMetadata(ctx, path)
	// if err != nil {
	// 	p.logger.Error("Errors encountered while fetching path metadata", zap.Error(err))
	// 	errors.AddPartial(0, err)
	// }

	// // Metrics can be partially collected (non-nil) even if there were partial errors reported
	// if blocksReadByTableMetrics == nil {
	// 	return
	// }
	// for _, table := range blocksReadByTableMetrics {
	// 	for sourceKey, source := range metadata.MapAttributeSource {
	// 		value, ok := table.stats[sourceKey]
	// 		if !ok {
	// 			// Data isn't present, error was already logged at a lower level
	// 			continue
	// 		}
	// 		i, err := p.parseInt(sourceKey, value)
	// 		if err != nil {
	// 			errors.AddPartial(0, err)
	// 			continue
	// 		}
	// 		p.mb.RecordPostgresqlBlocksReadDataPoint(now, i, table.database, table.table, source)
	// 	}
	// }
}
