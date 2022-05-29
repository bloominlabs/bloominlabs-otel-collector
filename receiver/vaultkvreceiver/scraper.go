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
	"time"

	"github.com/hashicorp/vault/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/scrapererror"
	"go.uber.org/zap"

	metricMetadata "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver/internal/metadata"
)

type vaultKVScraper struct {
	logger        *zap.Logger
	config        *Config
	clientFactory vaultKVClientFactory
	mb            *metricMetadata.MetricsBuilder
}

type vaultKVClientFactory interface {
	getClient(c *Config, mount string) (client, error)
}

type defaultClientFactory struct{}

func (d *defaultClientFactory) getClient(c *Config, mount string) (client, error) {
	config := api.DefaultConfig()
	if c.URL != "" {
		config.Address = c.URL
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return newVaultKVClient(client, c.Mount)
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
		mb:            metricMetadata.NewMetricsBuilder(config.Metrics, settings.BuildInfo),
	}
}

// scrape scrapes the metric stats, transforms them and attributes them into a metric slices.
func (p *vaultKVScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	client, err := p.clientFactory.getClient(p.config, p.config.Mount)
	if err != nil {
		p.logger.Error("Failed to initialize vault client", zap.Error(err))
		return pmetric.NewMetrics(), err
	}

	now := pcommon.NewTimestampFromTime(time.Now())

	var errors scrapererror.ScrapeErrors
	secretMetadata, err := client.listSecretMetadata(ctx)
	if err != nil {
		p.logger.Error("failed to list secret metadata", zap.Error(err))
		errors.AddPartial(0, err)
	}
	p.collectCreatedTime(ctx, now, secretMetadata, errors)

	return p.mb.Emit(), errors.Combine()
}

func (p *vaultKVScraper) collectCreatedTime(
	ctx context.Context,
	now pcommon.Timestamp,
	secretMetadata Secrets,
	errors scrapererror.ScrapeErrors,
) {
	// Metrics can be partially collected (non-nil) even if there were partial errors reported
	if secretMetadata == nil {
		return
	}

	for key, metadata := range secretMetadata {
		t := metadata.Type()
		if t == "" {
			p.mb.RecordVaultkvMetadataErrorDataPoint(
				now, 1, key, p.config.Mount,
				metricMetadata.AttributeMetadataErrorTypeMissingType,
			)
			continue
		}
		val, ok := metricMetadata.MapAttributeType[t]
		if ok {
			p.mb.RecordVaultkvCreatedOnDataPoint(
				now,
				metadata.CreatedTime.Unix(),
				key,
				p.config.Mount,
				val,
			)
		} else {
			p.mb.RecordVaultkvMetadataErrorDataPoint(
				now,
				1,
				key,
				p.config.Mount,
				metricMetadata.AttributeMetadataErrorTypeInvalidType,
			)
		}
	}
}
