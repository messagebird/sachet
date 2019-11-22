package kannel

import (
	"fmt"
	"net/url"
	"net/http"
	"time"

	"github.com/messagebird/sachet"
)

//KannelConfig configuration struct for Kannel Client
type KannelConfig struct {
	URL 	string `yaml:"url"`
	User  string `yaml:"username"`
	Pass 	string `yaml:"password"`
}

//KannelRequestTimeout  is the timeout for http request to Kannel
const KannelRequestTimeout = time.Second * 20

//Kannel is the exte Kannel
type Kannel struct {
	KannelConfig
}

//NewKannel creates a new
func NewKannel(config KannelConfig) *Kannel {
	Kannel := &Kannel{config}
	return Kannel
}

//Send send sms to n number of people using bulk sms api
func (c *Kannel) Send(message sachet.Message) (err error) {
	var request *http.Request
	var resp *http.Response

	queryParams := url.Values{"from": {message.From}, "to": message.To, "text": {message.Text}, "user": {c.User}, "pass": {c.Pass}}
	//preparing the request
	request, err = http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return
	}

	// request.SetBasicAuth(c.User, c.Pass)
	request.URL.RawQuery = queryParams.Encode()
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "SachetV1.0")
	//calling the endpoint
	fmt.Println(request.URL.String())
	httpClient := &http.Client{}
	httpClient.Timeout = KannelRequestTimeout

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
