package twiliogo

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCall = Call{
	Sid:            "testsid",
	ParentCallSid:  "",
	DateCreated:    "2013-05-11",
	DateUpdated:    "2013-05-11",
	AccountSid:     "AC3TestAccount",
	To:             "+15555555555",
	From:           "+16666666666",
	PhoneNumberSid: "",
	Status:         "queued",
	StartTime:      "5:00",
	EndTime:        "6:00",
	Duration:       "5",
	Price:          "4",
	PriceUnit:      "dollars",
	Direction:      "outbound-api",
	AnsweredBy:     "",
	ForwardedFrom:  "",
	CallerName:     "",
	Uri:            "/2010-04-01/Accounts/AC3TestAccount/Calls/testsid.json",
}

func TestNewCall(t *testing.T) {
	client := new(MockClient)

	callJson, _ := json.Marshal(testCall)

	params := url.Values{}
	params.Set("From", "6666666666")
	params.Set("To", "5555555555")
	params.Set("Url", "http://callback.com")

	client.On("post", params, "/Calls.json").Return(callJson, nil)

	call, err := NewCall(client, "6666666666", "5555555555", Callback("http://callback.com"))

	client.Mock.AssertExpectations(t)

	if assert.Nil(t, err, "Error unmarshaling call") {
		assert.Equal(t, call.Sid, "testsid", "Call malformed")
	}
}

func TestGetCall(t *testing.T) {
	client := new(MockClient)

	callJson, _ := json.Marshal(testCall)

	client.On("get", url.Values{}, "/Calls/testsid.json").Return(callJson, nil)

	call, err := GetCall(client, "testsid")

	client.Mock.AssertExpectations(t)

	if assert.Nil(t, err, "Error unmarshaling call") {
		assert.Equal(t, call.Sid, "testsid", "Call malformed")
	}
}
