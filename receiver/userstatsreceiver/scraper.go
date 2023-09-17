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
	"go.opentelemetry.io/collector/receiver/scrapererror"
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
	settings receiver.CreateSettings,
	config *Config,
	clientFactory userStatsClientFactory,
) *userStatsScraper {
	return &userStatsScraper{
		logger:  settings.Logger,
		rConfig: config,
		s3Config: &spaces.DigitalOceanSpacesConfig{
			Endpoint:        config.Endpoint,
			Region:          config.Region,
			AccessKeyID:     config.AccessKeyID,
			SecretAccessKey: config.SecretAccessKey,
		},
		clientFactory: clientFactory,
		mb:            metricMetadata.NewMetricsBuilder(config.MetricsBuilderConfig, settings),
	}
}

// scrape scrapes the metric stats, transforms them and attributes them into a metric slices.
func (p *userStatsScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
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
	backupsMap, err := client.listBackupsByUser(ctx)
	if err != nil {
		p.logger.Error("failed to list backups", zap.Error(err))
		return pmetric.NewMetrics(), err
	}

	for userID, size := range backupsMap {
		p.mb.RecordBackupsTotalSizeDataPoint(
			now,
			size,
			userID,
		)
	}

	return p.mb.Emit(), errors.Combine()
}
