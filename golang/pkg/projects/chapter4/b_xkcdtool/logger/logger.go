package logger

import (
	"errors"
	cfg "golang/pkg/projects/chapter4/b_xkcdtool/config"
	"log"
	"os"
)

func Init() *log.Logger {
	file, err := os.OpenFile(cfg.LogFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(errors.New("logger creation fail"))
	}

	return log.New(file, "-", log.Ldate|log.Ltime)
}
