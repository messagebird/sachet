package twiliogo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const ROOT = "https://api.twilio.com"
const VERSION = "2010-04-01"

type Client interface {
	AccountSid() string
	AuthToken() string
	RootUrl() string
	get(url.Values, string) ([]byte, error)
	post(url.Values, string) ([]byte, error)
	delete(string) error
}

type TwilioClient struct {
	accountSid string
	authToken  string
	rootUrl    string
}

var _ Client = &TwilioClient{}

func NewClient(accountSid, authToken string) *TwilioClient {
	rootUrl := ROOT + "/" + VERSION + "/Accounts/" + accountSid
	return &TwilioClient{accountSid, authToken, rootUrl}
}

func (client *TwilioClient) post(values url.Values, uri string) ([]byte, error) {
	req, err := http.NewRequest("POST", client.buildUri(uri), strings.NewReader(values.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.AccountSid(), client.AuthToken())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient := &http.Client{}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return body, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		if res.StatusCode == 500 {
			return body, Error{"Server Error"}
		} else {
			twilioError := new(TwilioError)
			json.Unmarshal(body, twilioError)
			return body, twilioError
		}
	}

	return body, err
}

func (client *TwilioClient) get(queryParams url.Values, uri string) ([]byte, error) {
	var params *strings.Reader

	if queryParams == nil {
		queryParams = url.Values{}
	}

	params = strings.NewReader(queryParams.Encode())
	req, err := http.NewRequest("GET", client.buildUri(uri), params)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(client.AccountSid(), client.AuthToken())
	httpClient := &http.Client{}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return body, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		if res.StatusCode == 500 {
			return body, Error{"Server Error"}
		} else {
			twilioError := new(TwilioError)
			json.Unmarshal(body, twilioError)
			return body, twilioError
		}
	}

	return body, err
}

func (client *TwilioClient) delete(uri string) error {
	req, err := http.NewRequest("DELETE", client.buildUri(uri), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(client.AccountSid(), client.AuthToken())
	httpClient := &http.Client{}

	res, err := httpClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != 204 {
		return fmt.Errorf("Non-204 returned from server for DELETE: %d", res.StatusCode)
	}

	return nil
}

func (client *TwilioClient) AccountSid() string {
	return client.accountSid
}

func (client *TwilioClient) AuthToken() string {
	return client.authToken
}

func (client *TwilioClient) RootUrl() string {
	return client.rootUrl
}

func (client *TwilioClient) buildUri(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	newParts := make([]string, 0, len(parts))
	// Check for "http" because sometimes we get raw URLs from following the metadata.
	if !strings.HasPrefix(parts[0], "http") {
		newParts = append(newParts, client.RootUrl())
	}
	for _, p := range parts {
		p = strings.Trim(p, "/")
		if p == "" {
			continue
		}
		newParts = append(newParts, p)
	}
	return strings.Join(newParts, "/")
}
