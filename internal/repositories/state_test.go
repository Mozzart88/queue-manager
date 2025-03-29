//go:build integration

package repos_test

import (
	repos "expat-news/queue-manager/internal/repositories"
	"testing"
)

func TestGetState(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id       *int
		name     *repos.State_t
		expected *repos.State
	}{
		{
			nil,
			ptr(repos.STATE_NEW),
			repos.NewState(1, repos.STATE_NEW),
		},
		{
			ptr(0),
			nil,
			repos.NewState(0, repos.STATE_DONE),
		},
		{
			ptr(2),
			ptr(repos.STATE_ACTIVE),
			repos.NewState(2, repos.STATE_ACTIVE),
		},
		{
			ptr(0),
			ptr(repos.STATE_ACTIVE),
			nil,
		},
	}

	for _, test := range tests {
		actual, err := repos.GetState(test.id, test.name)
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
