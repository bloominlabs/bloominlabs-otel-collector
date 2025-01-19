// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nomadprocessor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor/processortest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
)

var cfg = &Config{}

func createNomadTestHarness(t *testing.T) (*httptest.Server, *http.ServeMux, *api.Client) {
	t.Helper()

	// test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// Cloudflare client configured to use test server
	cfg := api.DefaultConfig()
	cfg.Address = server.URL
	client, err := api.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		server.Close()
	})

	return server, mux, client
}

type testLogMessage struct {
	body         pcommon.Value
	time         *time.Time
	observedTime *time.Time
	severity     plog.SeverityNumber
	attributes   *map[string]pcommon.Value
}

func parseTime(format, input string) *time.Time {
	val, _ := time.ParseInLocation(format, input, time.Local)
	return &val
}

func TestLogsTransformProcessor(t *testing.T) {
	baseMessage := pcommon.NewValueStr("2022-01-01 01:02:03 INFO this is a test message")

	tests := []struct {
		name           string
		config         *Config
		sourceMessages []testLogMessage
		parsedMessages []testLogMessage
	}{
		{
			name:   "simpleTest",
			config: cfg,
			sourceMessages: []testLogMessage{
				{
					body: baseMessage,
					attributes: &map[string]pcommon.Value{
						"nomad_allocation_id": pcommon.NewValueStr("test"),
					},
				},
				{
					body: baseMessage,
					attributes: &map[string]pcommon.Value{
						"nomad.allocation.id": pcommon.NewValueStr("nomad.allocation.id"),
					},
				},
				{
					body: baseMessage,
					attributes: &map[string]pcommon.Value{
						"nomad.allocation.id": pcommon.NewValueStr("add-metadata"),
					},
				},
			},
			parsedMessages: []testLogMessage{
				{
					body: baseMessage,
					attributes: &map[string]pcommon.Value{
						"nomad.job.name":      pcommon.NewValueStr("test"),
						"nomad_allocation_id": pcommon.NewValueStr("test"),
					},
				},
				{
					body: baseMessage,
					attributes: &map[string]pcommon.Value{
						"nomad.job.name":      pcommon.NewValueStr("nomad.allocation.id"),
						"nomad.allocation.id": pcommon.NewValueStr("nomad.allocation.id"),
					},
				},
				{
					body: baseMessage,
					attributes: &map[string]pcommon.Value{
						"nomad.job.name":           pcommon.NewValueStr("add-metadata"),
						"nomad.allocation.id":      pcommon.NewValueStr("add-metadata"),
						"nomad.job.meta.testkey":   pcommon.NewValueStr("testval"),
						"nomad.job.meta.user.id":   pcommon.NewValueStr("12345"),
						"nomad.job.meta.server.id": pcommon.NewValueStr("12345"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, nomadMux, nomadClient := createNomadTestHarness(t)
			nomadHandler := func(w http.ResponseWriter, r *http.Request) {
				assert.Regexp(t, regexp.MustCompile("/v1/allocation/.+"), r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
				allocationID := strings.TrimPrefix(r.URL.Path, "/v1/allocation/")
				resp := api.Allocation{
					ID:   allocationID,
					Name: "test",
					Job: &api.Job{
						Name: &allocationID,
					},
				}

				if allocationID == "add-metadata" {
					resp.Job.Meta = map[string]string{
						"testkey":   "testval",
						"user.id":   "12345",
						"server.id": "12345",
					}
				}

				w.Header().Set("content-type", "application/json")
				json.NewEncoder(w).Encode(resp)
			}
			nomadMux.HandleFunc("/", nomadHandler)

			tt.config.Client = nomadClient

			tln := new(consumertest.LogsSink)
			factory := NewFactory()
			ltp, err := factory.CreateLogs(context.Background(), processortest.NewNopSettings(), tt.config, tln)
			require.NoError(t, err)
			assert.True(t, ltp.Capabilities().MutatesData)

			err = ltp.Start(context.Background(), nil)
			require.NoError(t, err)

			sourceLogData := generateLogData(tt.sourceMessages)
			wantLogData := generateLogData(tt.parsedMessages)
			err = ltp.ConsumeLogs(context.Background(), sourceLogData)
			require.NoError(t, err)
			time.Sleep(200 * time.Millisecond)
			logs := tln.AllLogs()
			require.Len(t, logs, 1)
			assert.NoError(t, plogtest.CompareLogs(wantLogData, logs[0]))
		})
	}
}

func generateLogData(messages []testLogMessage) plog.Logs {
	ld := testdata.GenerateLogsOneEmptyResourceLogs()
	scope := ld.ResourceLogs().At(0).ScopeLogs().AppendEmpty()
	for _, content := range messages {
		log := scope.LogRecords().AppendEmpty()
		content.body.CopyTo(log.Body())
		if content.time != nil {
			log.SetTimestamp(pcommon.NewTimestampFromTime(*content.time))
		}
		if content.observedTime != nil {
			log.SetObservedTimestamp(pcommon.NewTimestampFromTime(*content.observedTime))
		}
		if content.severity != 0 {
			log.SetSeverityNumber(content.severity)
		}
		if content.attributes != nil {
			for k, v := range *content.attributes {
				v.CopyTo(log.Attributes().PutEmpty(k))
			}
		}
	}

	return ld
}
