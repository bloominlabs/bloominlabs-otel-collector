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
	"fmt"
	"sync"

	// "github.com/bloominlabs/baseplate-go/config"
	"github.com/hashicorp/nomad/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"k8s.io/utils/lru"
)

type resourceProcessor struct {
	allocationCache *lru.Cache
	client          *api.Client

	// Lock to allow updating the client when rotating credentials on disk
	sync.RWMutex
}

func (rp *resourceProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	rp.RLock()
	defer rp.Unlock()
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				allocationIDVal, ok := lr.Attributes().Get("nomad_allocation_id")
				if !ok {
					allocationIDVal, ok = lr.Attributes().Get("allocation_id")
				}
				if !ok {
					allocationIDVal, ok = lr.Attributes().Get("allocation.id")
				}
				if !ok {
					allocationIDVal, ok = lr.Attributes().Get("nomad.allocation.id")
				}

				var allocation *api.Allocation
				allocationID := allocationIDVal.AsString()
				allocationFromCache, ok := rp.allocationCache.Get(allocationID)
				if !ok {
					fmt.Println("cache miss for ", allocationIDVal.AsString())
					alloc, _, err := rp.client.Allocations().Info(allocationID, &api.QueryOptions{AllowStale: true})
					if err != nil {
						return ld, fmt.Errorf("failed to get allocation: %w", err)
					}
					allocation = alloc
					rp.allocationCache.Add(allocationID, alloc)
				} else {
					fmt.Println("cache hit for ", allocationID)
					allocation = allocationFromCache.(*api.Allocation)
				}

				lr.Attributes().InsertString("nomad.job.name", *allocation.Job.Name)
				for key, val := range allocation.Job.Meta {
					fmt.Println(key, val)
					lr.Attributes().InsertString("nomad.job.meta."+key, val)
				}
			}
		}
	}
	return ld, nil
}

func (s *resourceProcessor) Start(_ context.Context, host component.Host) error {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	return nil
}

func (s *resourceProcessor) Shutdown(context.Context) error {
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	fmt.Println("=================================")
	return nil
}
