package repos

import (
	"database/sql"
	"expat-news/queue-manager/pkg/logger"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func GetDBInstance() *sql.DB {
	once.Do(func() {
		var err error
		const env = "QDB_FILE"
		file, exist := os.LookupEnv(env)
		if !exist || len(file) == 0 {
			file = ":memory:"
			logger.Message("DB runs in memory")
		}
		db, err = sql.Open("sqlite3", file)
		if err != nil {
			log.Fatal("Fail to open database: ", err)
		}
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	})
	// defer db.Close()
	return db
}
