package format

import (
	"encoding/json"
	"fmt"

	"github.com/MaxxtonGroup/backup-validator/pkg/runtime"
)

type MongoFormatProvider struct {
	runtimeProvider runtime.RuntimeProvider
}

type MongoDatabasesResult struct {
	Databases []MongoDatabaseResult `json:"databases"`
}

type MongoDatabaseResult struct {
	Name string `json:"name"`
	Size uint64 `json:"sizeOnDisk"`
}

func (p MongoFormatProvider) Setup(testName string, dir string) error {
	return p.runtimeProvider.Setup(testName, dir)
}

func (p MongoFormatProvider) Destroy(testName string, dir string) error {
	return p.runtimeProvider.Destroy(testName, dir)
}

func (p MongoFormatProvider) ImportData(testName string, dir string, options []string) error {
	_, err := p.runtimeProvider.Exec(testName, "mongorestore", options...)
	return err
}

func (p MongoFormatProvider) GetDatabaseSize(testName string, database string) (*uint64, error) {
	output, err := p.runtimeProvider.Exec(testName, "mongo", "--eval=db.adminCommand( { listDatabases: 1 } )", "--quiet")
	if err != nil {
		return nil, err
	}

	result := MongoDatabasesResult{}
	err = json.Unmarshal([]byte(*output), &result)
	if err != nil {
		return nil, err
	}

	for _, databaseResult := range result.Databases {
		if databaseResult.Name == database {
			return &databaseResult.Size, nil
		}
	}
	return nil, fmt.Errorf("Database %s not found", database)
}

func (p MongoFormatProvider) ListDatabases(testName string) ([]string, error) {
	output, err := p.runtimeProvider.Exec(testName, "mongo", "--eval=db.adminCommand( { listDatabases: 1 } )", "--quiet")
	if err != nil {
		return nil, err
	}

	result := MongoDatabasesResult{}
	err = json.Unmarshal([]byte(*output), &result)
	if err != nil {
		return nil, err
	}

	databaseNames := []string{}
	for _, databaseResult := range result.Databases {
		databaseNames = append(databaseNames, databaseResult.Name)
	}
	return databaseNames, nil
}

func (p MongoFormatProvider) ListTables(testName string, database string) ([]string, error) {
	output, err := p.runtimeProvider.Exec(testName, "mongo", "--eval=db.getCollectionNames()", "--quiet", database)
	if err != nil {
		return nil, err
	}

	result := []string{}
	err = json.Unmarshal([]byte(*output), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p MongoFormatProvider) QueryRecord(testName string, database string, query string) (map[string]interface{}, error) {
	output, err := p.runtimeProvider.Exec(testName, "mongo", "--eval="+query, "--quiet", database)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{}
	err = json.Unmarshal([]byte(*output), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewMongoFormatProvider(runtimeProvider runtime.RuntimeProvider) MongoFormatProvider {
	bongoFormatProvider := MongoFormatProvider{
		runtimeProvider: runtimeProvider,
	}
	return bongoFormatProvider
}
