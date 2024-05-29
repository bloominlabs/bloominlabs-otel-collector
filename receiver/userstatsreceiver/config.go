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

package userstatsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver"

import (
	"fmt"

	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver/internal/metadata"
)

const (
	ErrInvalidMountTmpl = "invalid 'mount' (%s) provided. should be in the form '<mount>/'"
)

type Config struct {
	Endpoint                                string `mapstructure:"endpoint"`
	Region                                  string `mapstructure:"region"`
	AccessKeyID                             string `mapstructure:"access_key_id"`
	SecretAccessKey                         string `mapstructure:"secret_access_key"`
	Bucket                                  string `mapstructure:"bucket_name"`
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	metadata.MetricsBuilderConfig           `mapstructure:",squash"`
}

func (cfg *Config) Validate() error {
	if cfg.Region == "" {
		return fmt.Errorf("config.Region must be specified")
	}
	if cfg.Region == "" {
		return fmt.Errorf("config.Bucket must be specified")
	}
	if cfg.AccessKeyID == "" {
		return fmt.Errorf("config.AccessKeyID must be specified")
	}
	if cfg.SecretAccessKey == "" {
		return fmt.Errorf("config.SecretAccessKey must be specified")
	}

	return nil
}
