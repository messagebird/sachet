package sachet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CMConfig configuration struct for CM Client
type CMConfig struct {
	ProductToken string `yaml:"producttoken"`
}

// CM is the exte CM
type CM struct {
	CMConfig
}

var cmHTTPClient = &http.Client{Timeout: time.Second * 20}

//NewCM creates a new
func NewCM(config CMConfig) *CM {
	return &CM{config}
}

// Send send SMS to n number of people using Bulk SMS API
func (c *CM) Send(message Message) error {
	smsURL := "https://gw.cmtelecom.com/v1.0/message"

	var dataMap map[string]interface{}

	dataMap["Authentication"] = map[string]string{"ProductToken": c.CMConfig.ProductToken}
	dataMap["From"] = message.From

	var to []map[string]string
	for _, recipient := range message.To {
		to = append(to, map[string]string{"Number": recipient})
	}
	dataMap["To"] = to

	dataMap["Body"] = map[string]string{"Content": message.Text}

	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", smsURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Sachet")

	response, err := cmHTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var body []byte
	response.Body.Read(body)
	if response.StatusCode == http.StatusOK && err == nil {
		return nil
	}

	return fmt.Errorf("Failed sending sms. Reason: %s, statusCode: %d", string(body), response.StatusCode)
}
