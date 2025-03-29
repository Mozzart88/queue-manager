//go:build integration

package repos_test

import (
	repos "expat-news/queue-manager/internal/repositories"
	"testing"
)

func TestGetPublisher(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id       *int
		name     *string
		expected *repos.Publisher
	}{
		{
			nil,
			ptr("perfil"),
			repos.NewPublisher(2, "perfil"),
		},
		{
			ptr(1),
			nil,
			repos.NewPublisher(1, "pagina12"),
		},
		{
			ptr(2),
			ptr("perfil"),
			repos.NewPublisher(2, "perfil"),
		},
		{
			ptr(1),
			ptr("perfil"),
			nil,
		},
	}

	for _, test := range tests {
		actual, err := repos.GetPublisher(test.id, test.name)
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
		if actual.ID() != test.expected.ID() {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
		if actual.Name() != test.expected.Name() {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
	}
}

func TestAddPublisher(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		name     string
		expected *repos.Publisher
	}{
		{
			"some",
			repos.NewPublisher(4, "some"),
		},
	}

	for _, test := range tests {
		newId, err := repos.AddPublisher(test.name)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		actual, err := repos.GetPublisher(&newId, nil)
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
		if actual.ID() != test.expected.ID() {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
		if actual.Name() != test.expected.Name() {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
	}
}

func TestUpdatePublisher(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id       int
		name     string
		expected *repos.Publisher
	}{
		{
			1,
			"some",
			repos.NewPublisher(1, "some"),
		},
	}

	for _, test := range tests {
		if err := repos.UpdatePublisher(test.id, test.name); err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		actual, err := repos.GetPublisher(&test.id, &test.name)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if actual == nil {
			t.Errorf("expected %v, got: nil", *test.expected)
			continue
		}
		if actual.ID() != test.expected.ID() {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
		if actual.Name() != test.expected.Name() {
			t.Errorf("expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
	}
}

func TestDeletePublisher(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id   *int
		name *string
	}{
		{
			nil,
			ptr("perfil"),
		},
		{
			ptr(1),
			nil,
		},
		{
			ptr(2),
			ptr("perfil"),
		},
		{
			ptr(1),
			ptr("perfil"),
		},
	}

	for _, test := range tests {
		if err := repos.DeletePublisher(test.id, test.name); err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		actual, err := repos.GetPublisher(test.id, test.name)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if actual != nil {
			t.Errorf("expected nil, got: %v", actual)
			continue
		}
	}
}
