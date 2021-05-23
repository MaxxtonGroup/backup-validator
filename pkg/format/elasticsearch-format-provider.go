package format

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MaxxtonGroup/backup-validator/pkg/runtime"
)

type ElasticsearchSnapshotRepository struct {
	Type     string                                 `yaml:"type" json:"type"`
	Settings map[string]string                      `yaml:"settings" json:"settings"`
	Keystore *map[string]ElasticsearchKeystoreValue `yaml:"keystore" json:"keystore"`
}

type ElasticsearchKeystoreValue struct {
	Value    *string `yaml:"value"`
	FromFile *string `yaml:"fromFile"`
}

type ElasticsearchFormatProvider struct {
	runtimeProvider runtime.RuntimeProvider
	repository      ElasticsearchSnapshotRepository
}

type ElasticsearchQueryResult struct {
	Hits *ElasticsearchQueryHit `json:"hits"`
}

type ElasticsearchQueryHit struct {
	Hits []*ElasticsearchQueryDocument `json:"hits"`
}

type ElasticsearchQueryDocument struct {
	Source map[string]interface{} `json:"_source"`
}

func (p ElasticsearchFormatProvider) Setup(testName string, dir string) error {
	err := p.runtimeProvider.Setup(testName, dir)
	if err != nil {
		return err
	}

	// Configure keystore of elasticsearch node
	if p.repository.Keystore != nil {
		// create keystore
		log.Printf("[%s] Create Keystore", testName)
		_, err = p.runtimeProvider.Exec(testName, "bash", "-c", "if [[ ! -f /usr/share/elasticsearch/config/elasticsearch.keystore ]] ; then /usr/share/elasticsearch/bin/elasticsearch-keystore create; else true; fi")
		if err != nil {
			return err
		}

		for key, value := range *p.repository.Keystore {
			// get key value
			var keyValue string
			if value.Value != nil {
				keyValue = *value.Value
			} else if value.FromFile != nil {
				bytes, err := ioutil.ReadFile(*value.FromFile)
				if err != nil {
					return err
				}
				keyValue = string(bytes)
			} else {
				return fmt.Errorf("keystore '%s' doesn't has a 'value' or 'fromFile' field", key)
			}

			// Save key value in a tmp file that is mounted in the container
			keyFile, err := ioutil.TempFile(dir, ".keystore")
			if err != nil {
				return err
			}
			defer os.Remove(keyFile.Name())
			err = ioutil.WriteFile(keyFile.Name(), []byte(keyValue), os.ModePerm)
			if err != nil {
				return err
			}

			// Store value in keystore
			log.Printf("[%s] Store %s in keystore", testName, key)
			_, err = p.runtimeProvider.Exec(testName, "/usr/share/elasticsearch/bin/elasticsearch-keystore", "add-file", "-f", key, "/mnt/host/"+filepath.Base(keyFile.Name()))
			if err != nil {
				return err
			}
		}

		// Reload keystore
		log.Printf("[%s] Reload keystore", testName)
		_, err = p.runtimeProvider.Exec(testName, "curl", "--fail", "-X", "POST", "http://localhost:9200/_nodes/reload_secure_settings?pretty", "-H", "Content-Type: application/json", "-d", "{}")
		if err != nil {
			return err
		}

	}

	// Configure snapshot repository
	bytes, err := json.Marshal(p.repository)
	if err != nil {
		return err
	}
	log.Printf("[%s] Configure snapshot repository", testName)
	output, err := p.runtimeProvider.Exec(testName, "curl", "--output", "/dev/stdout", "--write-out", "%{http_code}", "-X", "PUT", "http://localhost:9200/_snapshot/backup", "-H", "Content-Type: application/json", "-d", string(bytes))
	if err != nil {
		return err
	}
	if !strings.HasSuffix(*output, "200") {
		return fmt.Errorf("the requested URL returned error: %s", *output)
	}

	return nil
}

func (p ElasticsearchFormatProvider) Destroy(testName string, dir string) error {
	return p.runtimeProvider.Destroy(testName, dir)
}

func (p ElasticsearchFormatProvider) ImportData(testName string, dir string, options []string) error {
	// Handled by the ElasticsearchBackupProvider
	return nil
}

func (p ElasticsearchFormatProvider) ListDatabases(testName string) ([]string, error) {
	output, err := p.runtimeProvider.Exec(testName, "curl", "--fail", "-X", "GET", "http://localhost:9200/_cat/indices?h=index")
	if err != nil {
		return nil, err
	}

	return strings.Split(*output, "\n"), nil
}

func (p ElasticsearchFormatProvider) ListTables(testName string, database string) ([]string, error) {
	output, err := p.runtimeProvider.Exec(testName, "curl", "--fail", "-X", "GET", "http://localhost:9200/"+database+"/_search?size=1")
	if err != nil {
		return nil, err
	}

	result := ElasticsearchQueryResult{}
	err = json.Unmarshal([]byte(*output), &result)
	if err != nil {
		return nil, err
	}

	fields := []string{}
	if len(result.Hits.Hits) > 0 {
		for key := range result.Hits.Hits[0].Source {
			fields = append(fields, key)
		}
	}
	return fields, nil
}

func (p ElasticsearchFormatProvider) GetDatabaseSize(testName string, database string) (*uint64, error) {
	output, err := p.runtimeProvider.Exec(testName, "curl", "--fail", "-X", "GET", "http://localhost:9200/_cat/indices/"+database+"?h=store.size&bytes=b")
	if err != nil {
		return nil, err
	}

	sizeString := strings.TrimSpace(*output)
	size, err := strconv.ParseUint(sizeString, 10, 64)
	if err != nil {
		return nil, err
	}
	return &size, nil
}

func (p ElasticsearchFormatProvider) QueryRecord(testName string, database string, query string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("[%s] QueryRecord not supported for postgresql yet", testName)
}

func NewElasticsearchFormatProvider(runtimeProvider runtime.RuntimeProvider, elasticsearchSnapshotRepository ElasticsearchSnapshotRepository) ElasticsearchFormatProvider {
	elasticsarchFormatProvider := ElasticsearchFormatProvider{
		runtimeProvider: runtimeProvider,
		repository:      elasticsearchSnapshotRepository,
	}
	return elasticsarchFormatProvider
}
