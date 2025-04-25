//go:guild integration

package db_test

import (
	db "expat-news/queue-manager/internal/db"
	"expat-news/queue-manager/internal/repositories/db_test_utils"
	"expat-news/queue-manager/pkg/utils"
	"strings"
	"testing"
)

func compare_test_db_publishers(a, b *db.Publisher) bool {
	if a == nil && a != b {
		return false
	}
	if *a.Id != *b.Id {
		return false
	}
	if *a.Name != *b.Name {
		return false
	}
	return true
}

type whantError struct {
	msg string
}

func TestPublisherGet(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		p          db.Publisher
		expected   *db.Publisher
		whantError *whantError
	}{
		{
			db.Publisher{utils.Ptr(1), nil},
			&db.Publisher{utils.Ptr(1), utils.Ptr("pagina12")},
			nil,
		},
		{
			db.Publisher{nil, utils.Ptr("perfil")},
			&db.Publisher{utils.Ptr(2), utils.Ptr("perfil")},
			nil,
		},
		{
			db.Publisher{utils.Ptr(3), utils.Ptr("lapolitica")},
			&db.Publisher{utils.Ptr(3), utils.Ptr("lapolitica")},
			nil,
		},
		{
			db.Publisher{nil, nil},
			nil,
			&whantError{"empty id and name"},
		},
		{
			db.Publisher{utils.Ptr(4), nil},
			nil,
			&whantError{"unregistered publisher"},
		},
		{
			db.Publisher{nil, utils.Ptr("some")},
			nil,
			&whantError{"unregistered publisher"},
		},
	}
	for i, test := range tests {
		err := test.p.Get()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_publishers(test.expected, &test.p) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.p)
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

func TestPublisherRegister(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		p          db.Publisher
		expected   *db.Publisher
		whantError *whantError
	}{
		{
			db.Publisher{utils.Ptr(1), nil},
			&db.Publisher{utils.Ptr(1), utils.Ptr("pagina12")},
			&whantError{"name is undefined"},
		},
		{
			db.Publisher{nil, utils.Ptr("some")},
			&db.Publisher{utils.Ptr(4), utils.Ptr("some")},
			nil,
		},
		{
			db.Publisher{utils.Ptr(9), utils.Ptr("other")},
			&db.Publisher{utils.Ptr(5), utils.Ptr("other")},
			nil,
		},
	}
	for i, test := range tests {
		err := test.p.Register()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_publishers(test.expected, &test.p) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.p)
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
func TestPublisherDelete(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		p          db.Publisher
		expected   *db.Publisher
		whantError *whantError
	}{
		{
			db.Publisher{nil, nil},
			nil,
			&whantError{"empty id and name"},
		},
		{
			db.Publisher{nil, utils.Ptr("some")},
			nil,
			&whantError{"unregistered publisher with id"},
		},
		{
			db.Publisher{utils.Ptr(10), utils.Ptr("some")},
			nil,
			&whantError{"unregistered publisher with id"},
		},
		{
			db.Publisher{utils.Ptr(1), nil},
			nil,
			nil,
		},
	}
	for i, test := range tests {
		err := test.p.Delete()
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if test.p.Id != nil || test.p.Name != nil {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.p)
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
func TestPublisherUpdate(t *testing.T) {
	db_test_utils.SetupDB(t)
	tests := []struct {
		p          db.Publisher
		expected   *db.Publisher
		whantError *whantError
	}{
		{
			db.Publisher{nil, nil},
			nil,
			&whantError{"id is undefined"},
		},
		{
			db.Publisher{utils.Ptr(10), nil},
			nil,
			&whantError{"unregistered publisher with id"},
		},
		{
			db.Publisher{utils.Ptr(1), nil},
			&db.Publisher{utils.Ptr(1), utils.Ptr("some")},
			nil,
		},
	}
	newName := "some"
	for i, test := range tests {
		err := test.p.Update(newName)
		if test.whantError == nil {
			if err != nil {
				t.Errorf("test %d: unexpected error occured: %v", i, err)
				continue
			}
			if !compare_test_db_publishers(&test.p, test.expected) {
				t.Errorf("test %d: expected\n%v\nactual\n%v", i, test.expected, test.p)
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
