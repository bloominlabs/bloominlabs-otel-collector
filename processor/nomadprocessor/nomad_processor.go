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

	"github.com/hashicorp/nomad/api"
	"go.opentelemetry.io/collector/pdata/plog"
	"k8s.io/utils/lru"
)

type resourceProcessor struct {
	allocationCache *lru.Cache
	client          *api.Client
}

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
				// TODO, use error
				allocationID, _ := lr.Attributes().Get("allocation_id")

				var allocation *api.Allocation

				allocationFromCache, ok := rp.allocationCache.Get(allocationID)
				if !ok {
					// TODO: use error
					alloc, _, _ := rp.client.Allocations().Info(allocationID.AsString(), &api.QueryOptions{AllowStale: true})
					allocation = alloc
					rp.allocationCache.Add(allocationID, alloc)
				} else {
					allocation = allocationFromCache.(*api.Allocation)
				}

				for key, val := range allocation.Job.Meta {
					lr.Attributes().InsertString(key, val)
				}
			}
		}
	}
	return ld, nil
}
