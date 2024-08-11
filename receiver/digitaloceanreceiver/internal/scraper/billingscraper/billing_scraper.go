// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package billingscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper"

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"

	// "go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/billingscraper/internal/metadata"
)

// scraper for Billing Metrics
type scraper struct {
	settings receiver.Settings
	config   *Config
	mb       *metadata.MetricsBuilder
}

// newBillingScraper creates a set of billing related metrics
func newBillingScraper(_ context.Context, settings receiver.Settings, cfg *Config) *scraper {
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

	balance, _, err := s.config.Client.Balance.Get(ctx)
	if err != nil {
		return s.mb.Emit(), err
	}
	if balance != nil {
		m2dUsage, err := strconv.ParseFloat(balance.MonthToDateUsage, 32)
		if err != nil {
			return s.mb.Emit(), fmt.Errorf("failed to convert month-to-day usage '%s': %w", balance.MonthToDateUsage, err)
		}
		m2dBalance, err := strconv.ParseFloat(balance.MonthToDateBalance, 32)
		if err != nil {
			return s.mb.Emit(), fmt.Errorf("failed to convert month-to-day balance '%s': %w", balance.MonthToDateUsage, err)
		}
		s.mb.RecordDigitaloceanBillingGeneratedAtDataPoint(now, balance.GeneratedAt.Unix())
		s.mb.RecordDigitaloceanBillingUsageDataPoint(now, m2dUsage)
		s.mb.RecordDigitaloceanBillingBalanceDataPoint(now, m2dBalance)
	}

	return s.mb.Emit(), nil
}
