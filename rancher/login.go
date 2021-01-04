package rancher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//"github.com/davecgh/go-spew/spew"
	"github.com/rancher/log"
	v3public "github.com/rancher/types/client/management/v3public"
)

const (
	FiveMinutes      = 5 * 60 * 1000
	TokenDescription = "Token for logs-collector-rancher2"

	contentType = "application/json"
)

func DoLocalLogin(url, username, password string) (string, error) {
	loginPostURL := url + "/v3-public/localProviders/local?action=login"
	return DoLogin(loginPostURL, username, password)
}

func DoLogin(loginPostURL, username, password string) (string, error) {
	payload := map[string]interface{}{
		"username":    username,
		"password":    password,
		"description": TokenDescription,
		"ttl":         FiveMinutes,
	}

	marshaledBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Debugf("marshaledBytes: %v", string(marshaledBytes))
	// TODO: Deal with self signed certs
	//if InsecureClient != nil {
	//	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//}

	log.Debugf("loginPostURL: %v", loginPostURL)
	resp, err := http.Post(loginPostURL,
		contentType,
		bytes.NewBuffer(marshaledBytes),
	)
	if err != nil {
		log.Errorf("error logging in: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Debugf("resp.Body: %+v", resp.Body)
	log.Debugf("resp.Status: %+v", resp.Status)
	log.Debugf("resp.StatusCode: %+v", resp.StatusCode)
	//var result map[string]interface{}
	result := v3public.Token{}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshaling body: %v", err)
	}
	log.Debugf("result: %+v", result)

	if resp.StatusCode != http.StatusCreated {
		// TODO: Figure out the error message from the response object
		return "", fmt.Errorf("unexpected status code %v, message: TODO", resp.StatusCode)
	}

	log.Debugf("Logged in and rancher returned userId result as %+v, token: %v", result.UserID, result.Token)
	return result.Token, nil
}

func (s *Server) doLogin() error {
	// Check what kind of auth providers are available
	authProviders, err := s.V3PublicClient.AuthProvider.List(nil)
	if err != nil {
		return fmt.Errorf("error listing auth providers: %v", err)
	}
	log.Debugf("found %v authProviders", len(authProviders.Data))

	authOptions := []string{"token"}
	for _, authProvider := range authProviders.Data {
		log.Debugf("AuthProvider: %+v", authProvider)
		authOptions = append(authOptions, authProvider.ID)
	}

	selectedAuth, err := askForAuthenticationMethod(authOptions)
	if err != nil {
		return err
	}
	// TODO: deal with other auth types like github, etc
	switch selectedAuth {
	// TODO: Introduce constants
	case authTypeLocal:
		credentials, err := askForUserCredentials()
		if err != nil {
			return err
		}
		token, err := DoLocalLogin(s.URL, credentials.Username, credentials.Password)
		if err != nil {
			return err
		}
		log.Infof("Successfully logged in and generated a token for the session")
		// TODO: Should we delete the token after the session is done?

		// TODO: Delete this
		log.Debugf("token: %+v", token)
		s.Token = token

	case authTypeToken:
		token, err := getTokenAuthCredentials()
		if err != nil {
			return err
		}
		s.Token = token

	default:
		return fmt.Errorf("unhandled authentication %v", selectedAuth)
	}

	// Assumption at this point is, the token field is populated,
	// no matter what authentication mechanism and we work with it

	// TODO:
	log.Infof("TODO: improve me: checking token")
	// TODO: Fetch /v3/users?me=true to make sure the token is valid
	// 		 and also to get the name of the user to personalize the
	//		 experience.
	return nil
}
