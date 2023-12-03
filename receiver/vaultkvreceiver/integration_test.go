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

//go:build integration
// +build integration

package vaultkvreceiver

import (
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
)

type configFunc func(hostname string) *Config

type testCase struct {
	name         string
	cfg          configFunc
	expectedFile string
}

func TestVaultIntegration(t *testing.T) {
	testCases := []testCase{
		{
			name: "single_db",
			cfg: func(addr string) *Config {
				f := NewFactory()
				cfg := f.CreateDefaultConfig().(*Config)
				cfg.URL = addr
				return cfg
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srv, client := createTestVault(t)
			defer srv.Close()

			expectedFile := filepath.Join("testdata", "integration", tc.name+".json")
			expectedMetrics, err := golden.ReadMetrics(expectedFile)
			require.NoError(t, err)

			f := NewFactory()
			consumer := new(consumertest.MetricsSink)
			settings := componenttest.NewNopReceiverCreateSettings()
			rcvr, err := f.CreateMetricsReceiver(context.Background(), settings, tc.cfg(client.Address()), consumer)
			require.NoError(t, err, "failed creating metrics receiver")
			require.NoError(t, rcvr.Start(context.Background(), componenttest.NewNopHost()))
			require.Eventuallyf(t, func() bool {
				return consumer.DataPointCount() > 0
			}, 2*time.Minute, 1*time.Second, "failed to receive more than 0 metrics")

			actualMetrics := consumer.AllMetrics()[0]

			require.NoError(t, pmetrictest.CompareMetrics(expectedMetrics, actualMetrics,
				pmetrictest.IgnoreMetricValues(), pmetrictest.IgnoreMetricDataPointsOrder(),
				pmetrictest.IgnoreStartTimestamp(), pmetrictest.IgnoreTimestamp()))
		})
	}
}

func createTestVault(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	coreConfig := &vault.CoreConfig{}

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, _, token := vault.TestCoreUnsealedWithConfig(t, coreConfig)

	// Start an HTTP server for the core.
	srv, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	return srv, client
}
