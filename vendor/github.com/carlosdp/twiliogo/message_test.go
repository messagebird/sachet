package twiliogo

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testMessage = Message{
	Sid:         "testsid",
	DateCreated: "2013-05-11",
	DateUpdated: "2013-05-11",
	DateSent:    "2013-05-11",
	AccountSid:  "AC3TestAccount",
	From:        "+15555555555",
	To:          "+16666666666",
	Body:        "TestBody",
	NumSegments: "1",
	Status:      "queued",
	Direction:   "outbound-api",
	Price:       "4",
	PriceUnit:   "dollars",
	ApiVersion:  "2008-04-01",
	Uri:         "/2010-04-01/Accounts/AC3TestAccount/Messages/testsid.json",
}

func TestNewMessage(t *testing.T) {
	client := new(MockClient)

	messageJson, _ := json.Marshal(testMessage)

	params := url.Values{}
	params.Set("From", "6666666666")
	params.Set("To", "5555555555")
	params.Set("Body", "TestBody")

	client.On("post", params, "/Messages.json").Return(messageJson, nil)

	message, err := NewMessage(client, "6666666666", "5555555555", Body("TestBody"))

	client.Mock.AssertExpectations(t)

	if assert.Nil(t, err, "Error unmarshaling message") {
		assert.Equal(t, message.Sid, "testsid", "Message malformed")
	}
}

func TestGetMessage(t *testing.T) {
	client := new(MockClient)

	messageJson, _ := json.Marshal(testMessage)

	client.On("get", url.Values{}, "/Messages/testsid.json").Return(messageJson, nil)

	message, err := GetMessage(client, "testsid")

	client.Mock.AssertExpectations(t)

	if assert.Nil(t, err, "Error unmarshaling message") {
		assert.Equal(t, message.Sid, "testsid", "Message malformed")
	}
}
