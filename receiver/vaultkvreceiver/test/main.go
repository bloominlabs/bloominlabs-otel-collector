package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/vault/api"
)

func listKeys(ctx context.Context, c *api.Client) ([]string, error) {
	fmt.Println("initial path", "infra/")
	return recursiveListKeys(ctx, c, "infra/")
}

func recursiveListKeys(ctx context.Context, client *api.Client, path string) ([]string, error) {
	var secretListPath []string

	fmt.Println("path", path)
	secretList, err := listSecret(ctx, client, path)
	if err == nil && secretList != nil {
		for _, secret := range secretList.Data["keys"].([]interface{}) {
			if strings.HasSuffix(secret.(string), "/") {
				keys, err := recursiveListKeys(ctx, client, path+secret.(string))
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

func listSecret(ctx context.Context, client *api.Client, path string) (*api.Secret, error) {
	secret, err := client.Logical().List(path)
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

func main() {
	client, _ := api.NewClient(api.DefaultConfig())

	keys, err := listKeys(context.Background(), client)
	fmt.Println(keys, err)
}
