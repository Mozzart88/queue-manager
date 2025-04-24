package statementgenerators

import (
	"errors"
	"fmt"
	"strings"
)

type Order struct {
	Fields []string
	Order  string
}

type Where struct {
	Fields []string
	Union  string
}

type Fields []string
type Values []any

type Limit int

func genWhereStatement(dst *string, w Where) error {
	var fields string
	if len(w.Fields) == 0 {
		return errors.New("empty where.fields")
	}
	if len(w.Fields) == 1 {
		fields = w.Fields[0]
	} else {
		if !strings.EqualFold(w.Union, "AND") && !strings.EqualFold(w.Union, "OR") {
			return fmt.Errorf("invalid union: %s", w.Union)
		}
		union := " " + strings.ToUpper(w.Union) + " "
		fields = strings.Join(w.Fields, union)
	}
	*dst = fmt.Sprintf("WHERE %s", fields)
	return nil
}

func genOrderStatement(dst *string, o Order) error {
	var fields string
	if len(o.Fields) == 0 {
		return errors.New("empty order.fields")
	}
	if len(o.Fields) == 1 {
		fields = o.Fields[0]
	} else {
		fields = strings.Join(o.Fields, ", ")
	}
	if len(o.Order) > 0 && !strings.EqualFold(o.Order, "ASC") && !strings.EqualFold(o.Order, "DESC") {
		return fmt.Errorf("invalid order: %s", o.Order)
	}
	*dst = strings.Trim(fmt.Sprintf("ORDER BY %s %s", fields, strings.ToUpper(o.Order)), " ")
	return nil
}

func SelectStatement(dst *string, table string, fields *Fields, where *Where, order *Order, limit *Limit) error {
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

func insertVals(v *Values) string {
	var vals []string
	for _, value := range *v {
		if value == "?" {
			vals = append(vals, value.(string))
		} else if val, ok := value.(int); ok {
			vals = append(vals, fmt.Sprintf("%d", val))
		} else {
			vals = append(vals, fmt.Sprintf("'%v'", value))
		}
	}
	return fmt.Sprintf("(%s)", strings.Join(vals, ", "))
}

func InsertStatement(dst *string, table string, f *Fields, v Values) error {
	return InsertManyStatement(dst, table, f, &[]Values{v})
}

func InsertManyStatement(dst *string, table string, f *Fields, v *[]Values) error {
	var valuesArr []string
	if len(table) == 0 {
		return fmt.Errorf("empty table name")
	}
	if v == nil || len(*v) == 0 {
		return fmt.Errorf("empty values")
	}
	for _, values := range *v {
		valuesArr = append(valuesArr, insertVals(&values))
	}
	fields := ""
	if f != nil && len(*f) > 0 {
		fields = fmt.Sprintf(" (%s)", strings.Join(*f, ", "))
	}
	sql := fmt.Sprintf("INSERT INTO %s%s VALUES %s", table, fields, strings.Join(valuesArr, ", "))
	*dst = sql
	return nil
}

func UpdateStatement(dst *string, table string, f *Fields, w *Where) error {
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

func DeleteStatement(dst *string, table string, w *Where) error {
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
	*dst = sql
	return nil
}
