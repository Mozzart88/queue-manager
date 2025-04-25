//go:build integration

package repos_test

import (
	repos "expat-news/queue-manager/internal/repositories"
	"expat-news/queue-manager/internal/repositories/db_test_utils"
	"expat-news/queue-manager/internal/test_utils"
	"expat-news/queue-manager/pkg/utils"
	"strings"
	"testing"
)

func TestGetPublisher(t *testing.T) {
	db_test_utils.SetupDB(t)

	tests := []struct {
		id       *int
		name     *string
		expected *repos.Publisher
	}{
		{
			nil,
			utils.Ptr("perfil"),
			repos.NewPublisher(2, "perfil"),
		},
		{
			utils.Ptr(1),
			nil,
			repos.NewPublisher(1, "pagina12"),
		},
		{
			utils.Ptr(2),
			utils.Ptr("perfil"),
			repos.NewPublisher(2, "perfil"),
		},
		{
			utils.Ptr(1),
			utils.Ptr("perfil"),
			nil,
		},
	}

	for i, test := range tests {
		actual, err := repos.GetPublisher(test.id, test.name)
		if err != nil {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		if test.expected == nil && actual != test.expected {
			test_utils.Fail(t, i, "expected nil, got: %v", actual)
			continue
		}
		if actual == nil {
			continue
		}
		if actual.ID() != test.expected.ID() {
			test_utils.Fail(t, i, "expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
		if actual.Name() != test.expected.Name() {
			test_utils.Fail(t, i, "expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
	}
}

func TestAddPublisher(t *testing.T) {
	db_test_utils.SetupDB(t)

	tests := []struct {
		name     string
		expected *repos.Publisher
	}{
		{
			"some",
			repos.NewPublisher(4, "some"),
		},
	}

	for i, test := range tests {
		newId, err := repos.AddPublisher(test.name)
		if err != nil {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		actual, err := repos.GetPublisher(&newId, nil)
		if err != nil {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		if test.expected == nil && actual != test.expected {
			test_utils.Fail(t, i, "expected nil, got: %v", actual)
			continue
		}
		if actual == nil {
			continue
		}
		if actual.ID() != test.expected.ID() {
			test_utils.Fail(t, i, "expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
		if actual.Name() != test.expected.Name() {
			test_utils.Fail(t, i, "expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
	}
}

func TestUpdatePublisher(t *testing.T) {
	db_test_utils.SetupDB(t)

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

	for i, test := range tests {
		if err := repos.UpdatePublisher(test.id, test.name); err != nil {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		actual, err := repos.GetPublisher(&test.id, &test.name)
		if err != nil {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		if actual == nil {
			test_utils.Fail(t, i, "expected %v, got: nil", *test.expected)
			continue
		}
		if actual.ID() != test.expected.ID() {
			test_utils.Fail(t, i, "expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
		if actual.Name() != test.expected.Name() {
			test_utils.Fail(t, i, "expected\n%v\ngot\n%v", test.expected, actual)
			continue
		}
	}
}

func TestDeletePublisher(t *testing.T) {
	db_test_utils.SetupDB(t)

	tests := []struct {
		id   *int
		name *string
	}{
		{
			nil,
			utils.Ptr("perfil"),
		},
		{
			utils.Ptr(1),
			nil,
		},
		{
			utils.Ptr(2),
			utils.Ptr("perfil"),
		},
		{
			utils.Ptr(1),
			utils.Ptr("perfil"),
		},
	}

	for i, test := range tests {
		if err := repos.DeletePublisher(test.id, test.name); err != nil && !strings.HasPrefix(err.Error(), "unregistered publisher with id") {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		actual, err := repos.GetPublisher(test.id, test.name)
		if err != nil {
			test_utils.Fail(t, i, "error occured: %v", err)
			continue
		}
		if actual != nil {
			test_utils.Fail(t, i, "expected nil, got: %v", actual)
			continue
		}
	}
}

type whantError struct {
	errMsg string
}

func TestDeletePublisher_negative(t *testing.T) {
	db_test_utils.SetupDB(t)

	tests := []struct {
		id         *int
		name       *string
		whantError whantError
	}{
		{
			nil,
			utils.Ptr("perfi"),
			whantError{"unregistered publisher with id: nil and name: perfi"},
		},
		{
			utils.Ptr(256),
			nil,
			whantError{"unregistered publisher with id: 256 and name: nil"},
		},
		{
			nil,
			nil,
			whantError{"empty id and name"},
		},
	}

	for i, test := range tests {
		err := repos.DeletePublisher(test.id, test.name)
		if err == nil {
			test_utils.Fail(t, i, "expected error, got nil", i)
			continue
		}
		if err.Error() != test.whantError.errMsg {
			test_utils.Fail(t, i, "expected error mssage:\n%s\ngot\n%s", test.whantError.errMsg, err.Error())
		}
	}
}
