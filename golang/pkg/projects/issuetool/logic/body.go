package issuetool

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	cfg "golang/pkg/projects/issuetool/config"
	lgr "golang/pkg/projects/issuetool/logger"
	ui "golang/pkg/projects/issuetool/ui"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	Logger *log.Logger

	methods                    map[string]Empty
	methodKey                  string
	token, method, owner, repo *string
	issueNumber                *int
)

func init() {
	// Logger initialization
	{
		Logger = lgr.LoggerInit()
	}

	// Required params flags initialization
	{
		token = flag.String("token", "", "GitHub personal access token")
		method = flag.String("method", "", "Method (create, read, update, close, reopen)")
		owner = flag.String("owner", "", "Repo owner")
		repo = flag.String("repo", "", "Repo name")
		issueNumber = flag.Int("issuenumber", -1, "for create/get/edit methods")
	}

	// Methods mapping
	methods = map[string]Empty{
		"create": {},
		"read":   {},
		"edit":   {},
		"close":  {},
		"reopen": {},
	}
}

/* Controller */
func Execute() {
	Logger.Println("Session created")
	defer Logger.Println("Session closed")

	flag.Parse()

	// Required params initialization
	if len(os.Args[1:]) > cfg.ValidParamsCount {
		Logger.Println(errors.New("invalid parameters count").Error())
	}
	Logger.Printf("USER DATA: \ntoken: %s\nmethod: %s\nowner: %s\nrepo: %s\nissue number: %v\n", *token, *method, *owner, *repo, *issueNumber)

	// Validations
	err := methodValidation()
	if err != nil {
		Logger.Println(err.Error())
		return
	}
	err = paramsValidation()
	if err != nil {
		Logger.Println(err.Error())
		return
	}
	err = tokenValidation()
	if err != nil {
		ui.NonAuthPrint()
		return
	}
	ui.AuthPrint()

	// Choose a method
	switch *method {
	case "create":
		err = create(&Issue{
			Title: "Дорогой друг...",
			Body:  `Это я - твой единственный зритель. Я на протяжении многих лет создавал иллюзию того, что тебя смотрят много людей, но это был я. Сейчас напишу это сообщение со всех аккаунтов.`,
		})
		if err != nil {
			Logger.Println(err)
		}

	case "read":
		err = read()
		if err != nil {
			Logger.Println(err)
		}
	case "edit":
		err = edit(&Issue{Labels: []string{"documentation", "good first "}})
		if err != nil {
			Logger.Println(err)
		}
	case "close":
		err = close(&Issue{State: "closed"})
		if err != nil {
			Logger.Println(err)
		}
	case "reopen":
		err = reopen(&Issue{State: "open"})
		if err != nil {
			Logger.Println(err)
		}
	}
}

/* Main methods */
func create(issue *Issue) error {
	// Marshalling params
	requestBody, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	// Filling request
	request, err := http.NewRequest("POST", getCreationLink(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	setHeaders(request)

	// Do a request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return errors.New(response.Status)
	}

	// Write response data
	var respUnmarshaledBody Issue
	if err := json.NewDecoder(response.Body).Decode(&respUnmarshaledBody); err != nil {
		return err
	}
	Logger.Printf("RESPONSE INFO:\nURL: %s\nCreation time: %v\nState: %s\nTitle: %s\nBody:\n%s\nUser login: %s\nUser URL: %s\n",
		respUnmarshaledBody.HTMLURL,
		respUnmarshaledBody.CreatedAt,
		respUnmarshaledBody.State,
		respUnmarshaledBody.Title,
		respUnmarshaledBody.Body,
		respUnmarshaledBody.User.Login,
		respUnmarshaledBody.User.HTMLURL,
	)

	return nil
}

func read() error {
	// Filling request
	request, err := http.NewRequest("GET", getGeneralPurposeLink(), nil)
	if err != nil {
		return err
	}
	setHeaders(request)

	// Do a request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	// Write response data
	var respUnmarshaledBody Issue
	if err := json.NewDecoder(response.Body).Decode(&respUnmarshaledBody); err != nil {
		return err
	}
	fmt.Printf("RESPONSE INFO:\nURL: %s\nCreation time: %v\nState: %s\nTitle: %s\nBody:\n%s\n",
		respUnmarshaledBody.HTMLURL,
		respUnmarshaledBody.CreatedAt,
		respUnmarshaledBody.State,
		respUnmarshaledBody.Title,
		respUnmarshaledBody.Body,
	)

	return nil
}

func edit(issue *Issue) error {
	requestBody, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PATCH", getGeneralPurposeLink(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	setHeaders(request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	} else {
		Logger.Println("Edition is done!")
	}

	return nil
}

func close(issue *Issue) error {
	requestBody, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PATCH", getGeneralPurposeLink(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	setHeaders(request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	} else {
		Logger.Println("Close is done!")
	}

	return nil
}

func reopen(issue *Issue) error {
	requestBody, err := json.Marshal(issue)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PATCH", getGeneralPurposeLink(), bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	setHeaders(request)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	} else {
		Logger.Println("Reopening is done!")
	}

	return nil
}

/* Validation methods */
func tokenValidation() error {
	var tokenChecker TokenChecker
	tokenChecker.init(*token)

	return tokenChecker.IsTokenValidForQuery()
}

func methodValidation() error {
	for existingMethod := range methods {
		if existingMethod == *method {
			return nil
		}
	}

	return errors.New("no appropriate method parameter")
}

func paramsValidation() error {
	if *repo == "" || *owner == "" {

		return errors.New("nequired parameter empty")
	}

	fullParams, isCreateMehodChosen := len(os.Args[1:]) == cfg.ValidParamsCount, *method == "create"
	if !fullParams {
		return nil
	}

	if !isCreateMehodChosen {
		if *issueNumber < 0 {
			return errors.New("invalid issue number value")
		}
	}

	return nil

}

/* Subroutines */
func setHeaders(request *http.Request) {
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+*token)
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}

func getCreationLink() string {
	return fmt.Sprintf(cfg.CreationLink, *owner, *repo)
}

func getGeneralPurposeLink() string {
	return fmt.Sprintf(cfg.DeletionLink, *owner, *repo, strconv.Itoa(*issueNumber))
}
