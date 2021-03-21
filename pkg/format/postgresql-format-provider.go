package format

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/MaxxtonGroup/backup-validator/pkg/runtime"
)

type PostgresqlFormatProvider struct {
	runtimeProvider runtime.RuntimeProvider
}

type PostgresqlDatabasesResult struct {
	Databases []PostgresqlDatabaseResult `json:"databases"`
}

type PostgresqlDatabaseResult struct {
	Name string `json:"name"`
	Size uint64 `json:"sizeOnDisk"`
}

func (p PostgresqlFormatProvider) Setup(dir string) error {
	return p.runtimeProvider.Setup(dir)
}

func (p PostgresqlFormatProvider) Destroy(dir string) error {
	return p.runtimeProvider.Destroy(dir)
}

func (p PostgresqlFormatProvider) ImportData(dir string, options []string) error {
	_, err := p.runtimeProvider.Exec("pg_restore", options...)
	log.Println("Import complete")
	return err
}

func (p PostgresqlFormatProvider) GetDatabaseSize(database string) (*uint64, error) {
	psqlUser, err := p.getPostgresUser()
	if err != nil {
		return nil, err
	}
	psqlDatabase, err := p.getPostgresDatabase()
	if err != nil {
		return nil, err
	}

	output, err := p.runtimeProvider.Exec("psql", "--username="+*psqlUser, *psqlDatabase, "-t", "-c", "select pg_database_size('"+database+"');")
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

func (p PostgresqlFormatProvider) ListDatabases() ([]string, error) {
	psqlUser, err := p.getPostgresUser()
	if err != nil {
		return nil, err
	}
	psqlDatabase, err := p.getPostgresDatabase()
	if err != nil {
		return nil, err
	}

	output, err := p.runtimeProvider.Exec("psql", "--username="+*psqlUser, *psqlDatabase, "-t", "-c", "select datname from pg_database;")
	if err != nil {
		return nil, err
	}
	databaseNames := []string{}
	databases := strings.Split(*output, "\n")
	for _, database := range databases {
		dbName := strings.TrimSpace(database)
		if dbName != "posgres" && dbName != "" {
			databaseNames = append(databaseNames, dbName)
		}
	}

	return databaseNames, nil
}

func (p PostgresqlFormatProvider) ListTables(database string) ([]string, error) {
	psqlUser, err := p.getPostgresUser()
	if err != nil {
		return nil, err
	}

	output, err := p.runtimeProvider.Exec("psql", "--username="+*psqlUser, database, "-t", "-c", "SELECT table_name FROM information_schema.tables WHERE table_catalog='"+database+"' AND table_type='BASE TABLE';")
	if err != nil {
		return nil, err
	}
	tableNames := []string{}
	tables := strings.Split(*output, "\n")
	for _, table := range tables {
		tableName := strings.TrimSpace(table)
		if tableName != "" {
			tableNames = append(tableNames, tableName)
		}
	}

	return tableNames, nil
}

func (p PostgresqlFormatProvider) QueryRecord(database string, query string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("QueryRecord not supported for postgresql yet")
}

func (p PostgresqlFormatProvider) getPostgresUser() (*string, error) {
	envs, err := p.runtimeProvider.Exec("env")
	if err != nil {
		return nil, err
	}
	envList := strings.Split(*envs, "\n")
	psqlUser := "postgres"
	for _, env := range envList {
		if strings.HasPrefix(env, "POSTGRES_USER=") {
			psqlUser = strings.TrimPrefix(env, "POSTGRES_USER=")
			break
		}
	}
	return &psqlUser, nil
}

func (p PostgresqlFormatProvider) getPostgresDatabase() (*string, error) {
	envs, err := p.runtimeProvider.Exec("env")
	if err != nil {
		return nil, err
	}
	envList := strings.Split(*envs, "\n")
	psqlUser := "postgres"
	for _, env := range envList {
		if strings.HasPrefix(env, "POSTGRES_DB=") {
			psqlUser = strings.TrimPrefix(env, "POSTGRES_DB=")
			break
		}
	}
	return &psqlUser, nil
}

func NewPostgresqlFormatProvider(runtimeProvider runtime.RuntimeProvider) PostgresqlFormatProvider {
	bongoFormatProvider := PostgresqlFormatProvider{
		runtimeProvider: runtimeProvider,
	}
	return bongoFormatProvider
}
