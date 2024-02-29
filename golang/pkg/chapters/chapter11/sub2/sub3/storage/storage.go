package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

// Email sender configuration.
// NOTE: never put passwords in source code
const (
	sender   = "notifications@example.com"
	password = "correcthorsebatterystaple"
	hostname = "smtp.example.com"

	template = `Warning: you're using %d bytes of storage,
	%d%% of your quota`
)

var (
	usage = make(map[string]int64)

	// default realization, that will be replaced in test
	notifyUser = func(username, msg string) {
		auth := smtp.PlainAuth("", sender, password, hostname)
		err := smtp.SendMail(hostname+"587", auth, sender, []string{username}, []byte(msg))
		if err != nil {
			log.Printf("smtp.SendMail(%s) failed: %s", username, err)
		}

	}
)

func bytesInUse(username string) int64 {
	return usage[username]
}

func CheckQuota(username string) {

	const (
		hundredPercents = 100
		quota           = 10e9

		quotaLimitInPercents = 90
	)

	var (
		percent int64

		msg string
	)

	used := bytesInUse(username)
	percent = hundredPercents * used / quota

	if percent < quotaLimitInPercents {
		// OK
		return
	}

	msg = fmt.Sprintf(template, used, percent)
	notifyUser(username, msg)
}
