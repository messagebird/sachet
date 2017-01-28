package sachet

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

//InfobipConfig configuration struct for Infobip Client
type InfobipConfig struct {
	Token  string `yaml:"token"`
	Secret string `yaml:"secret"`
}

//InfobipRequestTimeout  is the timeout for http request to Infobip
const InfobipRequestTimeout = time.Second * 20

//Infobip is the exte Infobip
type Infobip struct {
	InfobipConfig
}

//NewInfobip creates a new
func NewInfobip(config InfobipConfig) *Infobip {
	Infobip := &Infobip{config}
	return Infobip
}

//Send send sms to n number of people using bulk sms api
func (c *Infobip) Send(message Message) (err error) {
	smsURL := "https://api.infobip.com/sms/1/text/single"
	var request *http.Request
	var resp *http.Response

	data := fmt.Sprintf(`{"from": "%s",
			"to" : "%s",
			"text" : "%s"}`, message.From, message.To[0], message.Text)

	//preparing the request
	request, err = http.NewRequest("POST", smsURL, bytes.NewBufferString(data))
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
