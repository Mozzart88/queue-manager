package statementgenerators

import (
	"testing"
)

func TestGenWhereStatement(t *testing.T) {
	tests := []struct {
		where      Where
		expected   string
		whantError bool
	}{
		{Where{[]string{}, ""}, "empty where.fields", true},
		{Where{[]string{"some = value", "other like '%%value'"}, ""}, "invalid union: ", true},
		{Where{[]string{"some = value", "other like '%%value'"}, "some"}, "invalid union: some", true},
		{Where{[]string{"some = value"}, ""}, "WHERE some = value", false},
		{Where{[]string{"some = value", "other like '%value'"}, "AND"}, "WHERE some = value AND other like '%value'", false},
		{Where{[]string{"some = value", "other like '%value'"}, "and"}, "WHERE some = value AND other like '%value'", false},
		{Where{[]string{"some = value", "other like '%value'"}, "OR"}, "WHERE some = value OR other like '%value'", false},
		{Where{[]string{"some = value", "other like '%value'", "one_more < 1"}, "AND"}, "WHERE some = value AND other like '%value' AND one_more < 1", false},
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
		order      Order
		expected   string
		whantError bool
	}{
		{Order{[]string{}, ""}, "empty order.fields", true},
		{Order{[]string{"some", "other"}, "some"}, "invalid order: some", true},
		{Order{[]string{"some"}, ""}, "ORDER BY some", false},
		{Order{[]string{"some", "other"}, ""}, "ORDER BY some, other", false},
		{Order{[]string{"some", "other"}, "DESC"}, "ORDER BY some, other DESC", false},
		{Order{[]string{"some", "other"}, "desc"}, "ORDER BY some, other DESC", false},
		{Order{[]string{"some", "other"}, "ASC"}, "ORDER BY some, other ASC", false},
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
	var lim Limit = 1
	tests := []struct {
		table      string
		fields     *Fields
		where      *Where
		order      *Order
		limit      *Limit
		expected   string
		whantError bool
	}{
		{"some", nil, nil, nil, nil, "SELECT * FROM some", false},
		{"some", &Fields{"id", "name"}, nil, nil, nil, "SELECT id, name FROM some", false},
		{"some", &Fields{"id", "name as pub_name"}, nil, nil, nil, "SELECT id, name as pub_name FROM some", false},
		{"some", &Fields{}, nil, nil, nil, "SELECT * FROM some", false},
		{"some", nil, &Where{[]string{"id = ?"}, ""}, nil, nil, "SELECT * FROM some WHERE id = ?", false},
		{"some", nil, &Where{[]string{"id = ?"}, ""}, &Order{[]string{"name", "timestamp"}, ""}, nil, "SELECT * FROM some WHERE id = ? ORDER BY name, timestamp", false},
		{"some", nil, &Where{[]string{"id = ?"}, ""}, &Order{[]string{"name", "timestamp"}, ""}, &lim, "SELECT * FROM some WHERE id = ? ORDER BY name, timestamp LIMIT 1", false},
		{"some", &Fields{"id", "name"}, &Where{[]string{"id = ?"}, ""}, &Order{[]string{"name", "timestamp"}, ""}, &lim, "SELECT id, name FROM some WHERE id = ? ORDER BY name, timestamp LIMIT 1", false},
		{"", &Fields{"id", "name"}, &Where{[]string{"id = ?"}, ""}, &Order{[]string{"name", "timestamp"}, ""}, &lim, "empty table name", true},
	}

	for _, test := range tests {
		var actual string
		err := SelectStatement(&actual, test.table, test.fields, test.where, test.order, test.limit)
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
		fields     *Fields
		values     *Values
		expected   string
		whantError bool
	}{
		{"some", &Fields{"name"}, &Values{"bob"}, "INSERT INTO some (name) VALUES ('bob')", false},
		{"some", nil, &Values{"bob", 1}, "INSERT INTO some VALUES ('bob', 1)", false},
	}

	for _, test := range tests {
		var actual string
		err := InsertStatement(&actual, test.table, test.fields, *test.values)
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
		fields     *Fields
		values     *[]Values
		expected   string
		whantError bool
	}{
		{"", nil, nil, "empty table name", true},
		{"some", nil, nil, "empty values", true},
		{"some", &Fields{"name"}, nil, "empty values", true},
		{"some", &Fields{"name"}, &[]Values{}, "empty values", true},
		{"some", &Fields{"name"}, &[]Values{{"bob"}}, "INSERT INTO some (name) VALUES ('bob')", false},
		{"some", &Fields{"name", "age"}, &[]Values{{"bob", 1}}, "INSERT INTO some (name, age) VALUES ('bob', 1)", false},
		{"some", nil, &[]Values{{"bob", 1}}, "INSERT INTO some VALUES ('bob', 1)", false},
		{"some", &Fields{}, &[]Values{{"bob", 1}}, "INSERT INTO some VALUES ('bob', 1)", false},
		{"some", &Fields{"name", "age"}, &[]Values{{"bob", 1}, {"ana", 2}}, "INSERT INTO some (name, age) VALUES ('bob', 1), ('ana', 2)", false},
		{"some", &Fields{"name"}, &[]Values{{"bob"}, {"ana"}}, "INSERT INTO some (name) VALUES ('bob'), ('ana')", false},
	}

	for _, test := range tests {
		var actual string
		err := InsertManyStatement(&actual, test.table, test.fields, test.values)
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
		where      *Where
		expected   string
		whantError bool
	}{
		{"", nil, "empty table name", true},
		{"some", nil, "DELETE FROM some", false},
		{"some", &Where{[]string{"name = 'bob'"}, ""}, "DELETE FROM some WHERE name = 'bob'", false},
	}

	for _, test := range tests {
		var actual string
		err := DeleteStatement(&actual, test.table, test.where)
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
		fields     *Fields
		where      *Where
		expected   string
		whantError bool
	}{
		{"", nil, nil, "empty table name", true},
		{"some", nil, nil, "empty fields", true},
		{"some", &Fields{"name = 'bob'"}, nil, "UPDATE some SET name = 'bob'", false},
		{"some", &Fields{"name = 'bob'", "age = 30"}, nil, "UPDATE some SET name = 'bob', age = 30", false},
		{"some", &Fields{"name = 'bob'", "age = 30"}, &Where{[]string{"id = 10"}, ""}, "UPDATE some SET name = 'bob', age = 30 WHERE id = 10", false},
	}

	for _, test := range tests {
		var actual string
		err := UpdateStatement(&actual, test.table, test.fields, test.where)
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
