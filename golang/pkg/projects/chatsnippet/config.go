package chatsnippet

import (
	"strings"
)

const (
	logFilePath = "/Users/dmitriymamykin/Desktop/goprojects/golang/pkg/projects/chatsnippet/log.txt"

	sessionStartedMsg = "SESSION STARTED"
	sessionCloseddMsg = "SESSION CLOSED"

	welcomeMsgTemplate = "USER %s ENTERED THE CHAT"
	goodbyeMsgTemplate = "USER: %s LEFT THE CHAT"
)

var (
	logSplitter = strings.Repeat("-", 128)
)
