//go:build integration

package repos_test

import (
	repos "expat-news/queue-manager/internal/repositories"
	"testing"
)

func TestAddMessage(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		publisherId int
		content     string
		expected    repos.QueueMessage
	}{
		{
			3,
			"some new message",
			*repos.NewQueueMessage(
				7,
				"some new message",
				3,
				"lapolitica",
				1,
				"new",
			),
		},
	}
	for i, test := range tests {
		newId, err := repos.AddMessage(test.content, test.publisherId)
		if err != nil {
			t.Errorf("test %d AddMessage: error occured on %v", i, err)
			continue
		}
		actual, err := repos.GetUniqQueueMessage(newId)
		if err != nil {
			t.Errorf("test %d GetUniqQueueMessage: error occured on %v", i, err)
			continue
		}
		if actual == nil {
			t.Errorf("test %d: fail to get new message %v", i, test.expected)
			continue
		}
		if !compare_test_queue_message(actual, &test.expected) {
			t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, actual)
		}
	}
}
func TestAddMessages(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		publisherId int
		content     []string
		expected    int
	}{
		{
			3,
			[]string{"some new message"},
			1,
		},
		{
			3,
			[]string{"second new message message", "third new message"},
			2,
		},
	}
	for i, test := range tests {
		added, err := repos.AddMessages(test.publisherId, &test.content)
		if err != nil {
			t.Errorf("test %d AddMessages: error occured on %v", i, err)
			continue
		}
		if test.expected != added {
			t.Errorf("test %d AddMessages: expected %d, got %d", i, test.expected, added)
		}
	}
}

func TestUpdateStateMessage(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id       int
		state    repos.State_t
		expected repos.QueueMessage
	}{
		{
			1,
			repos.STATE_NEW,
			*repos.NewQueueMessage(
				1,
				"some post from Pagina 12 that already Done",
				1,
				"pagina12",
				1,
				"new",
			),
		},
	}
	for i, test := range tests {
		if err := repos.UpdateStateMessage(test.id, test.state); err != nil {
			t.Errorf("test %d AddMessage: error occured on %v", i, err)
			continue
		}
		actual, err := repos.GetUniqQueueMessage(test.id)
		if err != nil {
			t.Errorf("test %d GetUniqQueueMessage: error occured on %v", i, err)
			continue
		}
		if actual == nil {
			t.Errorf("test %d: fail to get new message %v", i, test.expected)
			continue
		}
		if !compare_test_queue_message(actual, &test.expected) {
			t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, actual)
		}
	}
}
func TestDeleteMessage(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		id int
	}{
		{
			3,
		},
	}
	for i, test := range tests {
		if err := repos.DeleteMessage(test.id); err != nil {
			t.Errorf("test %d AddMessage: error occured on %v", i, err)
			continue
		}
		actual, err := repos.GetUniqQueueMessage(test.id)
		if err != nil {
			t.Errorf("test %d GetUniqQueueMessage: error occured on %v", i, err)
			continue
		}
		if actual != nil {
			t.Errorf("test %d: message with given id (%d) still in base: %v", i, test.id, actual)
		}
	}

}
