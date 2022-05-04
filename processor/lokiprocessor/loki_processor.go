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

	"go.opentelemetry.io/collector/pdata/plog"
)

type resourceProcessor struct{}

func (rp *resourceProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				body := lr.Body().MapVal()
				message, exists := body.Get("MESSAGE")
				if !exists {
					continue
				}
				unit, exists := body.Get("_SYSTEMD_UNIT")
				if !exists {
					continue
				}
				hostname, exists := body.Get("_HOSTNAME")
				if !exists {
					continue
				}

				lr.Body().SetStringVal(message.StringVal())
				lr.Attributes().Insert("unit", unit)
				lr.Attributes().Insert("hostname", hostname)
			}
		}
	}
	return ld, nil
}
