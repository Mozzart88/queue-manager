package repos

import (
	"database/sql"
	"log"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func getDBInstance() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", "data/data.db")
		if err != nil {
			log.Fatal("Fail to open database: ", err)
		}
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	})
	return db
}
