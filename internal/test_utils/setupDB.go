package test_utils

import (
	repos "expat-news/queue-manager/internal/repositories"
	"testing"

	"fmt"
	"io"
	"os"
	"strings"
)

func loadDBFile(filename string) error {
	file, err := os.Open(filename)
	db := repos.GetDBInstance()
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
	if err := loadDBFile("../../data/schema.sql"); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
	if err := loadDBFile("../../data/test_data.sql"); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
}
