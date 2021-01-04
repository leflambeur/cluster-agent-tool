package rancher

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/rancher/log"
	"net/url"
)

const (
	authTypeToken    = "token"
	authTypeLocal    = "local"
)

type userCredentials struct {
	Username string
	Password string
}

func validateURL(val interface{}) error {
	inputURL, ok := val.(string)
	if !ok {
		return fmt.Errorf("non-string input")
	}
	if _, err := url.ParseRequestURI(inputURL); err != nil {
		return fmt.Errorf("invalid url: %v", err)
	}
	return nil
}

func askForRancherServerDetails() (string, error) {
	url := ""
	q := &survey.Input{
		Message: "Enter Rancher URL: ",
		Help:    "Enter the URL of your Rancher server. Example: https://rancher.company.com",
	}
	err := survey.AskOne(q, &url, survey.WithValidator(survey.Required), survey.WithValidator(validateURL))
	if err != nil {
		return "", fmt.Errorf("error asking for rancher server details: %v", err)
	}
	return url, nil
}

func askForAuthenticationMethod(options []string) (string, error) {
	selectedAuth := ""
	q := &survey.Select{
		Message: "Select authentication method (Note: only Token and Local Supported)",
		Options: options,
		Help:    "TODO: Select one of the various authentication methods available (Note only Token and Local Supported).",
	}
	err := survey.AskOne(q, &selectedAuth, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("error asking for authentication method: %v", err)
	}
	log.Debugf("selected auth: %v", selectedAuth)
	return selectedAuth, nil
}

func askForUserCredentials() (*userCredentials, error) {
	answers := userCredentials{}
	qs := []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Enter your username",
			},
			Validate: survey.Required,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Enter your password",
			},
			Validate: survey.Required,
		},
	}
	if err := survey.Ask(qs, &answers); err != nil {
		return nil, fmt.Errorf("error fetching credentials: %v", err)
	}
	return &answers, nil
}

func getTokenAuthCredentials() (string, error) {
	// TODO: Add help as to what kind of token is needed (Global scope, etc)
	token := ""
	q := &survey.Input{
		Message: "Enter your token:",
		Help:    "Enter your token (Example: token-abcde:somevalue)",
	}
	err := survey.AskOne(q, &token, survey.WithValidator(survey.Required))
	if err != nil {
		return "", fmt.Errorf("error asking for token: %v", err)
	}
	// TODO: Add validator for token format
	return token, nil
}

