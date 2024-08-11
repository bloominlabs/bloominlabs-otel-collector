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

// https://github.com/FalcoSuessgott/vkv

import (
	"context"
	"fmt"
	"strings"

	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hashicorp/go-multierror"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/userstatsreceiver/internal/metadata"
)

type client interface {
	listBackupsByUser(ctx context.Context) (map[string]map[metadata.AttributeType]int64, error)
}

type userStatsClient struct {
	client *s3.Client
	bucket string
}

var _ client = (*userStatsClient)(nil)

func newBackupsUtilizationClient(client *s3.Client, bucket string) (*userStatsClient, error) {
	return &userStatsClient{
		client,
		bucket,
	}, nil
}

func (c *userStatsClient) listBackupsByUser(ctx context.Context) (map[string]map[metadata.AttributeType]int64, error) {
	backupsMap := make(map[string]map[metadata.AttributeType]int64)
	client := c.client

	params := &s3.ListObjectVersionsInput{
		Bucket:  &c.bucket,
		MaxKeys: aws.Int32(1000),
	}

	maxPages := 1000
	pageNum := -1
	p := s3.NewListObjectVersionsPaginator(client, params)
	var pageErrors error
	for p.HasMorePages() && pageNum < maxPages {
		if ctx.Err() != nil {
			return backupsMap, pageErrors
		}
		pageNum = pageNum + 1
		pageCtx, cancel := context.WithTimeout(ctx, time.Second*2)
		output, err := p.NextPage(pageCtx)
		cancel()
		if output != nil {
			fmt.Println(p.HasMorePages(), pageNum, len(output.Versions), len(output.DeleteMarkers))
		} else {
			fmt.Println(p.HasMorePages(), pageNum, err)
		}
		if err != nil {
			pageErrors = multierror.Append(
				pageErrors,
				fmt.Errorf("failed to list page %d from '%s': %w", pageNum, c.bucket, err),
			)
			continue
		}

		for _, version := range output.Versions {
			if version.Key == nil {
				pageErrors = multierror.Append(
					pageErrors,
					fmt.Errorf("version.Key was nil on page %d from '%s'. could not safely index", pageNum, c.bucket),
				)
				continue
			}
			fullPath := strings.Split(strings.TrimSuffix(*version.Key, "/"), "/")
			if len(fullPath) < 3 {
				continue
			}
			userID := fullPath[0]
			// serverID := fullPath[1]
			backupType := metadata.MapAttributeType[strings.Split(fullPath[2], ".")[0]] // rmeove any .tar.gz prefix
			if version.Size == nil {
				pageErrors = multierror.Append(
					pageErrors,
					fmt.Errorf("version.Size was nil on page %d from '%s'. could not safely index", pageNum, c.bucket),
				)
				continue
			}
			if backupsMap[userID] == nil {
				backupsMap[userID] = make(map[metadata.AttributeType]int64)
			}
			backupsMap[userID][backupType] += *version.Size
		}
	}

	return backupsMap, pageErrors
}
