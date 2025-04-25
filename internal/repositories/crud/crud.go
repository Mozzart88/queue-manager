package crud

import (
	"database/sql"
	sg "expat-news/queue-manager/internal/repositories/statement_generators"
	"expat-news/queue-manager/pkg/utils"
	"fmt"
)

type Values []any
type Fields []string

func execSql(sql string, args ...any) (sql.Result, error) {
	db := GetDBInstance()
	res, err := db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func beginTx() error {
	_, err := execSql("BEGIN TRANSACTION")
	return err
}

func commitTx() error {
	_, err := execSql("COMMIT")
	return err
}

func rollbackTx() error {
	_, err := execSql("ROLLBACK")
	return err
}

func tx(callback func() (int, error)) (int, error) {
	if err := beginTx(); err != nil {
		return 0, err
	}
	val, err := callback()
	if err != nil {
		if err := rollbackTx(); err != nil {
			return 0, err
		}
		return 0, err
	}
	if err := commitTx(); err != nil {
		return 0, err
	}
	return val, nil
}

func placeholders(v *Values) sg.Values {
	n := len(*v)
	if n <= 0 {
		return sg.Values{}
	}
	parts := make(sg.Values, n)
	for i := range parts {
		parts[i] = "?"
	}
	return parts
}

func (f *Fields) prepare() *sg.Fields {
	if f == nil {
		return nil
	}
	return utils.Ptr(sg.Fields(*f))
}

func Insert(table string, f *Fields, v *Values) (int, error) {
	var sql string
	if v == nil || len(*v) == 0 {
		return -1, fmt.Errorf("values cant be empty in insert routine")
	}
	if err := sg.InsertStatement(&sql, table, f.prepare(), (placeholders(v))); err != nil {
		return -1, err
	}
	res, err := execSql(sql, *v...)
	if err != nil {
		return -1, err
	}
	resId, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}
	return int(resId), nil
}

func InsertMany(table string, f *Fields, v *[]Values) (int, error) {
	var sql string
	var result int64 = 0
	if v == nil {
		return -1, fmt.Errorf("values cant be empty in insert routine")
	}
	if err := sg.InsertStatement(&sql, table, f.prepare(), placeholders(&((*v)[0]))); err != nil {
		return -1, err
	}
	db := GetDBInstance()
	stmt, err := db.Prepare(sql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()
	for _, vals := range *v {
		res, err := stmt.Exec(vals...)
		if err != nil {
			return -1, err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return -1, nil
		}
		result += rowsAffected
	}
	return int(result), nil
}

func (o *Order) prepare() *sg.Order {
	if o == nil {
		return nil
	}
	res := sg.Order{Fields: []string{}, Order: o.Order}
	res.Fields = append(res.Fields, o.Fields...)
	return &res
}

func Update(table string, f *Fields, w *Where) (int, error) {
	var sql string
	if f == nil {
		return -1, fmt.Errorf("fields cant be empty in update routine")
	}
	if err := sg.UpdateStatement(&sql, table, f.prepare(), w.prepare()); err != nil {
		return -1, err
	}
	res, err := execSql(sql, w.values()...)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, nil
	}
	return int(affected), nil
}

func Delete(table string, w *Where) (int64, error) {
	var sql string
	if err := sg.DeleteStatement(&sql, table, w.prepare()); err != nil {
		return -1, err
	}
	res, err := execSql(sql, w.values()...)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, nil
	}
	return affected, nil
}

type cell struct {
	value any
}

func (c cell) Get() any {
	return c.value
}

type QueryRow map[string]cell

func rowsToMap(rows *sql.Rows) ([]QueryRow, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var results []QueryRow

	for rows.Next() {
		values := make([]any, len(columns))
		valuesPtrs := make([]any, len(columns))

		for i := range values {
			valuesPtrs[i] = &values[i]
		}

		if err := rows.Scan(valuesPtrs...); err != nil {
			return nil, err
		}

		rowMap := make(QueryRow)
		for i, colName := range columns {
			rowMap[colName] = cell{values[i]}
		}
		results = append(results, rowMap)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

type Limit int
type Order struct {
	Fields []string
	Order  string
}

func Get(table string, f *Fields, w *Where, o *Order, l *Limit) ([]QueryRow, error) {
	db := GetDBInstance()
	var sql string
	var limit *sg.Limit = nil
	if l != nil {
		limit = utils.Ptr(sg.Limit(*l))
	}
	if err := sg.SelectStatement(&sql, table, f.prepare(), w.prepare(), o.prepare(), limit); err != nil {
		return nil, err
	}
	args := w.values()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMap(rows)
}

func GetOne(table string, f *Fields, w *Where, o *Order) (QueryRow, error) {
	var l Limit = 1
	res, err := Get(table, f, w, o, &l)
	if len(res) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return res[0], nil
}
