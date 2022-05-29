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

package vaultkvreceiver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc                  string
		defaultConfigModifier func(cfg *Config)
		expected              error
	}{
		{
			desc: "just /",
			defaultConfigModifier: func(cfg *Config) {
				cfg.Mount = "/"
			},
			expected: fmt.Errorf(ErrInvalidMountTmpl, "/"),
		},
		{
			desc: "leading /",
			defaultConfigModifier: func(cfg *Config) {
				cfg.Mount = "/test"
			},
			expected: fmt.Errorf(ErrInvalidMountTmpl, "/test"),
		},
		{
			desc: "missing /",
			defaultConfigModifier: func(cfg *Config) {
				cfg.Mount = "test"
			},
			expected: fmt.Errorf(ErrInvalidMountTmpl, "test"),
		},
		{
			desc: "valid/",
			defaultConfigModifier: func(cfg *Config) {
				cfg.Mount = "test/"
			},
			expected: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig().(*Config)
			tC.defaultConfigModifier(cfg)
			actual := cfg.Validate()
			require.Equal(t, tC.expected, actual)
		})
	}
}
