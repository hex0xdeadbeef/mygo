package issuetool

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	validationURL       = "https://api.github.com/user"
	githubTokenPrefix   = "ghp_"
	invalidSpaceSymbols = "\t\n\v\f\r \u0085\u00A0"
)

var (
	testToken = "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"
)

type User struct {
	personalToken string
}

func isTokenStyleValid(userInput string) (string, error) {
	trimmedInput := strings.Trim(userInput, invalidSpaceSymbols)
	if strings.HasPrefix(trimmedInput, githubTokenPrefix) && len(trimmedInput) < 50 {
		return trimmedInput, nil
	}

	return "", errors.New("Invalid token style")
}

func IsTokenValidForQuery() bool {
	token, err := isTokenStyleValid(testToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured: %s\n", err)
		return false
	}

	client := http.Client{}
	tokenValidationRequest, err := http.NewRequest("GET", validationURL, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "While creating request error occured: %s\n", err)
		return false
	}

	tokenValidationRequest.Header.Set("Authorization", "Bearer "+token)

	response, err := client.Do(tokenValidationRequest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "While making request error occured: %s\n", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "%s\n", response.Status)
		return false
	}

	fmt.Fprintf(os.Stdout, "AUTHORIZATION DONE!")
	return true
}
