//go:build integration

package db_test_utils

import (
	"expat-news/queue-manager/internal/repositories/crud"
	"testing"

	"fmt"
	"io"
	"os"
	"strings"
)

func loadDBFile(filename string) error {
	file, err := os.Open(filename)
	db := crud.GetDBInstance()
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to execute query: %s, error: %w", query, err)
		}
	}

	return nil
}

func SetupDB(t *testing.T) {
	const env = "QDB_DIR"
	dataDir, exist := os.LookupEnv(env)
	if !exist || len(dataDir) == 0 {
		t.Fatal("undefined QDB_DIR env variable")
	}
	if err := loadDBFile(dataDir + "/schema.sql"); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
	if err := loadDBFile(dataDir + "/test_data.sql"); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
}
