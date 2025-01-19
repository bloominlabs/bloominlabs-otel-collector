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

package lokiprocessor

import (
	"context"
	"testing"
	"time"

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
	baseMessage := "2022-01-01 01:02:03 INFO this is a test message"
	pBaseMessage := pcommon.NewValueStr(baseMessage)

	val := pcommon.NewValueMap()
	val.Map().PutStr("_SYSTEMD_UNIT", "test_systemd_unit")
	val.Map().PutStr("MESSAGE", baseMessage)
	val.Map().PutStr("_HOSTNAME", "test_hostname")

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
					body: val,
				},
			},
			parsedMessages: []testLogMessage{
				{
					body: pBaseMessage,
					attributes: &map[string]pcommon.Value{
						"unit":     pcommon.NewValueStr("test_systemd_unit"),
						"hostname": pcommon.NewValueStr("test_hostname"),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
