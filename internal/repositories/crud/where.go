package crud

import (
	sg "expat-news/queue-manager/internal/repositories/statement_generators"
)

type Statement struct {
	Value      any
	Comparator func(feild string) string
}

type union int

const (
	U_Empty union = iota
	U_And
	U_Or
)

type Where struct {
	keys       []string
	statements map[string]Statement
	Union      union
}

func (s union) String() string {
	switch s {
	case U_Empty:
		return ""
	case U_And:
		return "AND"
	case U_Or:
		return "OR"
	default:
		return ""
	}
}

func (w *Where) prepare() *sg.Where {
	if w == nil {
		return nil
	}
	res := sg.Where{Fields: []string{}, Union: w.Union.String()}
	for _, key := range w.keys {
		stmt := w.statements[key]
		res.Fields = append(res.Fields, stmt.Comparator(key))
	}
	return &res
}

func NewWhere() *Where {
	w := &Where{}
	w.keys = []string{}
	w.statements = make(map[string]Statement)
	return w
}

func (w *Where) Add(colName string, value any, comporator func(field string) string) {
	if _, ok := w.statements[colName]; !ok {
		w.keys = append(w.keys, colName)
	}
	w.statements[colName] = Statement{value, comporator}
}

func (w *Where) Equals(colName string, value any) {
	w.Add(colName, value, Equals)
}

func (w Where) Len() int {
	return len(w.statements)
}

func (w *Where) values() []any {
	if w == nil {
		return nil
	}
	res := []any{}
	for _, key := range w.keys {
		stmt := w.statements[key]
		res = append(res, stmt.Value)
	}
	return res
}
