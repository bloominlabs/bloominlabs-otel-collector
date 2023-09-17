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

package backuputilizationreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/backuputilizationreceiver"

// https://github.com/FalcoSuessgott/vkv

import (
	"context"
	"fmt"
	"path/filepath"
	// "strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/go-multierror"
)

type client interface {
	listBackupsByUser(ctx context.Context) (map[string]int64, error)
}

type backupsUtilizationClient struct {
	client *s3.Client
	bucket string
}

var _ client = (*backupsUtilizationClient)(nil)

func newBackupsUtilizationClient(client *s3.Client, bucket string) (*backupsUtilizationClient, error) {
	return &backupsUtilizationClient{
		client,
		bucket,
	}, nil
}

func (c *backupsUtilizationClient) listBackupsByUser(ctx context.Context) (map[string]int64, error) {
	backupsMap := make(map[string]int64)
	client := c.client

	params := &s3.ListObjectVersionsInput{
		Bucket: &c.bucket,
	}

	maxPages := 1000
	pageNum := 0
	p := s3.NewListObjectVersionsPaginator(client, params)
	var pageErrors error
	for p.HasMorePages() && pageNum < maxPages {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		output, err := p.NextPage(ctx)
		cancel()
		if err != nil {
			pageErrors = multierror.Append(
				pageErrors,
				fmt.Errorf("failed to list page %d from '%s': %w", pageNum, c.bucket, err),
			)
			continue
		}

		for _, version := range output.Versions {
			userID := filepath.Dir(*version.Key)
			// serverID := strings.Split(filepath.Base(*version.Key), ".")[0]
			backupsMap[userID] += version.Size
		}
	}

	return backupsMap, pageErrors
}
