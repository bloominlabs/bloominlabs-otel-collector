// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package digitaloceanreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver"

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/scraper/scraperhelper"
	"go.uber.org/multierr"

	"github.com/digitalocean/godo"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/digitaloceanreceiver/internal"
)

const (
	scrapersKey = "scrapers"
)

// Config defines configuration for digitalocean receiver.
type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	Scrapers                       map[string]internal.Config `mapstructure:"scrapers"`

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default (empty) token.
	// Token or Tokenfile are only required if [Consul's ACL
	// System](https://www.consul.io/docs/security/acl/acl-system) is enabled.
	Token configopaque.String `mapstructure:"token"`

	// TokenFile is a file containing the current token to use for this client.
	// If provided it is read once at startup and never again.
	// Token or Tokenfile are only required if [Consul's ACL
	// System](https://www.consul.io/docs/security/acl/acl-system) is enabled.
	TokenFile string `mapstructure:"token_file"`
}

var (
	_ component.Config    = (*Config)(nil)
	_ confmap.Unmarshaler = (*Config)(nil)
)

// Validate checks the receiver configuration is valid
func (cfg *Config) Validate() error {
	var err error
	if len(cfg.Scrapers) == 0 {
		err = multierr.Append(err, errors.New("must specify at least one scraper when using digitalocean receiver"))
	}
	return err
}

// Unmarshal a config.Parser into the config struct.
func (cfg *Config) Unmarshal(componentParser *confmap.Conf) error {
	if componentParser == nil {
		return nil
	}

	// load the non-dynamic config normally
	err := componentParser.Unmarshal(cfg)
	if err != nil {
		return err
	}

	scrapersSection, err := componentParser.Sub(scrapersKey)
	if err != nil {
		return err
	}

	token := cfg.Token
	if cfg.TokenFile != "" {
		handle, err := os.Open(cfg.TokenFile)
		if err != nil {
			return fmt.Errorf("failed to open handle to %s: %w", cfg.TokenFile, err)
		}
		tokenBytes, err := io.ReadAll(handle)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", cfg.TokenFile, err)
		}
		token = configopaque.String(tokenBytes)
	}

	if token == "" {
		return fmt.Errorf("no token specified. please use the 'token' or 'token_file' configuration option")
	}

	client := godo.NewFromToken(string(token))

	for key := range scrapersSection.ToStringMap() {
		factory, ok := getScraperFactory(key)
		if !ok {
			return fmt.Errorf("invalid scraper key: %s", key)
		}

		collectorCfg := factory.CreateDefaultConfig()
		collectorViperSection, err := scrapersSection.Sub(key)
		if err != nil {
			return err
		}
		err = collectorViperSection.Unmarshal(collectorCfg)
		if err != nil {
			return fmt.Errorf("error reading settings for scraper type %q: %w", key, err)
		}

		collectorCfg.SetClient(client)

		cfg.Scrapers[key] = collectorCfg
	}

	return nil
}
