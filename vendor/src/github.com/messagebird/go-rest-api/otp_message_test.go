package messagebird

import (
	"testing"
	"time"
)

var OtpGenerateObject []byte = []byte(`{
    "id": "76d429606554b5827a46b12o29500954",
    "recipient": "31630174123",
    "reference": null,
    "status": "sent",
    "href": {
        "message": "https://rest.messagebird.com/messages/5b823a806554b5827b0cb66b58154766"
    },
    "createdDatetime": "2015-05-07T12:18:47+00:00",
    "validUntilDatetime": "2015-05-07T12:19:17+00:00"
}`)

func TestOtpGenerate(t *testing.T) {
	SetServerResponse(200, OtpGenerateObject)

	result, err := mbClient.OtpGenerate("31630174123", nil)
	if err != nil {
		t.Fatalf("Didn't expect error while generating an OTP: %s", err)
	}

	if result.Id != "76d429606554b5827a46b12o29500954" {
		t.Errorf("Unexpected OTP message id: %s", result.Id)
	}

	if result.Recipient != "31630174123" {
		t.Errorf("Unexpected OTP message recipient: %s", result.Recipient)
	}

	if result.Reference != "" {
		t.Errorf("Unexpected OTP message reference: %s", result.Reference)
	}

	if result.Status != "sent" {
		t.Errorf("Unexpected OTP message status: %s", result.Status)
	}

	if result.CreatedDatetime == nil || result.CreatedDatetime.Format(time.RFC3339) != "2015-05-07T12:18:47Z" {
		t.Errorf("Unexpected OTP message created datetime: %s", result.CreatedDatetime)
	}

	if result.ValidUntilDatetime == nil || result.ValidUntilDatetime.Format(time.RFC3339) != "2015-05-07T12:19:17Z" {
		t.Errorf("Unexpected OTP message valid until datetime: %s", result.ValidUntilDatetime)
	}

	message, exists := result.Href["message"]
	if !exists {
		t.Errorf("Unexpected OTP message href value: %s", result.Href)
	}

	if exists && message != "https://rest.messagebird.com/messages/5b823a806554b5827b0cb66b58154766" {
		t.Errorf("Unexpected OTP message href message: %s", message)
	}
}

var OtpVerifyObject []byte = []byte(`{
    "id": "8b912ea03554b7a3d5d6e22o95082672",
    "recipient": "31630174123",
    "reference": null,
    "status": "verified",
    "href": {
        "message": "https://rest.messagebird.com/messages/8668de203554b7a3d62b5e7b62054920"
    },
    "createdDatetime": "2015-05-07T14:44:13+00:00",
    "validUntilDatetime": "2015-05-07T14:44:43+00:00"
}`)

func TestOtpVerify(t *testing.T) {
	SetServerResponse(200, OtpVerifyObject)

	result, err := mbClient.OtpVerify("31630174123", "302443", nil)
	if err != nil {
		t.Fatalf("Didn't expect error while verifying an OTP: %s", err)
	}

	if result.Id != "8b912ea03554b7a3d5d6e22o95082672" {
		t.Errorf("Unexpected OTP message id: %s", result.Id)
	}

	if result.Recipient != "31630174123" {
		t.Errorf("Unexpected OTP message recipient: %s", result.Recipient)
	}

	if result.Reference != "" {
		t.Errorf("Unexpected OTP message reference: %s", result.Reference)
	}

	if result.Status != "verified" {
		t.Errorf("Unexpected OTP message status: %s", result.Status)
	}

	if result.CreatedDatetime == nil || result.CreatedDatetime.Format(time.RFC3339) != "2015-05-07T14:44:13Z" {
		t.Errorf("Unexpected OTP message created datetime: %s", result.CreatedDatetime)
	}

	if result.ValidUntilDatetime == nil || result.ValidUntilDatetime.Format(time.RFC3339) != "2015-05-07T14:44:43Z" {
		t.Errorf("Unexpected OTP message valid until datetime: %s", result.ValidUntilDatetime)
	}

	message, exists := result.Href["message"]
	if !exists {
		t.Errorf("Unexpected OTP message href value: %s", result.Href)
	}

	if exists && message != "https://rest.messagebird.com/messages/8668de203554b7a3d62b5e7b62054920" {
		t.Errorf("Unexpected OTP message href message: %s", message)
	}
}
