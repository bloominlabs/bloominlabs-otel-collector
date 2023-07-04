// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package dropletscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/dropletscraper"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"

	// "go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/digitalocean/godo"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal/scraper/dropletscraper/internal/metadata"
)

// scraper for Billing Metrics
type scraper struct {
	settings receiver.CreateSettings
	config   *Config
	mb       *metadata.MetricsBuilder
}

// newDropletScraper creates a set of billing related metrics
func newDropletScraper(_ context.Context, settings receiver.CreateSettings, cfg *Config) *scraper {
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

	list := []godo.Droplet{}

	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := s.config.Client.Droplets.List(ctx, opt)
		if err != nil {
			return s.mb.Emit(), err
		}

		// append the current page's droplets to our list
		list = append(list, droplets...)

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return s.mb.Emit(), err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	for _, droplet := range list {
		val := 1
		if droplet.Status != "active" {
			val = 0
		}
		s.mb.RecordDigitaloceanDropletUpDataPoint(now, int64(val), int64(droplet.ID), droplet.Name, droplet.Region.Slug)
	}

	return s.mb.Emit(), nil
}
