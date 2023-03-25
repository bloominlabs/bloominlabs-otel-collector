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
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/nomad/api"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"k8s.io/utils/lru"

	"github.com/bloominlabs/baseplate-go/config/filesystem"
)

type nomadProcessor struct {
	allocationCache *lru.Cache
	client          *api.Client
	watcher         *filesystem.Watcher
	tokenFile       *string

	// Lock to allow updating the client when rotating credentials on disk
	sync.RWMutex
}

func (rp *nomadProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
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
					if !ok {
						allocationIDVal, ok = lr.Attributes().Get("allocation.id")
						if !ok {
							allocationIDVal, _ = lr.Attributes().Get("nomad.allocation.id")
						}
					}
				}

				var allocation *api.Allocation
				allocationID := allocationIDVal.AsString()
				allocationFromCache, ok := rp.allocationCache.Get(allocationID)
				if !ok {
					rp.RLock()
					alloc, _, err := rp.client.Allocations().Info(allocationID, &api.QueryOptions{AllowStale: true})
					rp.RUnlock()
					if err != nil {
						return ld, fmt.Errorf("failed to get allocation: %w", err)
					}
					allocation = alloc
					rp.allocationCache.Add(allocationID, alloc)
				} else {
					allocation = allocationFromCache.(*api.Allocation)
				}

				if allocation == nil {
					return ld, fmt.Errorf("failed to extract job name from allocation id '%s'. allocation is empty", allocationID)
				}
				if allocation.Job == nil {
					return ld, fmt.Errorf("failed to extract job name from allocation id '%s'. job is empty", allocationID)
				}
				if allocation.Job.Name == nil {
					return ld, fmt.Errorf("failed to extract job name from allocation id '%s'. job name is empty", allocationID)
				}

				lr.Attributes().PutStr("nomad.job.name", *allocation.Job.Name)
				for key, val := range allocation.Job.Meta {
					lr.Attributes().PutStr("nomad.job.meta."+key, val)
				}
			}
		}
	}

	return ld, nil
}

func (s *nomadProcessor) Start(ctx context.Context, host component.Host) error {
	s.RLock()
	defer s.RUnlock()

	if s.watcher != nil {
		(*s.watcher).Start(ctx)
		go func() {
			for range (*s.watcher).EventsCh() {
				secretID, err := os.ReadFile(*s.tokenFile)
				if err != nil {
					panic(err)
				}
				s.Lock()
				s.client.SetSecretID(strings.TrimSpace(string(secretID)))
				s.Unlock()
			}
		}()
	}

	return nil
}

func (s *nomadProcessor) Shutdown(context.Context) error {
	if s.watcher != nil {
		(*s.watcher).Stop()
	}

	return nil
}
