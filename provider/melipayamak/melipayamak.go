package melipayamak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	Endpoint string `yaml:"endpoint"`
}

var _ (sachet.Provider) = (*Melipayamak)(nil)

type Melipayamak struct {
	Config
	HTTPClient *http.Client
}

func NewMelipayamak(config Config) *Melipayamak {
	return &Melipayamak{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

func (mp *Melipayamak) Send(message sachet.Message) error {

	Payload := map[string]string{
		"username": mp.Username,
		"password": mp.Password,
		"to":       strings.Join(message.To, ","),
		"from":     message.From,
		"text":     message.Text,
	}
	data, _ := json.Marshal(Payload)
	request, err := http.NewRequest("POST", mp.Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("content-type", "application/json")
	request.Header.Set("User-Agent", "Sachet")
	response, err := mp.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode != http.StatusOK {

		return fmt.Errorf(
			"SMS sending failed. HTTP status code: %d, Response body: %s",
			response.StatusCode,
			body,
		)
	}
	return nil
}
