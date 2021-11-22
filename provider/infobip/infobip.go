package infobip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

//Config configuration struct for Infobip Client
type Config struct {
	Token  string `yaml:"token"`
	Secret string `yaml:"secret"`
}

//InfobipRequestTimeout  is the timeout for http request to Infobip
const InfobipRequestTimeout = time.Second * 20

//Infobip is the exte Infobip
type Infobip struct {
	Config
}

type InfobipDestination struct {
    To string `json:"to"`
}

type InfobipMessage struct {
    From string `json:"from"`
    Destinations []InfobipDestination `json:"destinations"`
    Text string `json:"text"`
}

type InfobipPayload struct {
    Messages []InfobipMessage `json:"messages"`
}

//NewInfobip creates a new
func NewInfobip(config Config) *Infobip {
	Infobip := &Infobip{config}
	return Infobip
}

//Send send sms to n number of people using bulk sms api
func (c *Infobip) Send(message sachet.Message) (err error) {
	smsURL := "https://api.infobip.com/sms/2/text/advanced"
	//smsURL = "http://requestb.in/pwf2ufpw"
	var request *http.Request
	var resp *http.Response

    payload := InfobipPayload{}
    payload.Messages = append(payload.Messages, InfobipMessage{})
    payload.Messages[0].From = message.From
    payload.Messages[0].Text = message.Text

    for _, destination := range message.To {
        payload.Messages[0].Destinations = append(
            payload.Messages[0].Destinations,
            InfobipDestination{
                To: destination,
            },
        )
    }

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	//preparing the request
	request, err = http.NewRequest("POST", smsURL, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	request.SetBasicAuth(c.Token, c.Secret)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "SachetV1.0")
	//calling the endpoint
	httpClient := &http.Client{}
	httpClient.Timeout = InfobipRequestTimeout

	resp, err = httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var body []byte
	resp.Body.Read(body)
	if resp.StatusCode == http.StatusOK && err == nil {
		return
	}
	return fmt.Errorf("Failed sending sms:Reason: %s , StatusCode : %d", string(body), resp.StatusCode)
}
