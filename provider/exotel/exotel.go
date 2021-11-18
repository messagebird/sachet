package exotel

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

//Config configuration struct for exotel Client
type Config struct {
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
}

//ExotelRequestTimeout  is the timeout for http request to exotel
const ExotelRequestTimeout = time.Second * 20

//Exotel is the exte Exotel
type Exotel struct {
	AccountSid string
	Token      string
}

//NewExotel creates a new
func NewExotel(config Config) *Exotel {
	Exotel := &Exotel{AccountSid: config.AccountSID, Token: config.AuthToken}
	return Exotel
}

//Send send sms to n number of people using bulk sms api
func (c *Exotel) Send(message sachet.Message) (err error) {
	smsURL := fmt.Sprintf("https://twilix.exotel.in/v1/Accounts/%s/Sms/send.json", c.AccountSid)
	var request *http.Request
	var resp *http.Response

	form := url.Values{"From": {message.From}, "Body": {message.Text}, "To": message.To}

	//preparing the request
	request, err = http.NewRequest("POST", smsURL, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	request.SetBasicAuth(c.AccountSid, c.Token)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "SachetV1.0")
	//calling the endpoint
	httpClient := &http.Client{}
	httpClient.Timeout = ExotelRequestTimeout

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
