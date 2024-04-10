package chatsnippet

import (
	"fmt"
	"log"
	"os"
)

type CustomLogger struct {
	logger  *log.Logger
	logFile *os.File
}

func (cl *CustomLogger) Close() error {
	if err := cl.logFile.Close(); err != nil {
		return fmt.Errorf("closing log file: %v", err)
	}

	return nil
}

func getLogger() *CustomLogger {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("creating log file: %v", err)
	}

	logger := log.New(file, "", log.LUTC|log.Ldate|log.Lshortfile|log.Ltime)

	return &CustomLogger{logger: logger, logFile: file}
}
