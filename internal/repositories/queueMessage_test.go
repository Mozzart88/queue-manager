//go:build integration

package repos_test

import (
	repos "expat-news/queue-manager/internal/repositories"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
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

func setupDB() error {
	if err := loadDBFile("../../data/schema.sql"); err != nil {
		return err
	}
	if err := loadDBFile("../../data/test_data.sql"); err != nil {
		return err
	}
	return nil
}

func ptr[T any](v T) *T {
	return &v
}

func compare_test_queue_message(a, b *repos.QueueMessage) bool {
	if a == nil && a != nil {
		return false
	}
	if a.Content() != b.Content() {
		return false
	}
	if a.ID() != b.ID() {
		return false
	}
	if a.PublisherId() != b.PublisherId() {
		return false
	}
	if a.PublisherName() != b.PublisherName() {
		return false
	}
	if a.StateId() != b.StateId() {
		return false
	}
	if a.StateName() != b.StateName() {
		return false
	}
	return true
}

func TestGetUniqQueueMessage(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id       int
		expected *repos.QueueMessage
	}{
		{
			2,
			repos.NewQueueMessage(
				2,
				"some post from La Politica that processing right now",
				3,
				"lapolitica",
				2,
				"active",
			),
		},
		{
			256,
			nil,
		},
	}

	for _, test := range tests {
		actual, err := repos.GetUniqQueueMessage(test.id)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if test.expected == nil && actual != test.expected {
			t.Errorf("expected nil, got: %v", actual)
			continue
		}
		if actual == nil {
			continue
		}
		if !compare_test_queue_message(actual, test.expected) {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
		}
	}
}

func TestGetQueueMessage(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		publisher string
		state     repos.State_t
		oldest    *bool
		expected  *repos.QueueMessage
	}{
		{
			"pagina12",
			repos.STATE_NEW,
			nil,
			repos.NewQueueMessage(
				3,
				"some post from Pagina 12",
				1,
				"pagina12",
				1,
				"new",
			),
		},
		{
			"pagina12",
			repos.STATE_NEW,
			ptr(true),
			repos.NewQueueMessage(
				3,
				"some post from Pagina 12",
				1,
				"pagina12",
				1,
				"new",
			),
		},
		{
			"pagina12",
			repos.STATE_NEW,
			ptr(false),
			repos.NewQueueMessage(
				4,
				"some other post from Pagina 12",
				1,
				"pagina12",
				1,
				"new",
			),
		},
		{
			"pagina12",
			repos.STATE_ACTIVE,
			nil,
			nil,
		},
	}

	for i, test := range tests {
		actual, err := repos.GetQueueMessage(test.publisher, test.state, test.oldest)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if test.expected == nil && actual != test.expected {
			t.Errorf("expected nil, got: %v", actual)
			continue
		}
		if actual == nil {
			continue
		}
		if !compare_test_queue_message(actual, test.expected) {
			t.Errorf("for test case %d expected\n%v\ngot\n%v", i, test.expected, actual)
		}
	}
}

func TestGetMessages(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		publisher string
		state     repos.State_t
		oldest    *bool
		expected  int
	}{
		{
			"pagina12",
			repos.STATE_NEW,
			nil,
			2,
		},
		{
			"pagina12",
			repos.STATE_ACTIVE,
			nil,
			0,
		},
	}

	for _, test := range tests {
		actual, err := repos.GetMessages(test.publisher, test.state, test.oldest)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if test.expected != len(actual) {
			t.Errorf("expected len %d, got: %v", test.expected, len(actual))
			continue
		}
	}
}

func TestAddQueueMessage(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		content   string
		publisher string
		expected  *repos.QueueMessage
	}{
		{
			"some new message",
			"perfil",
			repos.NewQueueMessage(
				7,
				"some new message",
				2,
				"perfil",
				1,
				"new",
			),
		},
	}

	for _, test := range tests {
		actual, err := repos.AddQueueMessage(test.content, test.publisher)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if actual == nil {
			t.Errorf("expected new message, got: %v", actual)
			continue
		}
		if !compare_test_queue_message(actual, test.expected) {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
		}
	}
}

func TestAddQueueMessage_negative(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		content   string
		publisher string
		expected  string
	}{
		{
			"some new message",
			"unregistered",
			"unregistered publisher: unregistered",
		},
	}

	for _, test := range tests {
		_, err := repos.AddQueueMessage(test.content, test.publisher)
		if err == nil {
			t.Errorf("expected error got nil")
			continue
		}
		if err.Error() != test.expected {
			t.Errorf("expected error '%s' got '%v'", test.expected, err)
			continue
		}
	}
}
