package twiliogo

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testNumber = IncomingPhoneNumber{
	Sid:                  "testsid",
	AccountSid:           "AC3TestAccount",
	FriendlyName:         "testname",
	PhoneNumber:          "+1 (444) 444-4444",
	VoiceUrl:             "http://test.com",
	VoiceMethod:          "POST",
	VoiceFallbackUrl:     "http://fail.com",
	VoiceFallbackMethod:  "GET",
	VoiceCallerIdLookup:  true,
	StatusCallback:       "http://status.com",
	StatusCallbackMethod: "GET",
	SmsUrl:               "http://sms.com",
	SmsMethod:            "GET",
	DateCreated:          "2013-05-11",
	DateUpdated:          "2013-05-11",
	Capabilities:         Capabilites{true, true, true},
	ApiVersion:           "2008-04-01",
	Uri:                  "/2010-04-01/Accounts/AC3TestAccount/Messages/testsid.json",
}

func TestBuyPhoneNumber(t *testing.T) {
	client := new(MockClient)

	numberJson, _ := json.Marshal(testNumber)

	params := url.Values{}
	params.Set("PhoneNumber", "4444444444")

	client.On("post", params, "/IncomingPhoneNumbers.json").Return(numberJson, nil)

	number, err := BuyPhoneNumber(client, PhoneNumber("4444444444"))

	client.Mock.AssertExpectations(t)

	if assert.Nil(t, err, "Error unmarshaling number") {
		assert.Equal(t, number.Sid, "testsid", "Number malformed")
	}
}

func TestGetIncomingPhoneNumber(t *testing.T) {
	client := new(MockClient)

	numberJson, _ := json.Marshal(testNumber)

	client.On("get", url.Values{}, "/IncomingPhoneNumbers/testsid.json").Return(numberJson, nil)

	number, err := GetIncomingPhoneNumber(client, "testsid")

	client.Mock.AssertExpectations(t)

	if assert.Nil(t, err, "Error unmarshaling number") {
		assert.Equal(t, number.Sid, "testsid", "Number malformed")
	}
}
