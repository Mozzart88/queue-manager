//go:guild integration
package db_test

import (
	db "expat-news/queue-manager/internal/db"
	"expat-news/queue-manager/internal/repositories/db_test_utils"
	"expat-news/queue-manager/pkg/utils"
	"reflect"
	"strings"
	"testing"
)

func compare_test_db_message_fields[T any](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return reflect.DeepEqual(*a, *b)
}

func compare_test_db_message(a, b *db.Message) bool {
	if a == nil || b == nil {
		return a == b
	}
	if !compare_test_db_message_fields(a.Id, b.Id) {
		return false
	}
	if !compare_test_db_message_fields(a.Publisher, b.Publisher) {
		return false
	}
	if !compare_test_db_message_fields(a.Msg, b.Msg) {
		return false
	}
	if !compare_test_db_message_fields(a.State, b.State) {
		return false
	}
	return true
}

func TestMessageGet(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(1),
				nil,
				nil,
				nil,
			},
			&db.Message{
				utils.Ptr(1),
				utils.Ptr("pagina12"),
				utils.Ptr("some post from Pagina 12 that already Done"),
				utils.Ptr("done"),
			},
			nil,
		},
		{
			db.Message{
				utils.Ptr(10),
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"no message in queue with specified id:"},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"can't get message: id undefined"},
		},
	}
	for i, test := range tests {
		err := test.m.Get()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
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
func TestMessageSetState(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(1),
				nil,
				nil,
				utils.Ptr("new"),
			},
			&db.Message{
				utils.Ptr(1),
				nil,
				nil,
				utils.Ptr("new"),
			},
			nil,
		},
		{
			db.Message{
				utils.Ptr(10),
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"missing id and/or state"},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				utils.Ptr("new"),
			},
			nil,
			&whantError{"missing id and/or state"},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"missing id and/or state"},
		},
		{
			db.Message{
				utils.Ptr(10),
				nil,
				nil,
				utils.Ptr("done"),
			},
			nil,
			&whantError{"no message with id: "},
		},
	}
	for i, test := range tests {
		err := test.m.SetState()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
				continue
			}
			actual := db.Message{test.m.Id, nil, nil, nil}
			if err := actual.Get(); err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message_fields(actual.State, test.m.State) {
				t.Errorf("test %d: state doesn't updated in database: %v", i, actual)
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
func TestMessageSetNew(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(1),
				nil,
				nil,
				nil,
			},
			&db.Message{
				utils.Ptr(1),
				nil,
				nil,
				utils.Ptr("new"),
			},
			nil,
		},
		{
			db.Message{
				utils.Ptr(10),
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"no message with id: "},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				utils.Ptr("new"),
			},
			nil,
			&whantError{"missing id and/or state"},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"missing id and/or state"},
		},
	}
	for i, test := range tests {
		err := test.m.SetNew()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
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
func TestMessageSetActive(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(1),
				nil,
				nil,
				nil,
			},
			&db.Message{
				utils.Ptr(1),
				nil,
				nil,
				utils.Ptr("active"),
			},
			nil,
		},
	}
	for i, test := range tests {
		err := test.m.SetActive()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
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
func TestMessageRollback(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(1),
				nil,
				nil,
				nil,
			},
			&db.Message{
				utils.Ptr(1),
				nil,
				nil,
				utils.Ptr("new"),
			},
			nil,
		},
	}
	for i, test := range tests {
		err := test.m.Rollback()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
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
func TestMessageSetDone(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(4),
				nil,
				nil,
				nil,
			},
			&db.Message{
				nil,
				nil,
				nil,
				utils.Ptr("done"),
			},
			nil,
		},
		{
			db.Message{
				utils.Ptr(10),
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"no message with id: "},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"can't update state of unsaved message: id is nil"},
		},
	}
	for i, test := range tests {
		actual := db.Message{test.m.Id, nil, nil, nil}
		err := test.m.SetDone()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
				continue
			}
			if err := actual.Get(); err == nil {
				t.Errorf("test %d: state doesn't updated in database: %v", i, actual)
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
func TestMessageAdd(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				nil,
				utils.Ptr("perfil"),
				utils.Ptr("new message"),
				nil,
			},
			&db.Message{
				utils.Ptr(7),
				utils.Ptr("perfil"),
				utils.Ptr("new message"),
				utils.Ptr("new"),
			},
			nil,
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"required fields id and/or publisher are not defined"},
		},
		{
			db.Message{
				nil,
				utils.Ptr("prefil"),
				nil,
				nil,
			},
			nil,
			&whantError{"required fields id and/or publisher are not defined"},
		},
		{
			db.Message{
				nil,
				nil,
				utils.Ptr("prefil"),
				nil,
			},
			nil,
			&whantError{"required fields id and/or publisher are not defined"},
		},
	}
	for i, test := range tests {
		err := test.m.Add()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
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
func TestMessageDelete(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		m          db.Message
		expected   *db.Message
		whantError *whantError
	}{
		{
			db.Message{
				utils.Ptr(4),
				utils.Ptr("perfil"),
				utils.Ptr("some"),
				utils.Ptr("new"),
			},
			&db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
		},
		{
			db.Message{
				utils.Ptr(10),
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"no message with id: "},
		},
		{
			db.Message{
				nil,
				nil,
				nil,
				nil,
			},
			nil,
			&whantError{"can't delete unsaved message"},
		},
	}
	for i, test := range tests {
		actual := db.Message{test.m.Id, nil, nil, nil}
		err := test.m.Delete()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_message(test.expected, &test.m) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.m)
				continue
			}
			if err := actual.Get(); err == nil {
				t.Errorf("test %d: message hasn't deleted: %v", i, actual)
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
