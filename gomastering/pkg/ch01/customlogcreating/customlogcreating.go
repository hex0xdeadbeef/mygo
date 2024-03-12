package customlogcreating

import (
	"fmt"
	"log"
	"os"
	"unicode"
)

func newTempLog(fileName string) (*log.Logger, func(), error) {

	const (
		logSuffix = ".log"
	)

	if err := isFileNameValid(fileName); err != nil {
		return nil, nil, fmt.Errorf("validating file name: %s", err)
	}

	LOGFILE := fileName + logSuffix
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f.Close()
		return nil, nil, err
	}

	// we can use the "*log.Logger.SetFlags(...)" for the same purpose
	logger := log.New(f, "-> ", log.LstdFlags|log.Ldate|log.Lshortfile)

	return logger, func() { f.Close() }, nil
}

func isFileNameValid(fileName string) error {
	const (
		zeroLength = 0
		maxLength  = 32
	)

	if len(fileName) == zeroLength || len(fileName) > maxLength {
		return fmt.Errorf("invalid file name length: %d", len(fileName))
	}

	for _, r := range fileName {
		if !unicode.IsLetter(r) {
			return fmt.Errorf("invalid symbol in file name: %c", r)
		}
	}

	return nil
}

func Using() {
	logger, closer, err := newTempLog("new")
	defer closer()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger.Println("Hey there!")

}
