package kannel

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/messagebird/sachet"
)

// Config configuration struct for Kannel Client.
type Config struct {
	URL  string `yaml:"url"`
	User string `yaml:"username"`
	Pass string `yaml:"password"`
}

// KannelRequestTimeout  is the timeout for http request to Kannel.
const KannelRequestTimeout = time.Second * 20

// Kannel is the exte Kannel.
type Kannel struct {
	Config
}

// NewKannel creates a new.
func NewKannel(config Config) *Kannel {
	Kannel := &Kannel{config}
	return Kannel
}

// Send send sms to n number of people using bulk sms api.
func (c *Kannel) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		queryParams := url.Values{
			"from": {message.From},
			"to":   {recipient},
			"text": {message.Text},
			"user": {c.User},
			"pass": {c.Pass},
		}

		request, err := http.NewRequest("GET", c.URL, nil)
		if err != nil {
			return err
		}

		request.URL.RawQuery = queryParams.Encode()
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Set("User-Agent", "SachetV1.0")
		//	calling the endpoint - print out Kannel requested URL for debug purpose
		// fmt.Println(request.URL.String())
		httpClient := &http.Client{}
		httpClient.Timeout = KannelRequestTimeout

		response, err := httpClient.Do(request)
		if err != nil {
			return err
		}

		if response.StatusCode >= http.StatusBadRequest {
			return fmt.Errorf("Failed sending sms. statusCode: %d", response.StatusCode)
		}
	}

	return nil
}
