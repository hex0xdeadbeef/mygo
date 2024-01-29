package issuetool

import (
	"errors"
	cfg "golang/pkg/projects/a_issuetool/config"
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

func (tockenChecker *TokenChecker) IsTokenValidForQuery() error {
	err := tockenChecker.isTokenStyleValid()
	if err != nil {
		return errors.New("invalid token style")
	}

	client := http.Client{}
	tokenValidationRequest, err := http.NewRequest("GET", cfg.TokenValidationURL, nil)
	if err != nil {
		return err
	}

	tokenValidationRequest.Header.Set("Authorization", "Bearer "+tockenChecker.token)

	response, err := client.Do(tokenValidationRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	Logger.Println("User was successfully authorized")
	return nil
}

func (tokenChecker *TokenChecker) isTokenStyleValid() error {
	if len(tokenChecker.token) != 0 &&
		(strings.HasPrefix(tokenChecker.token, cfg.ClassicTokenPrefix) ||
			strings.HasPrefix(tokenChecker.token, cfg.FineGrainedPrefix)) &&
		len(tokenChecker.token) < 100 {
		return nil
	}

	return errors.New("no valid token")
}
