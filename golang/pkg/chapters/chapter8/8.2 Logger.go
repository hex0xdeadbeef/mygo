package chapter8

import (
	"log"
	"os"
)

const (
	serverLogFile = "serverlog.txt"
	clientLogFile = "clientlog.txt"
)

func NewLogger(fileName string) (*log.Logger, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return log.New(file, "-> ", log.Ldate|log.Ltime), nil
}
