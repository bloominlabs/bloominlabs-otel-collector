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

package certificatesreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/certificatesreceiver"

import (
	"fmt"

	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/certificatesreceiver/internal/metadata"
)

const (
	ErrInvalidMountTmpl = "invalid 'mount' (%s) provided. should be in the form '<mount>/'"
)

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	metadata.MetricsBuilderConfig           `mapstructure:",squash"`

	CertificateIncludeGlobs []string `mapstructure:"certificate_include_globs"`
	CertificateExcludeGlobs []string `mapstructure:"certificate_exclude_globs"`
}

func (cfg *Config) Validate() error {
	if len(cfg.CertificateIncludeGlobs) < 1 {
		return fmt.Errorf("config.certificatesIncludeGlob must be specified")
	}

	return nil
}
