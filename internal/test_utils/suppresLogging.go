package test_utils

import (
	"bytes"
	"log"
	"os"
)

func SuppressLogging() func() {
	oldStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stderr = w
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	return func() {
		os.Stderr = oldStderr
		log.SetOutput(nil)
	}
}
