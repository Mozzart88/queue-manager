package repos

import (
	"testing"
)

func TestGenWhereStatement(t *testing.T) {
	tests := []struct {
		where      where
		expected   string
		whantError bool
	}{
		{where{[]string{}, ""}, "empty where.fields", true},
		{where{[]string{"some = value", "other like '%%value'"}, ""}, "invalid union: ", true},
		{where{[]string{"some = value", "other like '%%value'"}, "some"}, "invalid union: some", true},
		{where{[]string{"some = value"}, ""}, "WHERE some = value", false},
		{where{[]string{"some = value", "other like '%value'"}, "AND"}, "WHERE some = value AND other like '%value'", false},
		{where{[]string{"some = value", "other like '%value'"}, "and"}, "WHERE some = value AND other like '%value'", false},
		{where{[]string{"some = value", "other like '%value'"}, "OR"}, "WHERE some = value OR other like '%value'", false},
		{where{[]string{"some = value", "other like '%value'", "one_more < 1"}, "AND"}, "WHERE some = value AND other like '%value' AND one_more < 1", false},
	}

	for _, test := range tests {
		var actual string
		err := genWhereStatement(&actual, test.where)
		if test.whantError {
			if err == nil {
				t.Errorf("genWhereStatement(res, %v) expected an error, got nil", test.where)
			} else if err.Error() != test.expected {
				t.Errorf("genWhereStatement(res, %v) expected an error: %s, got %s", test.where, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genWhereStatement(res, %v) = %s, want %s", test.where, actual, test.expected)
		}
	}
}

func TestGenOrderStatement(t *testing.T) {
	tests := []struct {
		order      order
		expected   string
		whantError bool
	}{
		{order{[]string{}, ""}, "empty order.fields", true},
		{order{[]string{"some", "other"}, "some"}, "invalid order: some", true},
		{order{[]string{"some"}, ""}, "ORDER BY some", false},
		{order{[]string{"some", "other"}, ""}, "ORDER BY some, other", false},
		{order{[]string{"some", "other"}, "DESC"}, "ORDER BY some, other DESC", false},
		{order{[]string{"some", "other"}, "desc"}, "ORDER BY some, other DESC", false},
		{order{[]string{"some", "other"}, "ASC"}, "ORDER BY some, other ASC", false},
	}

	for _, test := range tests {
		var actual string
		err := genOrderStatement(&actual, test.order)
		if test.whantError {
			if err == nil {
				t.Errorf("genOrderStatement(res, %v) expected an error, got nil", test.order)
			} else if err.Error() != test.expected {
				t.Errorf("genOrderStatement(res, %v) expected an error: %s, got %s", test.order, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genOrderStatement(res, %v) = %s, want %s", test.order, actual, test.expected)
		}
	}
}

func TestGenSelectStatement(t *testing.T) {
	var lim limit = 1
	tests := []struct {
		table      string
		fields     *fields
		where      *where
		order      *order
		limit      *limit
		expected   string
		whantError bool
	}{
		{"some", nil, nil, nil, nil, "SELECT * FROM some", false},
		{"some", &fields{"id", "name"}, nil, nil, nil, "SELECT id, name FROM some", false},
		{"some", &fields{"id", "name as pub_name"}, nil, nil, nil, "SELECT id, name as pub_name FROM some", false},
		{"some", &fields{}, nil, nil, nil, "SELECT * FROM some", false},
		{"some", nil, &where{[]string{"id = 1"}, ""}, nil, nil, "SELECT * FROM some WHERE id = 1", false},
		{"some", nil, &where{[]string{"id = 1"}, ""}, &order{[]string{"name", "timestamp"}, ""}, nil, "SELECT * FROM some WHERE id = 1 ORDER BY name, timestamp", false},
		{"some", nil, &where{[]string{"id = 1"}, ""}, &order{[]string{"name", "timestamp"}, ""}, &lim, "SELECT * FROM some WHERE id = 1 ORDER BY name, timestamp LIMIT 1", false},
		{"some", &fields{"id", "name"}, &where{[]string{"id = 1"}, ""}, &order{[]string{"name", "timestamp"}, ""}, &lim, "SELECT id, name FROM some WHERE id = 1 ORDER BY name, timestamp LIMIT 1", false},
		{"", &fields{"id", "name"}, &where{[]string{"id = 1"}, ""}, &order{[]string{"name", "timestamp"}, ""}, &lim, "empty table name", true},
	}

	for _, test := range tests {
		var actual string
		err := genSelectStatement(&actual, test.table, test.fields, test.where, test.order, test.limit)
		if test.whantError {
			if err == nil {
				t.Errorf("genSelectStatement(res, %v, %v, %v, %v, %v) expected an error, got nil", test.table, test.fields, test.where, test.order, test.limit)
			} else if err.Error() != test.expected {
				t.Errorf("genSelectStatement(res, %v, %v, %v, %v, %v) expected an error: %s, got %s", test.table, test.fields, test.where, test.order, test.limit, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genSelectStatement(res, %v, %v, %v, %v, %v) = %s, want %s", test.table, test.fields, test.where, test.order, test.limit, actual, test.expected)
		}
	}
}

func TestGenInsertStatement(t *testing.T) {
	tests := []struct {
		table      string
		fields     *fields
		values     *values
		expected   string
		whantError bool
	}{
		{"some", &fields{"name"}, &values{"bob"}, "INSERT INTO some (name) VALUES ('bob')", false},
		{"some", nil, &values{"bob", 1}, "INSERT INTO some  VALUES ('bob', 1)", false}, // !! Double space
	}

	for _, test := range tests {
		var actual string
		err := genInsertStatement(&actual, test.table, test.fields, test.values)
		if test.whantError {
			if err == nil {
				t.Errorf("genInsertStatement(res, %v, %v, %v) expected an error, got nil", test.table, test.fields, test.values)
			} else if err.Error() != test.expected {
				t.Errorf("genInsertStatement(res, %v, %v, %v) expected an error: %s, got %s", test.table, test.fields, test.values, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genInsertStatement(res, %v, %v, %v) = %s, want %s", test.table, test.fields, test.values, actual, test.expected)
		}
	}
}

func TestGenInsertManyStatement(t *testing.T) {
	tests := []struct {
		table      string
		fields     *fields
		values     *[]values
		expected   string
		whantError bool
	}{
		{"", nil, nil, "empty table name", true},
		{"some", nil, nil, "empty values", true},
		{"some", &fields{"name"}, nil, "empty values", true},
		{"some", &fields{"name"}, &[]values{}, "empty values", true},
		{"some", &fields{"name"}, &[]values{{"bob"}}, "INSERT INTO some (name) VALUES ('bob')", false},
		{"some", &fields{"name", "age"}, &[]values{{"bob", 1}}, "INSERT INTO some (name, age) VALUES ('bob', 1)", false},
		{"some", nil, &[]values{{"bob", 1}}, "INSERT INTO some  VALUES ('bob', 1)", false},       // !! Double space
		{"some", &fields{}, &[]values{{"bob", 1}}, "INSERT INTO some  VALUES ('bob', 1)", false}, // !! Double space
		{"some", &fields{"name", "age"}, &[]values{{"bob", 1}, {"ana", 2}}, "INSERT INTO some (name, age) VALUES ('bob', 1), ('ana', 2)", false},
		{"some", &fields{"name"}, &[]values{{"bob"}, {"ana"}}, "INSERT INTO some (name) VALUES ('bob'), ('ana')", false},
	}

	for _, test := range tests {
		var actual string
		err := genInsertManyStatement(&actual, test.table, test.fields, test.values)
		if test.whantError {
			if err == nil {
				t.Errorf("genInsertManyStatement(res, %v, %v, %v) expected an error, got nil", test.table, test.fields, test.values)
			} else if err.Error() != test.expected {
				t.Errorf("genInsertManyStatement(res, %v, %v, %v) expected an error: %s, got %s", test.table, test.fields, test.values, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genInsertManyStatement(res, %v, %v, %v) = %s, want %s", test.table, test.fields, test.values, actual, test.expected)
		}
	}
}

func TestGenDeleteStatement(t *testing.T) {
	tests := []struct {
		table      string
		where      *where
		expected   string
		whantError bool
	}{
		{"", nil, "empty table name", true},
		{"some", nil, "DELETE FROM some", false},
		{"some", &where{[]string{"name = 'bob'"}, ""}, "DELETE FROM some WHERE name = 'bob'", false},
	}

	for _, test := range tests {
		var actual string
		err := genDeleteStatement(&actual, test.table, test.where)
		if test.whantError {
			if err == nil {
				t.Errorf("genDeleteStatement(res, %v, %v) expected an error, got nil", test.table, test.where)
			} else if err.Error() != test.expected {
				t.Errorf("genDeleteStatement(res, %v, %v) expected an error: %s, got %s", test.table, test.where, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genDeleteStatement(res, %v, %v) = %s, want %s", test.table, test.where, actual, test.expected)
		}
	}
}

func TestGenUpdateStatement(t *testing.T) {
	tests := []struct {
		table      string
		fields     *fields
		where      *where
		expected   string
		whantError bool
	}{
		{"", nil, nil, "empty table name", true},
		{"some", nil, nil, "empty fields", true},
		{"some", &fields{"name = 'bob'"}, nil, "UPDATE some SET name = 'bob'", false},
		{"some", &fields{"name = 'bob'", "age = 30"}, nil, "UPDATE some SET name = 'bob', age = 30", false},
		{"some", &fields{"name = 'bob'", "age = 30"}, &where{[]string{"id = 10"}, ""}, "UPDATE some SET name = 'bob', age = 30 WHERE id = 10", false},
	}

	for _, test := range tests {
		var actual string
		err := genUpdateStatement(&actual, test.table, test.fields, test.where)
		if test.whantError {
			if err == nil {
				t.Errorf("genUpdateStatement(res, %v, %v, %v) expected an error, got nil", test.table, test.fields, test.where)
			} else if err.Error() != test.expected {
				t.Errorf("genUpdateStatement(res, %v, %v, %v) expected an error: %s, got %s", test.table, test.fields, test.where, test.expected, err.Error())
			}
		} else if test.expected != actual {
			t.Errorf("genUpdateStatement(res, %v, %v, %v) = %s, want %s", test.table, test.fields, test.where, actual, test.expected)
		}
	}
}
