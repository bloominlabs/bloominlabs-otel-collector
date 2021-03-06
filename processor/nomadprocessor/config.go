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
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for Resource processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct

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
}

var _ config.Processor = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	return nil
}
