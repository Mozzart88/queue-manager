//go:guild integration

package db_test

import (
	db "expat-news/queue-manager/internal/db"
	"expat-news/queue-manager/internal/test_utils"
	"expat-news/queue-manager/pkg/utils"
	"strings"
	"testing"
)

func TestQueueGetMessage(t *testing.T) {
	test_utils.SetupDB(t)
	tests := []struct {
		q          db.Queue
		o          *bool
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Queue{"pagina12", "new"},
			nil,
			&db.Message{
				utils.Ptr(3),
				utils.Ptr("pagina12"),
				utils.Ptr("some post from Pagina 12"),
				utils.Ptr("new"),
			},
			nil,
		},
		{
			db.Queue{"pagina12", "new"},
			utils.Ptr(false),
			&db.Message{
				utils.Ptr(4),
				utils.Ptr("pagina12"),
				utils.Ptr("some other post from Pagina 12"),
				utils.Ptr("new"),
			},
			nil,
		},
		{
			db.Queue{"perfil", "done"},
			utils.Ptr(false),
			nil,
			nil,
		},
	}
	for i, test := range tests {
		actual, err := test.q.GetMessage(test.o)
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, actual) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, actual)
				continue
			}
		} else {
			if err == nil {
				t.Errorf("test %d: expected error, got nil", i)
				continue
			}
			if !strings.HasPrefix(err.Error(), test.whantError.msg) && err.Error() != test.whantError.msg {
				t.Errorf("test %d: expected error\n%v\ngot\n%v", i, test.whantError.msg, err.Error())
				continue
			}
		}
	}
}
func TestQueueGetMessages(t *testing.T) {
	test_utils.SetupDB(t)
	tests := []struct {
		q          db.Queue
		o          *bool
		expected   *[]db.Message
		whantError *whantError
	}{
		{
			db.Queue{"pagina12", "new"},
			nil,
			&[]db.Message{
				{
					utils.Ptr(3),
					utils.Ptr("pagina12"),
					utils.Ptr("some post from Pagina 12"),
					utils.Ptr("new"),
				},
				{
					utils.Ptr(4),
					utils.Ptr("pagina12"),
					utils.Ptr("some other post from Pagina 12"),
					utils.Ptr("new"),
				},
			},
			nil,
		},
		{
			db.Queue{"pagina12", "new"},
			utils.Ptr(false),
			&[]db.Message{
				{
					utils.Ptr(4),
					utils.Ptr("pagina12"),
					utils.Ptr("some other post from Pagina 12"),
					utils.Ptr("new"),
				},
				{
					utils.Ptr(3),
					utils.Ptr("pagina12"),
					utils.Ptr("some post from Pagina 12"),
					utils.Ptr("new"),
				},
			},
			nil,
		},
		{
			db.Queue{"perfil", "done"},
			utils.Ptr(false),
			&[]db.Message{},
			nil,
		},
	}
	for i, test := range tests {
		msgs, err := test.q.GetMessages(test.o)
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if len(msgs) != len(*test.expected) {
				t.Errorf("test %d: expected len %d, got %d", i, len(*test.expected), len(msgs))
				continue
			}
			for ind, actual := range msgs {
				if !compare_test_db_message(&(*test.expected)[ind], &actual) {
					t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, actual)
					continue
				}
			}
		} else {
			if err == nil {
				t.Errorf("test %d: expected error, got nil", i)
				continue
			}
			if !strings.HasPrefix(err.Error(), test.whantError.msg) && err.Error() != test.whantError.msg {
				t.Errorf("test %d: expected error\n%v\ngot\n%v", i, test.whantError.msg, err.Error())
				continue
			}
		}
	}
}
func TestQueueAddMessages(t *testing.T) {
	test_utils.SetupDB(t)
	tests := []struct {
		q          db.Queue
		msgs       []string
		expected   *int
		whantError *whantError
	}{
		{
			db.Queue{"perfil", ""},
			[]string{"some new message", "some other new message"},
			utils.Ptr(2),
			nil,
		},
		{
			db.Queue{"some", ""},
			[]string{"some new message", "some other new message"},
			nil,
			&whantError{"unregistered publisher:"},
		},
	}
	for i, test := range tests {
		actual, err := test.q.AddMessages(&test.msgs)
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if len(test.msgs) != actual {
				t.Errorf("test %d: expected len %d, got %d", i, test.expected, actual)
				continue
			}
		} else {
			if err == nil {
				t.Errorf("test %d: expected error, got nil", i)
				continue
			}
			if !strings.HasPrefix(err.Error(), test.whantError.msg) && err.Error() != test.whantError.msg {
				t.Errorf("test %d: expected error\n%v\ngot\n%v", i, test.whantError.msg, err.Error())
				continue
			}
		}
	}
}
