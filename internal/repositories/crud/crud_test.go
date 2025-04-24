//_go:build integration

package crud

import (
	"expat-news/queue-manager/internal/test_utils"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func loadDBFile(filename string) error {
	file, err := os.Open(filename)
	db := GetDBInstance()
	if err != nil {
		return fmt.Errorf("failed to open SQL file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read SQL file: %w", err)
	}

	queries := strings.SplitSeq(string(content), ";")
	for query := range queries {
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
	if err := loadDBFile("../../../data/schema.sql"); err != nil {
		return err
	}
	if err := loadDBFile("../../../data/test_data.sql"); err != nil {
		return err
	}
	return nil
}

func TestExecSql(t *testing.T) {
	os.Setenv("QDB_FILE", "../../../data/data.db")
	_, err := execSql("SELECT * FROM message_status")
	if err != nil {
		t.Fatal("Query failed: ", err)
	}

}

func compare_test_get_results(a, b *QueryRow) bool {
	if len(*a) != len(*b) {
		return false
	}
	for key, valA := range *a {
		valB, exists := (*b)[key]
		if !exists || valA != valB {
			return false
		}
	}
	return true
}

func TestWhereNew(t *testing.T) {
	where := Where.New(Where{})
	if where.Statements == nil {
		t.Error("statments equals nil")
	} else if where.Union != U_Empty {
		t.Errorf("union value expected be empty string, got %v", where.Union)
	}
}

func TestWhereUnion(t *testing.T) {
	where := Where.New(Where{})
	if where.Union != U_Empty {
		t.Errorf("union value expected be empty, got %v", where.Union.String())
	}
	where = Where.New(Where{Union: U_Or})
	if where.Union != U_Or {
		t.Errorf("union value expected be 'OR', got %v", where.Union.String())
	}
}

func TestWhereAdd(t *testing.T) {
	where := Where.New(Where{})
	where.Add("some", 1, Equals)
	val, ok := where.Statements["some"]
	if !ok {
		t.Error("key 'some' not added to the statments")
	} else if val.Value != 1 {
		t.Errorf("statement value expected = 1, got %v", val.Value)
	}
}

func TestWhereEquals(t *testing.T) {
	where := Where.New(Where{})
	where.Equals("some", 1)
	val, ok := where.Statements["some"]
	if !ok {
		t.Error("key 'some' not added to the statments")
	} else if val.Value != 1 {
		t.Errorf("statement value expected = 1, got %v", val.Value)
	}
}

func TestWhereLen(t *testing.T) {
	where := Where.New(Where{})
	where.Equals("some", 1)
	if where.Len() != 1 {
		t.Errorf("len expected equals 1, got %d", where.Len())
		return
	}
	where.Equals("other", 1)
	if where.Len() != 2 {
		t.Errorf("len expected equals 2, got %d", where.Len())
		return
	}
}

func TestGet(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
	var l Limit = 1
	tests := []struct {
		w        *Where
		f        *Fields
		l        *Limit
		o        *Order
		expected *[]QueryRow
	}{
		{nil, nil, nil, nil, &[]QueryRow{
			{
				"id":           {int64(1)},
				"status_id":    {int64(0)},
				"publisher_id": {int64(1)},
				"content":      {"some post from Pagina 12 that already Done"},
			},
			{
				"id":           {int64(2)},
				"status_id":    {int64(2)},
				"publisher_id": {int64(3)},
				"content":      {"some post from La Politica that processing right now"},
			},
			{
				"id":           {int64(3)},
				"status_id":    {int64(1)},
				"publisher_id": {int64(1)},
				"content":      {"some post from Pagina 12"},
			},
			{
				"id":           {int64(4)},
				"status_id":    {int64(1)},
				"publisher_id": {int64(1)},
				"content":      {"some other post from Pagina 12"},
			},
			{
				"id":           {int64(5)},
				"status_id":    {int64(1)},
				"publisher_id": {int64(2)},
				"content":      {"some post from Perfil"},
			},
			{
				"id":           {int64(6)},
				"status_id":    {int64(1)},
				"publisher_id": {int64(3)},
				"content":      {"some post from La Politica"},
			},
		}},
		{nil, nil, &l, nil, &[]QueryRow{{
			"id":           {int64(1)},
			"status_id":    {int64(0)},
			"publisher_id": {int64(1)},
			"content":      {"some post from Pagina 12 that already Done"},
		}}},
		{nil, nil, &l, &Order{[]string{"id"}, "DESC"}, &[]QueryRow{{
			"id":           {int64(6)},
			"status_id":    {int64(1)},
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}}},
		{&Where{map[string]Statement{"id": {6, Equals}}, U_Empty}, &Fields{"publisher_id", "content"}, nil, nil, &[]QueryRow{{
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}}},
		{&Where{map[string]Statement{"id": {7, Equals}}, U_Empty}, &Fields{"publisher_id", "content"}, nil, nil, &[]QueryRow{}},
	}

	for _, test := range tests {
		res, err := Get("message", test.f, test.w, test.o, test.l)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		for i, m := range res {
			if !compare_test_get_results(&m, &(*test.expected)[i]) {
				t.Errorf("actual != expected:\n%v\n%v", res, test.expected)
				break
			}

		}
	}
}

func TestGetOne(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
	tests := []struct {
		w        *Where
		f        *Fields
		o        *Order
		expected *QueryRow
	}{
		{nil, nil, nil, &QueryRow{
			"id":           {int64(1)},
			"status_id":    {int64(0)},
			"publisher_id": {int64(1)},
			"content":      {"some post from Pagina 12 that already Done"},
		}},
		{nil, nil, nil, &QueryRow{
			"id":           {int64(1)},
			"status_id":    {int64(0)},
			"publisher_id": {int64(1)},
			"content":      {"some post from Pagina 12 that already Done"},
		}},
		{nil, nil, &Order{[]string{"id"}, "DESC"}, &QueryRow{
			"id":           {int64(6)},
			"status_id":    {int64(1)},
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}},
		{&Where{map[string]Statement{"id": {6, Equals}}, U_Empty}, &Fields{"publisher_id", "content"}, nil, &QueryRow{
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}},
		{&Where{map[string]Statement{"id": {7, Equals}}, U_Empty}, &Fields{"publisher_id", "content"}, nil, nil},
	}

	for _, test := range tests {
		res, err := GetOne("message", test.f, test.w, test.o)
		if err != nil {
			t.Errorf("error occured: %v", err)
			continue
		}
		if test.expected != nil && res == nil {
			t.Errorf("expected nil, got :\n%v", res)
			continue
		}
		if test.expected != nil && !compare_test_get_results(&res, test.expected) {
			t.Errorf("actual != expected:\n%v\n%v", res, test.expected)
		}
	}
}

func TestDelete(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	type deleteQ struct {
		w        *Where
		expected int64
	}
	type getQ struct {
		w        *Where
		expected QueryRow
	}
	tests := []struct {
		delete deleteQ
		get    getQ
	}{
		{
			deleteQ{
				&Where{map[string]Statement{"id": {1, Equals}}, U_Empty},
				1,
			},
			getQ{
				&Where{map[string]Statement{"id": {1, Equals}}, U_Empty},
				nil,
			},
		},
	}

	for i, test := range tests {
		{
			res, err := Delete("message", test.delete.w)
			if err != nil {
				test_utils.Fail(t, i, "unexpected error occured: %v", err)
				continue
			}
			if test.delete.expected != res {
				test_utils.Fail(t, i, "expected to delete  %d row, got %d", test.delete.expected, res)
				continue
			}
		}
		{
			res, err := GetOne("message", nil, test.get.w, nil)
			if err != nil {
				test_utils.Fail(t, i, "unexpected error occured: %v", err)
				continue
			}
			if res != nil {
				test_utils.Fail(t, i, "get should returns nil, got %v", res)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	type updateQ struct {
		w        *Where
		f        *Fields
		expected int
	}
	type getQ struct {
		w        *Where
		f        *Fields
		expected QueryRow
	}
	tests := []struct {
		update updateQ
		get    getQ
	}{
		{
			updateQ{
				&Where{map[string]Statement{"id": {1, Equals}}, U_Empty},
				&Fields{"status_id = 2"},
				1,
			},
			getQ{
				&Where{map[string]Statement{"id": {1, Equals}}, U_Empty},
				&Fields{"status_id"},
				QueryRow{
					"status_id": {int64(2)},
				},
			},
		},
		{
			updateQ{
				&Where{map[string]Statement{"id": {1, Equals}}, U_Empty},
				&Fields{"status_id = 1", "publisher_id = 3"},
				1,
			},
			getQ{
				&Where{map[string]Statement{"id": {1, Equals}}, U_Empty},
				&Fields{"status_id", "publisher_id"},
				QueryRow{
					"status_id":    {int64(1)},
					"publisher_id": {int64(3)},
				},
			},
		},
	}

	for _, test := range tests {
		{
			res, err := Update("message", test.update.f, test.update.w)
			if err != nil {
				t.Errorf("unexpected error occured: %v", err)
				continue
			}
			if test.update.expected != res {
				t.Logf("expected to delete  %d row, got %d", test.update.expected, res)
				continue
			}
		}
		{
			res, err := GetOne("message", test.get.f, test.get.w, nil)
			if err != nil {
				t.Errorf("unexpected error occured: %v", err)
				continue
			}
			if !compare_test_get_results(&res, &test.get.expected) {
				t.Errorf("get result not updated - expected:\n%v\ngot %v", test.get.expected, res)
			}
		}
	}
}

func TestInsert(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	type insertQ struct {
		f *Fields
		v Values
	}
	type getQ struct {
		f        *Fields
		expected QueryRow
	}
	tests := []struct {
		insert insertQ
		get    getQ
	}{
		{
			insertQ{
				&Fields{"name"},
				Values{"some"},
			},
			getQ{
				&Fields{"name"},
				QueryRow{
					"name": {"some"},
				},
			},
		},
		{
			insertQ{
				nil,
				Values{256, "some"},
			},
			getQ{
				&Fields{"id", "name"},
				QueryRow{
					"id":   {int64(256)},
					"name": {"some"},
				},
			},
		},
	}

	for i, test := range tests {
		id, err := Insert("publisher", test.insert.f, &test.insert.v)
		if err != nil {
			test_utils.Fail(t, i, "unexpected error occured: %v", err)
			continue
		}
		res, err := GetOne("publisher", test.get.f, &Where{
			map[string]Statement{"id": {id, Equals}},
			U_Empty,
		}, nil)
		if err != nil {
			test_utils.Fail(t, i, "unexpected error occured: %v", err)
			continue
		}
		if !compare_test_get_results(&res, &test.get.expected) {
			test_utils.Fail(t, i, "get result not inserted - expected:\n%v\ngot %v", test.get.expected, res)
		}
	}
}
func TestInsertMany(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		f        *Fields
		v        []Values
		expected int
	}{
		{
			&Fields{"name"},
			[]Values{{"some"}, {"other"}},
			2,
		},
	}

	for _, test := range tests {
		res, err := InsertMany("publisher", test.f, &test.v)
		if err != nil {
			t.Errorf("unexpected error occured: %v", err)
			continue
		}
		if res != test.expected {
			t.Logf("expected to be added %d row, got %d", test.expected, res)
			continue
		}
	}
}
