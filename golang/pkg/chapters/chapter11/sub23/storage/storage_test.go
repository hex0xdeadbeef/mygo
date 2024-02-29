package storage

import (
	"strings"
	"testing"
)

func TestCheckQuouta(t *testing.T) {
	const (
		user = "joe@example.org"

		wantSubstring = "98% of your quota"
	)
	var (
		notifiedUser string
		notifiedMsg  string

		// Save the original function value of "notifyUser()" so that it'll be restored after test function has returned.
		notifyUserOriginal = notifyUser
	)

	defer func() {
		notifyUser = notifyUserOriginal
	}()

	// Reinitialization of unexported package function variable
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}

	// Simulate a 980Mb-used condition
	usage[user] = 98e8

	// ops
	CheckQuota(user)

	if notifiedUser == "" && notifiedMsg == "" {
		t.Fatalf("notifyUser not called")
	}

	if notifiedUser != user {
		t.Errorf("wrong user (%s) notified, want %s", notifiedMsg, user)
	}

	if !strings.Contains(notifiedMsg, wantSubstring) {
		t.Errorf("unexpected notification message <<%s>>, want substring %q", notifiedMsg, wantSubstring)
	}

}
