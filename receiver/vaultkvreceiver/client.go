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

// https://github.com/FalcoSuessgott/vkv

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

const (
	Delimiter            = "/"
	mountEnginePath      = "sys/mounts/%s"
	readWriteSecretsPath = "%s/data/%s"
	listSecretsPath      = "%s/metadata/%s"
)

type MetadataResponse struct {
	CreatedTime    time.Time `mapstructure:"created_time"`
	CurrentVersion int       `mapstructure:"current_version"`
	MaxVersions    int       `mapstructure:"max_versions"`
	UpdatedTime    time.Time `mapstructure:"updated_time"`
}

// Secrets holds all recursive secrets of a certain path.
type Secrets map[string]MetadataResponse

type client interface {
	listSecretMetadata(ctx context.Context) (Secrets, error)
	listSecrets(rootPath, subPath string) ([]string, error)
	readSecrets(rootPath, subPath string) (MetadataResponse, error)
}

type vaultKVClient struct {
	client *api.Client
	mount  string
}

var _ client = (*vaultKVClient)(nil)

func newVaultKVClient(client *api.Client) (*vaultKVClient, error) {
	return &vaultKVClient{
		client: client,
		// TODO
		mount: "infra/",
	}, nil
}

func (c *vaultKVClient) listSecretMetadata(ctx context.Context) (Secrets, error) {
	secrets := make(Secrets)
	err := secrets.listRecursive(c, c.mount, "")
	return secrets, err
}

// ListRecursive returns secrets to a path recursive.
func (s *Secrets) listRecursive(client *vaultKVClient, rootPath, subPath string) error {
	keys, err := client.listSecrets(rootPath, subPath)
	if err != nil {
		// no sub directories in here, but lets check for normal kv pairs then..
		secrets, e := client.readSecrets(rootPath, subPath)
		if e == nil {
			(*s)[path.Join(rootPath, subPath)] = secrets

			return nil
		}

		return err
	}

	for _, k := range keys {
		if strings.HasSuffix(k, Delimiter) {
			if err := s.listRecursive(client, rootPath, path.Join(subPath, k)); err != nil {
				return err
			}
		} else {
			secrets, err := client.readSecrets(rootPath, path.Join(subPath, k))
			if err != nil {
				(*s)[path.Join(rootPath, subPath, k)] = MetadataResponse{}

				continue
			}

			(*s)[path.Join(rootPath, subPath, k)] = secrets
		}
	}

	return nil
}

// ListSecrets returns all keys from vault kv secret path.
func (c *vaultKVClient) listSecrets(rootPath, subPath string) ([]string, error) {
	data, err := c.client.Logical().List(fmt.Sprintf(listSecretsPath, rootPath, subPath))
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no secrets under path \"%s\" found", path.Join(rootPath, subPath))
	}

	if data.Data != nil {
		keys := []string{}

		k, ok := data.Data["keys"].([]interface{})
		if !ok {
			log.Fatalf("did not found any keys in %s/%s", rootPath, subPath)
		}

		for _, e := range k {
			keys = append(keys, fmt.Sprintf("%v", e))
		}

		return keys, nil
	}

	return nil, fmt.Errorf("no secrets in %s found", path.Join(rootPath, subPath))
}

// ReadSecrets returns a map with all secrets from a kv engine path.
func (c *vaultKVClient) readSecrets(rootPath, subPath string) (MetadataResponse, error) {
	var resp MetadataResponse

	data, err := c.client.Logical().Read(fmt.Sprintf(listSecretsPath, rootPath, subPath))
	if err != nil {
		return resp, err
	}

	if data == nil {
		return resp, fmt.Errorf("no secrets in %s found", path.Join(rootPath, subPath))
	}

	decoder, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result:     &resp,
			DecodeHook: mapstructure.StringToTimeHookFunc(time.RFC3339),
		},
	)
	if err != nil {
		return resp, err
	}

	err = decoder.Decode(data.Data)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
