package messagebird

import (
	"testing"
	"time"
)

var voiceMessageObject []byte = []byte(`{
  "id":"430c44a0354aab7ac9553f7a49907463",
  "href":"https:\/\/rest.messagebird.com\/voicemessages\/430c44a0354aab7ac9553f7a49907463",
  "originator":"MessageBird",
  "body":"Hello World",
  "reference":null,
  "language":"en-gb",
  "voice":"female",
  "repeat":1,
  "ifMachine":"continue",
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T16:11:24+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"calling",
        "statusDatetime":"2015-01-05T16:11:24+00:00"
      }
    ]
  }
}`)

func TestNewVoiceMessage(t *testing.T) {
	SetServerResponse(200, voiceMessageObject)

	message, err := mbClient.NewVoiceMessage([]string{"31612345678"}, "Hello World", nil)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new voice message: %s", err)
	}

	if message.Id != "430c44a0354aab7ac9553f7a49907463" {
		t.Errorf("Unexpected voice message id: %s", message.Id)
	}

	if message.HRef != "https://rest.messagebird.com/voicemessages/430c44a0354aab7ac9553f7a49907463" {
		t.Errorf("Unexpected voice message href: %s", message.HRef)
	}

	if message.Originator != "MessageBird" {
		t.Errorf("Unexpected voice message originator: %s", message.Originator)
	}

	if message.Body != "Hello World" {
		t.Errorf("Unexpected voice message body: %s", message.Body)
	}

	if message.Reference != "" {
		t.Errorf("Unexpected voice message reference: %s", message.Reference)
	}

	if message.Language != "en-gb" {
		t.Errorf("Unexpected voice message language: %s", message.Language)
	}

	if message.Voice != "female" {
		t.Errorf("Unexpected voice message voice: %s", message.Voice)
	}

	if message.Repeat != 1 {
		t.Errorf("Unexpected voice message repeat: %d", message.Repeat)
	}

	if message.IfMachine != "continue" {
		t.Errorf("Unexpected voice message ifmachine: %d", message.IfMachine)
	}

	if message.ScheduledDatetime != nil {
		t.Errorf("Unexpected voice message scheduled datetime: %s", message.ScheduledDatetime)
	}

	if message.CreatedDatetime == nil || message.CreatedDatetime.Format(time.RFC3339) != "2015-01-05T16:11:24Z" {
		t.Errorf("Unexpected voice message created datetime: %s", message.CreatedDatetime)
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 1 {
		t.Errorf("Unexpected number of total sent count: %d", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected voice message recipient: %d", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "calling" {
		t.Errorf("Unexpected voice message recipient status: %s", message.Recipients.Items[0].Status)
	}

	if message.Recipients.Items[0].StatusDatetime == nil || message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339) != "2015-01-05T16:11:24Z" {
		t.Errorf("Unexpected datetime status for voice message recipient: %s", message.Recipients.Items[0].StatusDatetime.Format(time.RFC3339))
	}

	if len(message.Errors) != 0 {
		t.Errorf("Unexpected number of errors in voice message: %d", len(message.Errors))
	}
}

var voiceMessageObjectWithParams []byte = []byte(`{
  "id":"430c44a0354aab7ac9553f7a49907463",
  "href":"https:\/\/rest.messagebird.com\/voicemessages\/430c44a0354aab7ac9553f7a49907463",
  "body":"Hello World",
  "reference":"MyReference",
  "language":"en-gb",
  "voice":"male",
  "repeat":5,
  "ifMachine":"hangup",
  "scheduledDatetime":null,
  "createdDatetime":"2015-01-05T16:11:24+00:00",
  "recipients":{
    "totalCount":1,
    "totalSentCount":1,
    "totalDeliveredCount":0,
    "totalDeliveryFailedCount":0,
    "items":[
      {
        "recipient":31612345678,
        "status":"calling",
        "statusDatetime":"2015-01-05T16:11:24+00:00"
      }
    ]
  }
}`)

func TestNewVoiceMessageWithParams(t *testing.T) {
	SetServerResponse(200, voiceMessageObjectWithParams)

	params := &VoiceMessageParams{
		Reference: "MyReference",
		Voice:     "male",
		Repeat:    5,
		IfMachine: "hangup",
	}

	message, err := mbClient.NewVoiceMessage([]string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new voice message: %s", err)
	}

	if message.Reference != "MyReference" {
		t.Errorf("Unexpected voice message reference: %s", message.Reference)
	}

	if message.Voice != "male" {
		t.Errorf("Unexpected voice message voice: %s", message.Voice)
	}

	if message.Repeat != 5 {
		t.Errorf("Unexpected voice message repeat: %d", message.Repeat)
	}

	if message.IfMachine != "hangup" {
		t.Errorf("Unexpected voice message ifmachine: %s", message.IfMachine)
	}
}

var voiceMessageObjectWithCreatedDatetime []byte = []byte(`{
  "id":"430c44a0354aab7ac9553f7a49907463",
  "href":"https:\/\/rest.messagebird.com\/voicemessages\/430c44a0354aab7ac9553f7a49907463",
  "body":"Hello World",
  "reference":null,
  "language":"en-gb",
  "voice":"female",
  "repeat":1,
  "ifMachine":"continue",
  "scheduledDatetime":"2015-01-05T16:12:24+00:00",
  "createdDatetime":"2015-01-05T16:11:24+00:00",
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

func TestNewVoiceMessageWithScheduledDatetime(t *testing.T) {
	SetServerResponse(200, voiceMessageObjectWithCreatedDatetime)

	scheduledDatetime, _ := time.Parse(time.RFC3339, "2015-01-05T16:12:24+00:00")

	params := &VoiceMessageParams{ScheduledDatetime: scheduledDatetime}
	message, err := mbClient.NewVoiceMessage([]string{"31612345678"}, "Hello World", params)
	if err != nil {
		t.Fatalf("Didn't expect error while creating a new voice message: %s", err)
	}

	if message.ScheduledDatetime.Format(time.RFC3339) != scheduledDatetime.Format(time.RFC3339) {
		t.Errorf("Unexpected scheduled datetime: %s", message.ScheduledDatetime.Format(time.RFC3339))
	}

	if message.Recipients.TotalCount != 1 {
		t.Fatalf("Unexpected number of total count: %d", message.Recipients.TotalCount)
	}

	if message.Recipients.TotalSentCount != 0 {
		t.Errorf("Unexpected number of total sent count: %d", message.Recipients.TotalSentCount)
	}

	if message.Recipients.Items[0].Recipient != 31612345678 {
		t.Errorf("Unexpected voice message recipient: %d", message.Recipients.Items[0].Recipient)
	}

	if message.Recipients.Items[0].Status != "scheduled" {
		t.Errorf("Unexpected voice message recipient status: %s", message.Recipients.Items[0].Status)
	}
}
