package crud

import (
	"testing"
)

func TestEquals(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` = ?"},
	}
	for i, test := range tests {
		if actual := Equals(test.field); actual != test.expected {
			t.Errorf("test %d: equals(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestNotEquals(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` != ?"},
	}
	for i, test := range tests {
		if actual := NotEquals(test.field); actual != test.expected {
			t.Errorf("test %d: notEquals(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestLike(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` LIKE ?"},
	}
	for i, test := range tests {
		if actual := Like(test.field); actual != test.expected {
			t.Errorf("test %d: like(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestNotLike(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` NOT LIKE ?"},
	}
	for i, test := range tests {
		if actual := NotLike(test.field); actual != test.expected {
			t.Errorf("test %d: notLike(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestLess(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` < ?"},
	}
	for i, test := range tests {
		if actual := Less(test.field); actual != test.expected {
			t.Errorf("test %d: less(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestLessEq(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` <= ?"},
	}
	for i, test := range tests {
		if actual := LessEq(test.field); actual != test.expected {
			t.Errorf("test %d: lessEq(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestGreater(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` > ?"},
	}
	for i, test := range tests {
		if actual := Greater(test.field); actual != test.expected {
			t.Errorf("test %d: greater(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}

func TestGreaterEq(t *testing.T) {
	tests := []struct {
		field    string
		expected string
	}{
		{"id", "`id` >= ?"},
	}
	for i, test := range tests {
		if actual := GreaterEq(test.field); actual != test.expected {
			t.Errorf("test %d: greaterEq(\"%s\") = %s, expected %s", i, test.field, actual, test.expected)
		}
	}
}
