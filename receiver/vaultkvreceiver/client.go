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
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
)

type client interface {
	listKeys(ctx context.Context) ([]string, error)
	recursiveListKeys(ctx context.Context, path string) ([]string, error)
	listSecret(ctx context.Context, path string) (*api.Secret, error)
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

type MetricStat struct {
	database string
	table    string
	stats    map[string]string
}

func (c *vaultKVClient) listKeys(ctx context.Context) ([]string, error) {
	fmt.Println("initial path", c.mount)
	return c.recursiveListKeys(ctx, c.mount)
}

func (c *vaultKVClient) recursiveListKeys(ctx context.Context, path string) ([]string, error) {
	var secretListPath []string

	fmt.Println("path", path)
	secretList, err := c.listSecret(ctx, path)
	if err == nil && secretList != nil {
		for _, secret := range secretList.Data["keys"].([]interface{}) {
			if strings.HasSuffix(secret.(string), "/") {
				keys, err := c.recursiveListKeys(ctx, path+secret.(string))
				if err != nil {
					return secretListPath, fmt.Errorf("failed to list keys under %s: %w", path, err)
				}
				secretListPath = append(secretListPath, keys...)
			} else {
				secretListPath = append([]string{strings.Replace(path, "metadata", "data", -1) + secret.(string)}, secretListPath...)
			}
		}
	}
	return secretListPath, nil
}

func (c *vaultKVClient) listSecret(ctx context.Context, path string) (*api.Secret, error) {
	secret, err := c.client.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't list from vault: %w", err)
	}

	if isNil(secret) {
		return nil, fmt.Errorf("listed %s but was nil", path)
	}
	return secret, err
}

func isNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
