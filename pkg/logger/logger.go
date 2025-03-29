package logger

import (
	"log"
	"os"
	"sync"
)

var (
	errorLogger *log.Logger
	once        sync.Once
)

func getErrorLogger() *log.Logger {
	once.Do(func() {
		errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	})
	return errorLogger
}

func Message(msg string) {
	log.Default().Println(msg)
}

func Error(err string) {
	getErrorLogger().Println(err)
}
