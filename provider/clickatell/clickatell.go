package clickatell

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
	"github.com/messagebird/sachet"
)

type ClickatellConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	ApiId    string `yaml:"api_id"`
}

const ClickatellRequestTimeout = time.Second * 60

type Clickatell struct {
	User     string
	Password string
	ApiId    string
}

func NewClickatell(config ClickatellConfig) *Clickatell {
	Clickatell := &Clickatell{User: config.User, Password: config.Password, ApiId: config.ApiId}
	return Clickatell
}

func (c *Clickatell) Send(message sachet.Message) (err error) {
	for _, number := range message.To {
		err = c.SendOne(message, number)
		if err != nil {
			return fmt.Errorf("Failed to make API call to clickatell:%s", err)
		}
	}
	return
}

func (c *Clickatell) SendOne(message sachet.Message, PhoneNumber string) (err error) {
	fmt.Printf("ALERT : %s\n", message.Text)
	encoded_message := url.QueryEscape(message.Text)
	smsURL := fmt.Sprintf("http://api.clickatell.com/http/sendmsg?user=%s&password=%s&api_id=%s&to=%s&text=%s", c.User, c.Password, c.ApiId, PhoneNumber, encoded_message)
	var request *http.Request
	var resp *http.Response
	request, err = http.NewRequest("GET", smsURL, nil)
	if err != nil {
		return
	}
	httpClient := &http.Client{}
	httpClient.Timeout = ClickatellRequestTimeout
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
	return fmt.Errorf("Failed sending sms:Reason: %s, StatusCode : %d", string(body), resp.StatusCode)
}
