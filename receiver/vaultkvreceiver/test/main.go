package main

import (
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

// ListRecursive returns secrets to a path recursive.
func (s *Secrets) ListRecursive(client *api.Client, rootPath, subPath string) error {
	keys, err := ListSecrets(client, rootPath, subPath)
	if err != nil {
		// no sub directories in here, but lets check for normal kv pairs then..
		secrets, e := ReadSecrets(client, rootPath, subPath)
		if e == nil {
			(*s)[path.Join(rootPath, subPath)] = secrets

			return nil
		}

		return err
	}

	for _, k := range keys {
		if strings.HasSuffix(k, Delimiter) {
			if err := s.ListRecursive(client, rootPath, path.Join(subPath, k)); err != nil {
				return err
			}
		} else {
			secrets, err := ReadSecrets(client, rootPath, path.Join(subPath, k))
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
func ListSecrets(client *api.Client, rootPath, subPath string) ([]string, error) {
	data, err := client.Logical().List(fmt.Sprintf(listSecretsPath, rootPath, subPath))
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
func ReadSecrets(client *api.Client, rootPath, subPath string) (MetadataResponse, error) {
	var resp MetadataResponse

	data, err := client.Logical().Read(fmt.Sprintf(listSecretsPath, rootPath, subPath))
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
	fmt.Println("err", err)
	if err != nil {
		return resp, err
	}

	err = decoder.Decode(data.Data)
	fmt.Println(err, data.Data)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	secrets := make(Secrets)

	secrets.ListRecursive(client, "infra/", "")
	fmt.Println(secrets)
}
