package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
)

type TableauClient struct {
	ApiUrl     string
	HTTPClient *http.Client
	AuthToken  string
}

type Site struct {
	ID         string `json:"id"`
	ContentUrl string `json:"contentUrl"`
}

type Credentials struct {
	TokenName   string `json:"personalAccessTokenName"`
	TokenSecret string `json:"personalAccessTokenSecret"`
	Site        Site   `json:"site"`
}

type SignInRequest struct {
	Credentials Credentials `json:"credentials"`
}

type SignInResponseData struct {
	Site                      Site   `json:"site"`
	Token                     string `json:"token"`
	EstimatedTimeToExpiration string `json:"estimatedTimeToExpiration"`
}

type SignInResponse struct {
	SignInResponseData SignInResponseData `json:"credentials"`
}

func NewTableauClient(serverAddress string, apiVersion string, site string, personalAccessTokenName string, personalAccessTokenSecret string) (*TableauClient, error) {
	tableauClient := &TableauClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	baseUrl := fmt.Sprintf("%s/api/%s", serverAddress, apiVersion)
	signInUrl := fmt.Sprintf("%s/auth/signin", baseUrl)

	// Create credentials
	credentials := Credentials{
		TokenName:   personalAccessTokenName,
		TokenSecret: personalAccessTokenSecret,
		Site: Site{
			ContentUrl: site,
		},
	}

	// Create sign in request
	authRequest := SignInRequest{
		Credentials: credentials,
	}

	// Marshal sign in request to JSON
	authRequestJson, err := json.Marshal(authRequest)
	if err != nil {
		return nil, err
	}

	// authenticate
	req, err := http.NewRequest("POST", signInUrl, strings.NewReader(string(authRequestJson)))
	if err != nil {
		return nil, err
	}

	// send request to Tableau API Server
	body, err := tableauClient.sendRequest(req)
	if err != nil {
		return nil, err
	}

	// Unmarshal response
	var signInResponse SignInResponse
	err = json.Unmarshal(body, &signInResponse)
	if err != nil {
		return nil, err
	}

	// Set API URL
	tableauClient.ApiUrl = fmt.Sprintf("%s/sites/%s", baseUrl, signInResponse.SignInResponseData.Site.ID)

	// Set auth token
	tableauClient.AuthToken = signInResponse.SignInResponseData.Token

	return tableauClient, nil
}

func (c *TableauClient) sendRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Tableau-Auth", c.AuthToken)

	body, err := retry.DoWithData(
		func() ([]byte, error) {
			res, err := c.HTTPClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			if (res.StatusCode != http.StatusOK) && (res.StatusCode != 201) && (res.StatusCode != 204) {
				return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
			}

			return body, nil
		},
		retry.Attempts(3),
		retry.Delay(5*time.Second),
	)

	if err != nil {
		return nil, err
	}

	return body, nil
}
