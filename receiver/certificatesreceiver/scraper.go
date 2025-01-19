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

package certificatesreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/certificatesreceiver"

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/scraper/scrapererror"
	"go.uber.org/zap"
	"golang.org/x/crypto/pkcs12"

	metricMetadata "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/certificatesreceiver/internal/metadata"
)

type certificatesScraper struct {
	logger *zap.Logger
	config *Config
	mb     *metricMetadata.MetricsBuilder
}

type defaultClientFactory struct{}

func newCertificatesScraper(
	settings receiver.Settings,
	config *Config,
) *certificatesScraper {
	return &certificatesScraper{
		logger: settings.Logger,
		config: config,
		mb:     metricMetadata.NewMetricsBuilder(config.MetricsBuilderConfig, settings),
	}
}

// scrape scrapes the metric stats, transforms them and attributes them into a metric slices.
func (p *certificatesScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return p.mb.Emit(), fmt.Errorf("failed to get hostname: %w", err)
	}
	now := pcommon.NewTimestampFromTime(time.Now())

	matches, err := p.getMatches()
	if err != nil {
		p.logger.Error("failed to get matches", zap.Error(err))
		return p.mb.Emit(), err
	}

	var errors scrapererror.ScrapeErrors
	for _, match := range matches {
		results, err := secondsToExpiryFromCertAsFile(match)
		if err != nil {
			errors.Add(fmt.Errorf("failed to get certificates metrics for %s: %w", match, err))
			continue
		}

		for _, result := range results {
			p.mb.RecordCertificatesCertExpiresInSecondsDataPoint(now, int64(result.durationUntilExpiry), match, hostname, result.cn)
		}
	}

	return p.mb.Emit(), errors.Combine()
}

func (p *certificatesScraper) getMatches() ([]string, error) {
	set := map[string]bool{}
	for _, glob := range p.config.CertificateIncludeGlobs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			p.logger.Error("failed to create glob", zap.String("CertificatesIncludeGlob", glob), zap.Error(err))
			continue
		}
		for _, match := range matches {
			set[match] = true
		}
	}

	for _, glob := range p.config.CertificateExcludeGlobs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			p.logger.Error("failed to create glob", zap.String("CertificatesExcludeGlob", glob), zap.Error(err))
			continue
		}

		for _, match := range matches {
			delete(set, match)
		}
	}

	res := make([]string, len(set))
	i := 0
	for k := range set {
		res[i] = k
		i++
	}
	return res, nil
}

type certMetric struct {
	durationUntilExpiry float64
	notAfter            float64
	issuer              string
	cn                  string
}

func secondsToExpiryFromCertAsFile(file string) ([]certMetric, error) {
	certBytes, err := os.ReadFile(file)
	if err != nil {
		return []certMetric{}, err
	}

	return secondsToExpiryFromCertAsBytes(certBytes)
}

func secondsToExpiryFromCertAsBase64String(s string) ([]certMetric, error) {
	certBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return []certMetric{}, err
	}

	return secondsToExpiryFromCertAsBytes(certBytes)
}

func secondsToExpiryFromCertAsBytes(certBytes []byte) ([]certMetric, error) {
	var metrics []certMetric

	parsed, metrics, err := parseAsPEM(certBytes)
	if parsed {
		return metrics, err
	}
	// Parse as PKCS ?
	parsed, metrics, err = parseAsPKCS(certBytes)
	if parsed {
		return metrics, nil
	}
	return nil, fmt.Errorf("failed to parse as pem and pkcs12: %w", err)
}

func getCertificatesMetrics(cert *x509.Certificate) certMetric {
	var metric certMetric
	metric.notAfter = float64(cert.NotAfter.Unix())
	metric.durationUntilExpiry = time.Until(cert.NotAfter).Seconds()
	metric.issuer = cert.Issuer.CommonName
	metric.cn = cert.Subject.CommonName
	return metric
}

func parseAsPKCS(certBytes []byte) (bool, []certMetric, error) {
	var metrics []certMetric
	var blocks []*pem.Block
	var last_err error

	pfx_blocks, err := pkcs12.ToPEM(certBytes, "")
	if err != nil {
		return false, nil, err
	}
	for _, b := range pfx_blocks {
		if b.Type == "CERTIFICATE" {
			blocks = append(blocks, b)
		}
	}

	for _, block := range blocks {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err == nil {
			metric := getCertificatesMetrics(cert)
			metrics = append(metrics, metric)
		} else {
			last_err = err
		}
	}
	return true, metrics, last_err
}

func parseAsPEM(certBytes []byte) (bool, []certMetric, error) {
	var metrics []certMetric
	var blocks []*pem.Block

	block, rest := pem.Decode(certBytes)
	if block == nil {
		return false, metrics, fmt.Errorf("failed to parse as a pem")
	}
	blocks = append(blocks, block)
	// Export the remaining certificatess in the certificates chain
	for len(rest) != 0 {
		block, rest = pem.Decode(rest)
		if block == nil {
			return true, metrics, fmt.Errorf("failed to parse intermediate as a pem")
		}
		if block.Type == "CERTIFICATE" {
			blocks = append(blocks, block)
		}
	}
	for _, block := range blocks {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return true, metrics, err
		}
		metric := getCertificatesMetrics(cert)
		metrics = append(metrics, metric)
	}
	return true, metrics, nil
}
