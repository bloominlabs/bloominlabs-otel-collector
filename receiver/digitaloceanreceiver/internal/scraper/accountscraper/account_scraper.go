// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package accountscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/accountscraper"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/accountscraper/internal/metadata"
)

// scraper for Account Metrics
type scraper struct {
	settings receiver.Settings
	config   *Config
	mb       *metadata.MetricsBuilder
}

// newDropletScraper creates a set of account related metrics
func newDropletScraper(_ context.Context, settings receiver.Settings, cfg *Config) *scraper {
	return &scraper{
		settings: settings,
		config:   cfg,
	}
}

func (s *scraper) start(context.Context, component.Host) error {
	s.mb = metadata.NewMetricsBuilder(s.config.MetricsBuilderConfig, s.settings, metadata.WithStartTime(pcommon.Timestamp(time.Now().Unix())))
	return nil
}

func (s *scraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	now := pcommon.NewTimestampFromTime(time.Now())
	account, _, err := s.config.Client.Account.Get(ctx)

	if account != nil {
		s.mb.RecordDigitaloceanAccountDropletLimitDataPoint(now, int64(account.DropletLimit))
		s.mb.RecordDigitaloceanAccountVolumeLimitDataPoint(now, int64(account.VolumeLimit))
	}

	return s.mb.Emit(), err
}
