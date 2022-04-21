package sfr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

// Config is the configuration struct for Sfr provider.
type Config struct {
	URL             string `yaml:"url"`
	SPACEID         string `yaml:"space_id"`
	SERVICEID       string `yaml:"service_id"`
	SERVICEPASSWORD string `yaml:"service_password"`
	LANG            string `yaml:"lang"`
	TPOA            string `yaml:"tpoa"`
}

type Authenticate struct {
	ServiceId       string `json:"serviceId"`
	ServicePassword string `json:"servicePassword"`
	SpaceId         string `json:"spaceId"`
	Lang            string `json:"lang"`
}

type MessageUnitaire struct {
	Media   string `json:"media"`
	TextMsg string `json:"textMsg"`
	To      string `json:"to"`
	From    string `json:"from"`
}

type ResponseBody struct {
	Success       bool   `json:"success"`
	ErrorCode     string `json:"errorCode"`
	ErrorDetail   string `json:"errorDetail"`
	Fatal         bool   `json:"fatal"`
	InvalidParams bool   `json:"invalidParams"`
	Response      int64  `json:response`
}

// Sap contains the necessary values for the Sfr provider.
type Sfr struct {
	Config
	HTTPClient *http.Client // The HTTP client to send requests on
}

// NewSfr creates and returns a new Sfr struct.
func NewSfr(config Config) *Sfr {
	if config.URL == "" {
		config.URL = "https://www.dmc.sfr-sh.fr/DmcWS/1.5.7/JsonService/MessagesUnitairesWS/addSingleCall"
	}
	return &Sfr{
		config,
		&http.Client{Timeout: time.Second * 20},
	}
}

// Send sends SMS to user registered in configuration.
func (c *Sfr) Send(message sachet.Message) error {
	// No \n in Text tolerated.
	msg := strings.ReplaceAll(message.Text, "\n", " - ")

	error := 0

	for _, dest := range message.To {

		request, err := http.NewRequest("GET", c.URL, nil)
		if err != nil {
			return err
		}

		authenticate := &Authenticate{
			ServiceId:       c.SERVICEID,
			ServicePassword: c.SERVICEPASSWORD,
			SpaceId:         c.SPACEID,
			Lang:            c.LANG,
		}

		params := request.URL.Query()

		messageUnitaire := &MessageUnitaire{
			Media:   "SMSLong",
			TextMsg: msg,
			To:      dest,
			From:    c.TPOA,
		}

		a, err := json.Marshal(authenticate)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		mU, err := json.Marshal(messageUnitaire)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		params.Add("authenticate", string(a))
		params.Add("messageUnitaire", string(mU))
		request.URL.RawQuery = params.Encode()

		response, err := c.HTTPClient.Do(request)
		if err != nil {
			error += 1
			fmt.Println(err)
		}
		defer response.Body.Close()
		dec := json.NewDecoder(response.Body)

		var responseBody ResponseBody
		for dec.More() {
			err := dec.Decode(&responseBody)
			if err != nil {
				fmt.Errorf("Can not decode JSON")
				error += 1
			}
		}

		if responseBody.Success != true {
			fmt.Println("API error :", responseBody)
			error += 1
		} else {
			fmt.Println("Successfully sent alert to ", dest)
		}
	}
	if error > 0 {
		return fmt.Errorf("Error with %d calls", error)
	}
	return nil
}
