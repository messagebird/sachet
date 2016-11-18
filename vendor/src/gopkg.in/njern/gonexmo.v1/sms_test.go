// +build messages

// Tests in this file will only be run if the build tag messages is set:
// `go test -tag messages`
// Test with only sending one message using:
// `go test -test.run SendText -tags messages`
package nexmo

import (
	"strconv"
	"testing"
	"time"
)

// TODO(inhies): Only create a Client once in an init() function.

func TestSendTextMessage(t *testing.T) {
	// TODO(inhies): Create an internal rate limiting system and do away with
	// this hacky 1 second delay.
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &SMSMessage{
		From:            TEST_FROM,
		To:              TEST_PHONE_NUMBER,
		Type:            Text,
		Text:            "Gonexmo test SMS message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           Standard,
	}

	messageResponse, err := nexmo.SMS.Send(message)
	if err != nil {
		t.Error("Failed to send text message with error:", err)
	}

	t.Logf("Sent SMS, response was: %#v\n", messageResponse)
}

func TestFlashMessage(t *testing.T) {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &SMSMessage{
		From:            TEST_FROM,
		To:              TEST_PHONE_NUMBER,
		Type:            Text,
		Text:            "Gonexmo test flash SMS message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           Flash,
	}

	messageResponse, err := nexmo.SMS.Send(message)
	if err != nil {
		t.Error("Failed to send flash message (class 0 SMS) with error:", err)
	}

	t.Logf("Sent Flash SMS, response was: %#v\n", messageResponse)
}
