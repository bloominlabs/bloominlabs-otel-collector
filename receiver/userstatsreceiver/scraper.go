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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bloominlabs/baseplate-go/config/spaces"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/scraper/scrapererror"
	"go.uber.org/zap"

	metricMetadata "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver/internal/metadata"
)

type userStatsScraper struct {
	logger        *zap.Logger
	rConfig       *Config
	s3Config      *spaces.DigitalOceanSpacesConfig
	clientFactory userStatsClientFactory
	mb            *metricMetadata.MetricsBuilder
}

type userStatsClientFactory interface {
	getClient(c *s3.Client, bucketName string) (client, error)
}

type defaultClientFactory struct{}

func (d *defaultClientFactory) getClient(c *s3.Client, bucketName string) (client, error) {
	return newBackupsUtilizationClient(c, bucketName)
}

func newBackupsUtilizationScraper(
	settings receiver.Settings,
	config *Config,
	clientFactory userStatsClientFactory,
) *userStatsScraper {
	s3Config := &spaces.DigitalOceanSpacesConfig{
		Region:           config.Region,
		AccessKeyID:      config.AccessKeyID,
		SecretAccessKey:  config.SecretAccessKey,
		InternalEndpoint: config.Endpoint,
	}

	scraper := &userStatsScraper{
		logger:        settings.Logger,
		rConfig:       config,
		s3Config:      s3Config,
		clientFactory: clientFactory,
		mb:            metricMetadata.NewMetricsBuilder(config.MetricsBuilderConfig, settings),
	}

	return scraper
}

// scrape scrapes the metric stats, transforms them and attributes them into a metric slices.
func (p *userStatsScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	p.logger.Info("starting scrape")
	s3Client, err := p.s3Config.GetClient()
	if err != nil {
		return pmetric.NewMetrics(), fmt.Errorf("failed to get s3 client: %w", err)
	}

	client, err := p.clientFactory.getClient(&s3Client, p.rConfig.Bucket)
	if err != nil {
		p.logger.Error("Failed to initialize client", zap.Error(err))
		return pmetric.NewMetrics(), err
	}

	now := pcommon.NewTimestampFromTime(time.Now())

	var errors scrapererror.ScrapeErrors
	p.logger.Info("starting about to list backups")
	backupsMap, err := client.listBackupsByUser(ctx)
	p.logger.Info("stopping list backups")

	if err != nil {
		p.logger.Error("failed to list backups", zap.Error(err))
		return pmetric.NewMetrics(), err
	}

	for userID, sizeByType := range backupsMap {
		for backupType, size := range sizeByType {
			p.mb.RecordBackupsTotalSizeDataPoint(
				now,
				size,
				userID,
				backupType,
			)
		}
	}

	return p.mb.Emit(), errors.Combine()
}
