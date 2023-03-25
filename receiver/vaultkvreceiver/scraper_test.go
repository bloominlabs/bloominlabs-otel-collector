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

package vaultkvreceiver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestUnsuccessfulScrape(t *testing.T) {
	originalVaultAddr := os.Getenv("VAULT_ADDR")
	os.Setenv("VAULT_ADDR", "")
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.Mount = "test/"

	scraper := newVaultKVScraper(receivertest.NewNopCreateSettings(), cfg, &defaultClientFactory{})
	actualMetrics, err := scraper.scrape(context.Background())
	require.Error(t, err)

	require.NoError(t, pmetrictest.CompareMetrics(pmetric.NewMetrics(), actualMetrics,
		pmetrictest.IgnoreMetricDataPointsOrder(), pmetrictest.IgnoreStartTimestamp(), pmetrictest.IgnoreTimestamp()))

	os.Setenv("VAULT_ADDR", originalVaultAddr)
}

func TestScraper(t *testing.T) {
	mockTime, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	testCases := []struct {
		desc    string
		secrets Secrets
	}{
		{
			desc: "missing type in metadata",
			secrets: Secrets{
				"infra/key1": MetadataResponse{
					CreatedTime: mockTime,
				},
			},
		},
		{
			desc: "invalid type in metadata",
			secrets: Secrets{
				"infra/key1": MetadataResponse{
					CreatedTime: mockTime,
					CustomMetadata: map[string]string{
						"type": "invalid",
					},
				},
			},
		},
		{
			desc: "valid type",
			secrets: Secrets{
				"infra/key1": MetadataResponse{
					CreatedTime: mockTime,
					CustomMetadata: map[string]string{
						"type": "digitalocean.api",
					},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			factory := new(mockClientFactory)
			factory.initMocks(tC.secrets)

			cfg := createDefaultConfig().(*Config)
			scraper := newVaultKVScraper(receivertest.NewNopCreateSettings(), cfg, factory)

			actualMetrics, err := scraper.scrape(context.Background())
			require.NoError(t, err)

			safeFileName := strings.ReplaceAll(tC.desc, " ", "_")
			expectedFile := filepath.Join("testdata", "scraper", "testScraper", safeFileName+".json")
			// golden.WriteMetrics(t, expectedFile, actualMetrics)
			expectedMetrics, err := golden.ReadMetrics(expectedFile)
			require.NoError(t, err)

			require.NoError(t, pmetrictest.CompareMetrics(actualMetrics, expectedMetrics,
				pmetrictest.IgnoreMetricDataPointsOrder(), pmetrictest.IgnoreStartTimestamp(), pmetrictest.IgnoreTimestamp()))
		})
	}
}

type mockClientFactory struct{ mock.Mock }
type mockClient struct{ mock.Mock }

func (m *mockClient) listSecretMetadata(_ context.Context) (Secrets, error) {
	args := m.Called()
	return args.Get(0).(Secrets), args.Error(1)
}

func (m *mockClient) listSecrets(rootPath, subPath string) ([]string, error) {
	args := m.Called(rootPath, subPath)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockClient) readSecrets(rootPath, subPath string) (MetadataResponse, error) {
	args := m.Called(rootPath, subPath)
	return args.Get(0).(MetadataResponse), args.Error(1)
}

func (m *mockClientFactory) getClient(c *Config, mount string) (client, error) {
	args := m.Called(mount)
	return args.Get(0).(client), args.Error(1)
}

func (m *mockClientFactory) initMocks(secrets Secrets) {
	client := new(mockClient)
	client.initMocks(secrets)
	m.On("getClient", "infra/").Return(client, nil)
}

func (m *mockClient) initMocks(secrets Secrets) {
	m.On("listSecretMetadata").Return(secrets, nil)
}
