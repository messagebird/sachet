package otc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

type Config struct {
	IdentityEndpoint string `yaml:"identity_endpoint"`
	DomainName       string `yaml:"domain_name"`
	ProjectName      string `yaml:"project_name"`
	UserName         string `yaml:"username"`
	Password         string `yaml:"password"`
	ProjectID        string `yaml:"project_id"`
	Insecure         bool   `yaml:"insecure"`
	token            string
	otcBaseURL       string
}

type smsRequest struct {
	Endpoint string `json:"endpoint"`
	Message  string `json:"message"`
}

type OTC struct {
	Config
}

func NewOTC(config Config) *OTC {
	OTC := &OTC{config}
	return OTC
}

func (c *OTC) loginRequest() error {
	type nameResponse struct {
		Name string `json:"name"`
	}

	type userResponse struct {
		Name     string       `json:"name"`
		Password string       `json:"password"`
		Domain   nameResponse `json:"domain"`
	}

	type passwordResponse struct {
		User userResponse `json:"user"`
	}
	type identityResponse struct {
		Methods  []string         `json:"methods"`
		Password passwordResponse `json:"password"`
	}

	type scopeResponse struct {
		Project nameResponse `json:"project"`
	}

	type authResponse struct {
		Identity identityResponse `json:"identity"`
		Scope    scopeResponse    `json:"scope"`
	}

	type loginResponse struct {
		Auth authResponse `json:"auth"`
	}

	userResp := userResponse{
		Name:     c.UserName,
		Password: c.Password,
		Domain: nameResponse{
			Name: c.DomainName,
		},
	}

	loginResp := loginResponse{
		Auth: authResponse{
			Identity: identityResponse{
				Methods: []string{"password"},
				Password: passwordResponse{
					User: userResp,
				},
			},
			Scope: scopeResponse{
				Project: nameResponse{
					Name: c.ProjectName,
				},
			},
		},
	}

	body, err := json.Marshal(loginResp)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.IdentityEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	tr := http.DefaultTransport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: c.Insecure}

	client := &http.Client{Timeout: 10 * time.Second, Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("OTC API request failed with HTTP status code %d", resp.StatusCode)
	}

	c.token = resp.Header.Get("X-Subject-Token")

	if c.token == "" {
		return fmt.Errorf("unable to get auth token")
	}

	type endpointResponse struct {
		Token struct {
			Catalog []struct {
				Type      string `json:"type"`
				Endpoints []struct {
					URL       string `json:"url"`
					Interface string `json:"interface"`
					Region    string `json:"region"`
				} `json:"endpoints"`
			} `json:"catalog"`
		} `json:"token"`
	}
	var endpointResp endpointResponse

	err = json.NewDecoder(resp.Body).Decode(&endpointResp)
	if err != nil {
		return err
	}

	for _, v := range endpointResp.Token.Catalog {
		if v.Type == "smn" {
			for _, endpoint := range v.Endpoints {
				c.otcBaseURL = fmt.Sprintf("%s%s", endpoint.URL, c.ProjectID)
				continue
			}
		}

		if c.otcBaseURL != "" {
			continue
		}
	}

	if c.otcBaseURL == "" {
		return fmt.Errorf("unable to find snm endpoint")
	}

	return nil
}

func (c *OTC) SendRequest(method, resource string, payload *smsRequest, attempts int) (io.Reader, error) {
	if len(c.token) == 0 {
		err := c.loginRequest()

		if err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf("%s/%s", c.otcBaseURL, resource)
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", c.token)

	tr := http.DefaultTransport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: c.Insecure}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		// Set empty token to force login
		c.token = ""
		if attempts--; attempts > 0 {
			return c.SendRequest(method, resource, payload, attempts)
		} else {
			return nil, err
		}
	} else if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("OTC API request %s failed with HTTP status code %d", url, resp.StatusCode)
	}

	body1, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body1), nil
}

//Send send sms to n number of people using bulk sms api
func (c *OTC) Send(message sachet.Message) (err error) {
	for _, recipent := range message.To {

		r1 := &smsRequest{
			Endpoint: recipent,
			Message:  message.Text,
		}
		_, err := c.SendRequest("POST", "notifications/sms", r1, 2)
		if err != nil {
			return err
		}
	}
	return nil
}
