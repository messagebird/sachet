package nexmo

import (
	"testing"
	"time"
)

func testSend(t *testing.T) *VerifyMessageResponse {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &VerifyMessageRequest{
		Number:   TEST_PHONE_NUMBER,
		Brand:    TEST_FROM,
		SenderID: TEST_FROM,
	}

	messageResponse, err := nexmo.Verify.Send(message)
	if err != nil {
		t.Error("Failed to send verification request with error:", err)
	}

	return messageResponse
}

func TestSend(t *testing.T) {
	messageResponse := testSend(t)
	t.Logf("Sent Verification SMS, response was: %#v\n", messageResponse)
}

func TestSendCheck(t *testing.T) {
	// We need the request id, so we have to run this first.
	sendResponse := testSend(t)

	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &VerifyCheckRequest{
		RequestID: sendResponse.RequestID,
		Code:      "1122", // Take a random code here, the number will not be verified properly though.
	}

	messageResponse, err := nexmo.Verify.Check(message)
	if err != nil {
		t.Error("Failed to send verification check request with error:", err)
	}

	t.Logf("Sent Verification SMS, response was: %#v\n", messageResponse)
}

func TestSendSearch(t *testing.T) {
	// We need the request id, so we have to run this first.
	sendResponse := testSend(t)

	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &VerifySearchRequest{
		RequestID: sendResponse.RequestID,
	}

	messageResponse, err := nexmo.Verify.Search(message)
	if err != nil {
		t.Error("Failed to send verification search request with error:", err)
	}

	t.Logf("Sent Verification SMS, response was: %#v\n", messageResponse)
}
