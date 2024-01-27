package issuetool

import (
	"flag"
	"fmt"
	cfg "golang/pkg/projects/issuetool/config"
	lgr "golang/pkg/projects/issuetool/logger"
	"log"
	"os"
)

var (
	Logger *log.Logger

	MethodConversion map[string]string
	MethodsParams    map[string]map[string]string

	token, method, owner, repo *string

	MethodKey   string
	RequestBody map[string]string
)

func init() {
	// Logger initialization
	{
		Logger = lgr.LoggerInit()
	}

	// Required params flags initialization
	{
		token = flag.String("token", "", "GitHub personal access token")
		method = flag.String("method", "", "Method (create, read, update, delete)")
		owner = flag.String("owner", "", "Repo owner")
		repo = flag.String("repo", "", "Repo name")
	}

	// Map initialization
	{
		// Methods mapping
		MethodConversion = map[string]string{
			"create": "POST",
			"edit":   "PATCH",
			"read":   "GET",
			"delete": "PUT",
		}

		MethodsParams = make(map[string]map[string]string)
		// Create
		MethodsParams["POST"] = map[string]string{
			"owner": "",
			"repo":  "",
			"title": "",

			"body":      "",
			"assignee":  "",
			"milestone": "",
			"labels":    "",
			"assignees": "",
		}

		// Get
		MethodsParams["GET"] = map[string]string{
			"owner":        "",
			"repo":         "",
			"issue_number": "",
		}

		// Edit
		MethodsParams["PATCH"] = map[string]string{
			"owner":        "",
			"repo":         "",
			"issue_number": "",

			"title":        "",
			"body":         "",
			"assignee":     "",
			"state":        "",
			"state_reason": "",
			"milestone":    "",
			"labels":       "",
			"assignees":    "",
		}

		// Delete
		MethodsParams["PUT"] = map[string]string{
			"owner":        "",
			"repo":         "",
			"issue_number": "",

			"lock_reason": "",
		}
	}

}

func Execute() {
	Logger.Printf("Session has created")
	flag.Parse()

	// Required params initialization
	if len(os.Args[1:]) != cfg.ValidParamsCount {
		Logger.Fatalf("Invalid parameters count\n")
	}

	// Checks
	tokenValidation()
	isMethodValid()
	areRequiredParamsValid()

	// Fill owner/repo
	requestBodyPreinit()

	// UI call
}

func tokenValidation() {
	var tokenChecker TokenChecker
	tokenChecker.init(*token)

	if tokenChecker.IsTokenValidForQuery() {
		fmt.Fprintf(os.Stdout, "\tYOU'VE BEEN SUCCESSFULLY AUTHORIZED\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "\n\tAUTHORIZATION FAILED\n\n")
	}
}

func isMethodValid() {
	for validMethod := range MethodConversion {
		if validMethod == *method {
			return
		}
	}
	log.Fatalf("No appropriate method parameter\n")
}

func areRequiredParamsValid() {
	if *repo == "" || *owner == "" {
		Logger.Fatalf("Required parameter empty\n")
	}
}

func requestBodyPreinit() {
	MethodKey = MethodConversion[*method]
	RequestBody := MethodsParams[MethodKey]

	RequestBody["owner"] = *owner
	RequestBody["repo"] = *repo
}
