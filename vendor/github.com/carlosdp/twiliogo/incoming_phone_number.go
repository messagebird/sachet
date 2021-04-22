package twiliogo

import (
	"encoding/json"
	"net/url"
)

type Capabilites struct {
	Voice bool `json:"voice"`
	SMS   bool `json:"SMS"`
	MMS   bool `json:"MMS"`
}

type IncomingPhoneNumber struct {
	Sid                  string      `json:"sid"`
	AccountSid           string      `json:"account_sid"`
	FriendlyName         string      `json:"friendly_name"`
	PhoneNumber          string      `json:"phone_number"`
	VoiceUrl             string      `json:"voice_url"`
	VoiceMethod          string      `json:"voice_method"`
	VoiceFallbackUrl     string      `json:"voice_fallback_url"`
	VoiceFallbackMethod  string      `json:"voice_fallback_method"`
	StatusCallback       string      `json:"status_callback"`
	StatusCallbackMethod string      `json:"status_callback_method"`
	VoiceCallerIdLookup  bool        `json:"voice_caller_id_lookup"`
	VoiceApplicationId   string      `json:"voice_application_id"`
	DateCreated          string      `json:"date_created"`
	DateUpdated          string      `json:"date_updated"`
	SmsUrl               string      `json:"sms_url"`
	SmsMethod            string      `json:"sms_method"`
	SmsFallbackUrl       string      `json:"sms_fallback_url"`
	SmsFallbackMethod    string      `json:"sms_fallback_method"`
	SmsApplicationId     string      `json:"sms_application_id"`
	Capabilities         Capabilites `json:"capabilities"`
	ApiVersion           string      `json:"api_version"`
	Uri                  string      `json:"uri"`
}

func GetIncomingPhoneNumber(client Client, sid string) (*IncomingPhoneNumber, error) {
	var incomingPhoneNumber *IncomingPhoneNumber

	res, err := client.get(url.Values{}, "/IncomingPhoneNumbers/"+sid+".json")

	if err != nil {
		return nil, err
	}

	incomingPhoneNumber = new(IncomingPhoneNumber)
	err = json.Unmarshal(res, incomingPhoneNumber)

	return incomingPhoneNumber, err
}

func BuyPhoneNumber(client Client, number Optional) (*IncomingPhoneNumber, error) {
	var incomingPhoneNumber *IncomingPhoneNumber

	if number == nil {
		return nil, Error{"Must input PhoneNumber or AreaCode"}
	}

	params := url.Values{}
	param, value := number.GetParam()
	params.Set(param, value)

	res, err := client.post(params, "/IncomingPhoneNumbers.json")

	if err != nil {
		return incomingPhoneNumber, err
	}

	incomingPhoneNumber = new(IncomingPhoneNumber)
	err = json.Unmarshal(res, incomingPhoneNumber)

	return incomingPhoneNumber, err
}
