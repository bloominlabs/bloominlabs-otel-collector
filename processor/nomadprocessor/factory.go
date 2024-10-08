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
	"io"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"k8s.io/utils/lru"

	"github.com/bloominlabs/baseplate-go/config/filesystem"
)

const (
	stability = component.StabilityLevelDevelopment
)

var (
	// The value of "type" key in configuration.
	typeStr = component.MustNewType("nomad")
)

var processorCapabilities = consumer.Capabilities{MutatesData: true}

// NewFactory returns a new factory for the Resource processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithLogs(createLogsProcessor, stability))
}

// Note: This isn't a valid configuration because the processor would do no work.
func createDefaultConfig() component.Config {
	return &Config{
		LRUCacheSize: 1000,
	}
}

func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (processor.Logs, error) {
	oCfg := cfg.(*Config)

	client, _ := api.NewClient(api.DefaultConfig())

	proc := &nomadProcessor{
		allocationCache: lru.New(oCfg.LRUCacheSize),
		client:          client,
	}
	if oCfg.Client != nil {
		proc.client = oCfg.Client
	}

	if oCfg.TokenFile != "" {
		w, err := filesystem.NewRateLimitedFileWatcher([]string{oCfg.TokenFile}, log.With().Logger().Output(io.Discard), time.Second)

		if err != nil {
			return nil, fmt.Errorf("failed to create file watcher: %w", err)
		}

		proc.watcher = &w
		proc.tokenFile = &oCfg.TokenFile
	}

	return processorhelper.NewLogsProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		proc.processLogs,
		processorhelper.WithCapabilities(processorCapabilities),
		processorhelper.WithStart(proc.Start),
		processorhelper.WithShutdown(proc.Shutdown))
}
