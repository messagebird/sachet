package messagebird

import (
	"testing"
	"time"
)

var messageObject []byte = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https:\/\/rest.messagebird.com\/messages\/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"sms",
  "originator":"TestName",
  "body":"Hello World",
  "reference":null,
  "validity":null,
  "gateway":239,
  "typeDetails":{
    
  },
  "datacoding":"plain",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessage(t *testing.T) {
	SetServerResponse(200, messageObject)

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", nil)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Id != "6fe65f90454aa61536e6a88b88972670" {
		t.Errorf("Unexpected message id: %s", message.Id)
	}

	if message.HRef != "https://rest.messagebird.com/messages/6fe65f90454aa61536e6a88b88972670" {
		t.Errorf("Unexpected message href: %s", message.HRef)
	}

	if message.Direction != "mt" {
		t.Errorf("Unexpected message direction: %s", message.Direction)
	}

	if message.Type != "sms" {
		t.Errorf("Unexpected message type: %s", message.Type)
	}

	if message.Originator != "TestName" {
		t.Errorf("Unexpected message originator: %s", message.Originator)
	}

	if message.Body != "Hello World" {
		t.Errorf("Unexpected message body: %s", message.Body)
	}

	if message.Reference != "" {
		t.Errorf("Unexpected message reference: %s", message.Reference)
	}

	if message.Validity != nil {
		t.Errorf("Unexpected message validity: %d", *message.Validity)
	}

	if message.Gateway != 239 {
		t.Errorf("Unexpected message gateway: %s", message.Gateway)
	}

	if len(message.TypeDetails) != 0 {
		t.Errorf("Unexpected number of items in message typedetails: %d", len(message.TypeDetails))
	}

	if message.DataCoding != "plain" {
		t.Errorf("Unexpected message datacoding: %s", message.DataCoding)
	}

	if message.MClass != 1 {
		t.Errorf("Unexpected message mclass: %s", message.MClass)
	}

	if message.ScheduledDatetime != nil {
		t.Errorf("Unexpected message scheduled datetime: %s", message.ScheduledDatetime)
	}

	if message.CreatedDatetime == nil || message.CreatedDatetime.Format(time.RFC3339) != "2015-01-05T10:02:59Z" {
		t.Errorf("Unexpected message created datetime: %s", message.CreatedDatetime)
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 1 {
		t.Errorf("Unexpected number of total sent count: %d", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected message recipient: %d", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "sent" {
		t.Errorf("Unexpected message recipient status: %s", message.Recipients.Items[0].Status)
	}

	if message.Recipients.Items[0].StatusDatetime == nil || message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339) != "2015-01-05T10:02:59Z" {
		t.Errorf("Unexpected datetime status for message recipient: %s", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}

	if len(message.Errors) != 0 {
		t.Errorf("Unexpected number of errors in message: %d", len(message.Errors))
	}
}

func TestNewMessageError(t *testing.T) {
	SetServerResponse(405, accessKeyErrorObject)

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", nil)
	if err != ErrResponse {
		t.Fatalf("Expected ErrResponse to be returned, instead I got %s", err)
	}

	if len(message.Errors) != 1 {
		t.Fatalf("Unexpected number of errors: %d", len(message.Errors))
	}

	if message.Errors[0].Code != 2 {
		t.Errorf("Unexpected error code: %d", message.Errors[0].Code)
	}

	if message.Errors[0].Parameter != "access_key" {
		t.Errorf("Unexpected error parameter: %s", message.Errors[0].Parameter)
	}
}

var messageWithParamsObject []byte = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https:\/\/rest.messagebird.com\/messages\/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"sms",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{
    
  },
  "datacoding":"unicode",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithParams(t *testing.T) {
	SetServerResponse(200, messageWithParamsObject)

	params := &MessageParams{
		Type:       "sms",
		Reference:  "TestReference",
		Validity:   13,
		Gateway:    10,
		DataCoding: "unicode",
	}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "sms" {
		t.Errorf("Unexpected message type: %s", message.Type)
	}

	if message.Reference != "TestReference" {
		t.Errorf("Unexpected message reference: %s", message.Reference)
	}

	if *message.Validity != 13 {
		t.Errorf("Unexpected message validity: %d", *message.Validity)
	}

	if message.Gateway != 10 {
		t.Errorf("Unexpected message gateway: %s", message.Gateway)
	}

	if message.DataCoding != "unicode" {
		t.Errorf("Unexpected message datacoding: %s", message.DataCoding)
	}
}

var binaryMessageObject []byte = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https:\/\/rest.messagebird.com\/messages\/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"binary",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{
    "udh":"050003340201"
  },
  "datacoding":"unicode",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithBinaryType(t *testing.T) {
	SetServerResponse(200, binaryMessageObject)

	params := &MessageParams{
		Type:        "binary",
		TypeDetails: TypeDetails{"udh": "050003340201"},
	}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "binary" {
		t.Errorf("Unexpected message type: %s", message.Type)
	}

	if len(message.TypeDetails) != 1 {
		t.Fatalf("Unexpected number of message typedetails: %d", len(message.TypeDetails))
	}

	if message.TypeDetails["udh"] != "050003340201" {
		t.Errorf("Unexpected 'udh' value in message typedetails: %s", message.TypeDetails["udh"])
	}
}

var premiumMessageObject []byte = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https:\/\/rest.messagebird.com\/messages\/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"premium",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{
    "tariff":150,
    "shortcode":1008,
    "keyword":"RESTAPI"
  },
  "datacoding":"unicode",
  "mclass":1,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithPremiumType(t *testing.T) {
	SetServerResponse(200, premiumMessageObject)

	params := &MessageParams{
		Type:        "premium",
		TypeDetails: TypeDetails{"keyword": "RESTAPI", "shortcode": 1008, "tariff": 150},
	}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "premium" {
		t.Errorf("Unexpected message type: %s", message.Type)
	}

	if len(message.TypeDetails) != 3 {
		t.Fatalf("Unexpected number of message typedetails: %d", len(message.TypeDetails))
	}

	if message.TypeDetails["tariff"] != 150.0 {
		t.Errorf("Unexpected 'tariff' value in message typedetails: %d", message.TypeDetails["tariff"])
	}

	if message.TypeDetails["shortcode"] != 1008.0 {
		t.Errorf("Unexpected 'shortcode' value in message typedetails: %d", message.TypeDetails["shortcode"])
	}

	if message.TypeDetails["keyword"] != "RESTAPI" {
		t.Errorf("Unexpected 'keyword' value in message typedetails: %s", message.TypeDetails["keyword"])
	}
}

var flashMessageObject []byte = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https:\/\/rest.messagebird.com\/messages\/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"flash",
  "originator":"TestName",
  "body":"Hello World",
  "reference":"TestReference",
  "validity":13,
  "gateway":10,
  "typeDetails":{

  },
  "datacoding":"unicode",
  "mclass":0,
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"sent",
        "statusDatetime":"2015-01-05T10:02:59+00:00"
      }
    ]
  }
}`)

func TestNewMessageWithFlashType(t *testing.T) {
	SetServerResponse(200, flashMessageObject)

	params := &MessageParams{Type: "flash"}

	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.Type != "flash" {
		t.Errorf("Unexpected message type: %s", message.Type)
	}
}

var messageObjectWithCreatedDatetime []byte = []byte(`{
  "id":"6fe65f90454aa61536e6a88b88972670",
  "href":"https:\/\/rest.messagebird.com\/messages\/6fe65f90454aa61536e6a88b88972670",
  "direction":"mt",
  "type":"sms",
  "originator":"TestName",
  "body":"Hello World",
  "reference":null,
  "validity":null,
  "gateway":239,
  "typeDetails":{
    
  },
  "datacoding":"plain",
  "mclass":1,
  "scheduledDatetime":"2015-01-05T10:03:59+00:00",
  "createdDatetime":"2015-01-05T10:02:59+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":0,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"scheduled",
        "statusDatetime":null
      }
    ]
  }
}`)

func TestNewMessageWithScheduledDatetime(t *testing.T) {
	SetServerResponse(200, messageObjectWithCreatedDatetime)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2015-01-05T10:03:59+00:00")

	params := &MessageParams{ScheduledDatetime: scheduledDatetime}
	message, err := mbClient.NewMessage("TestName", []string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new message: %s", err)
	}

	if message.ScheduledDatetime.Format(time.RFC3339) != scheduledDatetime.Format(time.RFC3339) {
		t.Errorf("Unexpected message scheduled datetime: %s", message.ScheduledDatetime.Format(time.RFC3339))
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 0 {
		t.Errorf("Unexpected number of total sent count: %d", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected message recipient: %d", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "scheduled" {
		t.Errorf("Unexpected message recipient status: %s", message.Recipients.Items[0].Status)
	}

	if message.Recipients.Items[0].StatusDatetime != nil {
		t.Errorf("Unexpected datetime status for message recipient: %s", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}
}
