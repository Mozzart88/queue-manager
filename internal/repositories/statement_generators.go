package repos

import (
	"errors"
	"fmt"
	"strings"
)

type order struct {
	fields []string
	order  string
}

type where struct {
	fields []string
	union  string
}

type fields []string
type values []any

type limit int

func genWhereStatement(dst *string, w where) error {
	var fields string
	if len(w.fields) == 0 {
		return errors.New("empty where.fields")
	}
	if len(w.fields) == 1 {
		fields = w.fields[0]
	} else {
		if !strings.EqualFold(w.union, "AND") && !strings.EqualFold(w.union, "OR") {
			return fmt.Errorf("invalid union: %s", w.union)
		}
		union := " " + strings.ToUpper(w.union) + " "
		fields = strings.Join(w.fields, union)
	}
	*dst = fmt.Sprintf("WHERE %s", fields)
	return nil
}

func genOrderStatement(dst *string, o order) error {
	var fields string
	if len(o.fields) == 0 {
		return errors.New("empty order.fields")
	}
	if len(o.fields) == 1 {
		fields = o.fields[0]
	} else {
		fields = strings.Join(o.fields, ", ")
	}
	if len(o.order) > 0 && !strings.EqualFold(o.order, "ASC") && !strings.EqualFold(o.order, "DESC") {
		return fmt.Errorf("invalid order: %s", o.order)
	}
	*dst = strings.Trim(fmt.Sprintf("ORDER BY %s %s", fields, strings.ToUpper(o.order)), " ")
	return nil
}

func genSelectStatement(dst *string, table string, fields *fields, where *where, order *order, limit *limit) error {
	if len(table) == 0 {
		return fmt.Errorf("empty table name")
	}
	f := "*"
	if fields != nil && len(*fields) > 0 {
		f = strings.Join(*fields, ", ")
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", f, table)
	if where != nil {
		var w string
		if err := genWhereStatement(&w, *where); err != nil {
			return err
		}
		sql = fmt.Sprintf("%s %s", sql, w)
	}
	if order != nil {
		var o string
		if err := genOrderStatement(&o, *order); err != nil {
			return err
		}
		sql = fmt.Sprintf("%s %s", sql, o)
	}
	if limit != nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, *limit)
	}
	*dst = sql
	return nil
}

func genInsertStatement(dst *string, table string, f *fields, v *values) error {
	if len(table) == 0 {
		return fmt.Errorf("empty table name")
	}
	if v == nil || len(*v) == 0 {
		return fmt.Errorf("empty values")
	}
	var vals []string
	for _, value := range *v {
		if val, ok := value.(int); ok {
			vals = append(vals, fmt.Sprintf("%d", val))
		} else {
			vals = append(vals, fmt.Sprintf("'%v'", value))
		}
	}
	var values string = fmt.Sprintf("(%s)", strings.Join(vals, ", "))
	fields := ""
	if f != nil && len(*f) > 0 {
		fields = fmt.Sprintf("(%s)", strings.Join(*f, ", "))
	}
	sql := fmt.Sprintf("INSERT INTO %s %s VALUES %s", table, fields, values)
	*dst = sql
	return nil
}

func genUpdateStatement(dst *string, table string, f *fields, w *where) error {
	if len(table) == 0 {
		return fmt.Errorf("empty table name")
	}
	if f == nil || len(*f) == 0 {
		return fmt.Errorf("empty fields")
	}
	fields := strings.Join(*f, ", ")
	sql := fmt.Sprintf("UPDATE %s SET %s", table, fields)
	if w != nil {
		var where string
		if err := genWhereStatement(&where, *w); err != nil {
			return err
		}
		sql = fmt.Sprintf("%s %s", sql, where)
	}
	*dst = sql
	return nil
}

func genDeleteStatement(dst *string, table string, w *where, l *limit) error {
	if len(table) == 0 {
		return fmt.Errorf("empty table name")
	}
	sql := fmt.Sprintf("DELETE FROM %s", table)
	if w != nil {
		var where string
		if err := genWhereStatement(&where, *w); err != nil {
			return err
		}
		sql = fmt.Sprintf("%s %s", sql, where)
	}
	if l != nil {
		sql = fmt.Sprintf("%s LIMIT %d", sql, *l)
	}
	*dst = sql
	return nil
}
