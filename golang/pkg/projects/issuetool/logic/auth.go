package issuetool

import (
	"errors"
	cfg "golang/pkg/projects/issuetool/config"
	"net/http"
	"strings"
)

type TokenChecker struct {
	token string
}

func (tokenChecker *TokenChecker) init(inputToken string) {
	tokenChecker.token = tokenChecker.getFormattedToken(inputToken)
}

func (tokenChecker *TokenChecker) getFormattedToken(inputToken string) string {
	return strings.Trim(inputToken, cfg.SpaceCutSet)
}

func (tockenChecker *TokenChecker) IsTokenValidForQuery() bool {
	if tockenChecker.isTokenStyleValid() {
		Logger.Printf("While checking style error occured: %s\n", errors.New("invalid token style"))
		return false
	}

	client := http.Client{}
	tokenValidationRequest, err := http.NewRequest("GET", cfg.TokenValidationURL, nil)
	if err != nil {
		Logger.Printf("While creating request error occured: %s\n", err)
		return false
	}

	tokenValidationRequest.Header.Set("Authorization", "Bearer "+tockenChecker.token)

	response, err := client.Do(tokenValidationRequest)
	if err != nil {
		Logger.Printf("While making request error occured: %s\n", err)
		return false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		Logger.Printf("Response status code: %s\n", response.Status)
		return false
	}

	Logger.Printf("User was successfully authorized\n")
	return true
}

func (tokenChecker *TokenChecker) isTokenStyleValid() bool {
	if len(tokenChecker.token) != 0 &&
		(strings.HasPrefix(tokenChecker.token, cfg.ClassicTokenPrefix) &&
			strings.HasPrefix(tokenChecker.token, cfg.FineGrainedPrefix)) &&
		len(tokenChecker.token) < 100 {
		return true
	}

	return false
}
