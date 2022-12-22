package mattermost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

type Config struct {
	Channel  string `yaml:"channel"`
	Username string `yaml:"username"`
	IconURL  string `yaml:"icon_url"`
}

var _ (sachet.Provider) = (*Mattermost)(nil)

type Mattermost struct {
	Config
	HTTPClient *http.Client
}

func NewMattermost(config Config) *Mattermost {
	return &Mattermost{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

func (mm *Mattermost) Send(message sachet.Message) error {

	Payload := map[string]string{
		"channel":  mm.Channel,
		"username": mm.Username,
		"icon_url": mm.IconURL,
		"text":     message.Text,
	}
	data, _ := json.Marshal(Payload)
	for _, endpoint := range message.To {
		request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		request.Header.Set("content-type", "application/json")
		request.Header.Set("User-Agent", "Sachet")
		response, err := mm.HTTPClient.Do(request)
		if err != nil {
			return err
		}
		if response.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(response.Body)
			return fmt.Errorf(
				"SMS sending failed. HTTP status code: %d, Response body: %s",
				response.StatusCode,
				body,
			)
		}
		response.Body.Close()
	}

	return nil
}
