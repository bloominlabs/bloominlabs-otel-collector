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
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for Resource processor.
type Config struct {
	// The size of the LRU cache for allocation IDs. The nomad processor will
	// cause opentelemetry to fetch each allocation from the nomad API which can
	// be expensive depending on log / batch volume.
	LRUCacheSize int `mapstructure:"cache_size"`

	// TokenFile is a file containing the current token to use for this client.
	// Token or Tokenfile are only required if Nomad's ACL System is enabled.
	TokenFile string `mapstructure:"token_file"`

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default (empty) token.
	// Token or Tokenfile are only required if Nomad's ACL System is enabled.
	Token string `mapstructure:"token"`

	// Client to use when performing requests to nomad. This will be created by
	// default; however, this allows you to pass in a custom nomad client for testing.
	Client *api.Client
}

var _ component.Component = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}

func (cfg *Config) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (cfg *Config) Shutdown(ctx context.Context) error {
	return nil
}
