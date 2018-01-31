package mediaburst

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/messagebird/sachet"
)

//MediaBurstConfig configuration struct for mediaburst Client
type MediaBurstConfig struct {
	APIKey string `yaml:"api_key"`
}

//MediaBurstRequestTimeout  is the timeout for http request to mediaburst
const MediaBurstRequestTimeout = time.Second * 20

//MediaBurst is the exte MediaBurst
type MediaBurst struct {
	MediaBurstConfig
}

//NewMediaBurst creates a new
func NewMediaBurst(config MediaBurstConfig) *MediaBurst {
	MediaBurst := &MediaBurst{config}
	return MediaBurst
}

//Send send sms to n number of people using bulk sms api
func (c *MediaBurst) Send(message sachet.Message) (err error) {
	smsURL := "https://api.clockworksms.com/http/send.aspx"
	var request *http.Request
	var resp *http.Response

	form := url.Values{"Key": {c.APIKey}, "From": {message.From}, "Content": {message.Text}, "To": message.To}

	//preparing the request
	request, err = http.NewRequest("GET", smsURL, strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "SachetV1.0")
	//calling the endpoint
	httpClient := &http.Client{}
	httpClient.Timeout = MediaBurstRequestTimeout

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
