package repos

import (
	"database/sql"
)

func execSql(sql string) (sql.Result, error) {
	db := GetDBInstance()
	res, err := db.Exec(sql)
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

func insert(table string, f *fields, v *values) (int, error) {
	var sql string
	if err := genInsertStatement(&sql, table, f, v); err != nil {
		return -1, err
	}
	res, err := execSql(sql)
	if err != nil {
		return -1, err
	}
	resId, err := res.LastInsertId()
	if err != nil {
		return -1, nil
	}
	return int(resId), nil
}

func insertMany(table string, f *fields, v *[]values) (int, error) {
	var sql string
	if err := genInsertManyStatement(&sql, table, f, v); err != nil {
		return -1, err
	}
	res, err := execSql(sql)
	if err != nil {
		return -1, err
	}
	result, err := res.RowsAffected()
	if err != nil {
		return -1, nil
	}
	return int(result), nil
}

func update(table string, f *fields, w *where) (int, error) {
	var sql string
	if err := genUpdateStatement(&sql, table, f, w); err != nil {
		return -1, err
	}
	res, err := execSql(sql)
	if err != nil {
		return -1, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return -1, nil
	}
	return int(affected), nil
}

func delete(table string, w *where) (int64, error) {
	var sql string
	if err := genDeleteStatement(&sql, table, w); err != nil {
		return -1, err
	}
	res, err := execSql(sql)
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

func (c *cell) Get() any {
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

func get(table string, f *fields, w *where, o *order, l *limit) ([]QueryRow, error) {
	db := GetDBInstance()
	var sql string
	if err := genSelectStatement(&sql, table, f, w, o, l); err != nil {
		return nil, err
	}
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMap(rows)
}

func getOne(table string, f *fields, w *where, o *order) (QueryRow, error) {
	var l limit = 1
	res, err := get(table, f, w, o, &l)
	if len(res) == 0 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return res[0], nil
}
