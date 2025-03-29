//go:build integration

package repos

import (
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

func TestExecSql(t *testing.T) {
	os.Setenv("QDB_FILE", "../../data/data.db")
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

func TestGet(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}
	var l limit = 1
	tests := []struct {
		w        *where
		f        *fields
		l        *limit
		o        *order
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
		{nil, nil, &l, &order{[]string{"id"}, "DESC"}, &[]QueryRow{{
			"id":           {int64(6)},
			"status_id":    {int64(1)},
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}}},
		{&where{[]string{"id = 6"}, ""}, &fields{"publisher_id", "content"}, nil, nil, &[]QueryRow{{
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}}},
		{&where{[]string{"id = 7"}, ""}, &fields{"publisher_id", "content"}, nil, nil, &[]QueryRow{}},
	}

	for _, test := range tests {
		res, err := get("message", test.f, test.w, test.o, test.l)
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
		w        *where
		f        *fields
		o        *order
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
		{nil, nil, &order{[]string{"id"}, "DESC"}, &QueryRow{
			"id":           {int64(6)},
			"status_id":    {int64(1)},
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}},
		{&where{[]string{"id = 6"}, ""}, &fields{"publisher_id", "content"}, nil, &QueryRow{
			"publisher_id": {int64(3)},
			"content":      {"some post from La Politica"},
		}},
		{&where{[]string{"id = 7"}, ""}, &fields{"publisher_id", "content"}, nil, nil},
	}

	for _, test := range tests {
		res, err := getOne("message", test.f, test.w, test.o)
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
		w        *where
		expected int64
	}
	type getQ struct {
		w        *where
		expected QueryRow
	}
	tests := []struct {
		delete deleteQ
		get    getQ
	}{
		{
			deleteQ{
				&where{[]string{"id = 1"}, ""},
				1,
			},
			getQ{
				&where{[]string{"id = 1"}, ""},
				nil,
			},
		},
	}

	for _, test := range tests {
		{
			res, err := delete("message", test.delete.w)
			if err != nil {
				t.Errorf("unexpected error occured: %v", err)
				continue
			}
			if test.delete.expected != res {
				t.Logf("expected to delete  %d row, got %d", test.delete.expected, res)
				continue
			}
		}
		{
			res, err := getOne("message", nil, test.get.w, nil)
			if err != nil {
				t.Errorf("unexpected error occured: %v", err)
				continue
			}
			if res != nil {
				t.Errorf("get should returns nil, got %v", res)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	type updateQ struct {
		w        *where
		f        *fields
		expected int
	}
	type getQ struct {
		w        *where
		f        *fields
		expected QueryRow
	}
	tests := []struct {
		update updateQ
		get    getQ
	}{
		{
			updateQ{
				&where{[]string{"id = 1"}, ""},
				&fields{"status_id = 2"},
				1,
			},
			getQ{
				&where{[]string{"id = 1"}, ""},
				&fields{"status_id"},
				QueryRow{
					"status_id": {int64(2)},
				},
			},
		},
		{
			updateQ{
				&where{[]string{"id = 1"}, ""},
				&fields{"status_id = 1", "publisher_id = 3"},
				1,
			},
			getQ{
				&where{[]string{"id = 1"}, ""},
				&fields{"status_id", "publisher_id"},
				QueryRow{
					"status_id":    {int64(1)},
					"publisher_id": {int64(3)},
				},
			},
		},
	}

	for _, test := range tests {
		{
			res, err := update("message", test.update.f, test.update.w)
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
			res, err := getOne("message", test.get.f, test.get.w, nil)
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
		f *fields
		v values
	}
	type getQ struct {
		f        *fields
		expected QueryRow
	}
	tests := []struct {
		insert insertQ
		get    getQ
	}{
		{
			insertQ{
				&fields{"name"},
				values{"some"},
			},
			getQ{
				&fields{"name"},
				QueryRow{
					"name": {"some"},
				},
			},
		},
		{
			insertQ{
				nil,
				values{256, "some"},
			},
			getQ{
				&fields{"id", "name"},
				QueryRow{
					"id":   {int64(256)},
					"name": {"some"},
				},
			},
		},
	}

	for _, test := range tests {
		id, err := insert("publisher", test.insert.f, &test.insert.v)
		if err != nil {
			t.Errorf("unexpected error occured: %v", err)
			continue
		}
		res, err := getOne("publisher", test.get.f, &where{
			[]string{fmt.Sprintf("id = %d", id)},
			"",
		}, nil)
		if err != nil {
			t.Errorf("unexpected error occured: %v", err)
			continue
		}
		if !compare_test_get_results(&res, &test.get.expected) {
			t.Errorf("get result not inserted - expected:\n%v\ngot %v", test.get.expected, res)
		}
	}
}
func TestInsertMany(t *testing.T) {
	if err := setupDB(); err != nil {
		t.Fatalf("fail to prepare database: %v", err)
	}

	tests := []struct {
		f        *fields
		v        []values
		expected int
	}{
		{
			&fields{"name"},
			[]values{{"some"}, {"other"}},
			2,
		},
	}

	for _, test := range tests {
		res, err := insertMany("publisher", test.f, &test.v)
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
