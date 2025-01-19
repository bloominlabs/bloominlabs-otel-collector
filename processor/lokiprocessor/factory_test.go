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

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestNewProcessor(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		metricsExporter      string
		lruCacheSize         int
		expectedLRUCacheSize int
	}{
		{
			name: "simplest config (use defaults)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare
			factory := NewFactory()

			creationParams := processortest.NewNopSettings()
			cfg := factory.CreateDefaultConfig().(*Config)

			// Test
			logProcessor, err := factory.CreateLogs(context.Background(), creationParams, cfg, consumertest.NewNop())

			// Verify
			assert.NoError(t, err)
			assert.NotNil(t, logProcessor)
		})
	}
}
