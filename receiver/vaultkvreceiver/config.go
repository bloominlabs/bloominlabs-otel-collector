// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vaultkvreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver"

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/vaultkvreceiver/internal/metadata"
)

const (
	ErrInvalidMountTmpl = "invalid 'mount' (%s) provided. should be in the form '<mount>/'"
)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	Mount                                   string                   `mapstructure:"mount"`
	URL                                     string                   `mapstructure:"url"`
	Metrics                                 metadata.MetricsSettings `mapstructure:"metrics"`
}

func (cfg *Config) Validate() error {
	split := strings.Split(cfg.Mount, "/")

	if len(split) != 2 || split[0] == "" || split[1] != "" {
		return fmt.Errorf(ErrInvalidMountTmpl, cfg.Mount)
	}
	return nil
}
