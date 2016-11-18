package twiliogo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationCallList(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(API_KEY, API_TOKEN)

	callList, err := GetCallList(client)

	if assert.Nil(t, err, "Failed to retrieve call list") {
		calls := callList.GetCalls()

		assert.NotNil(t, calls, "Failed to retrieve calls")
	}
}

func TestIntegrationMakingCall(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(TEST_KEY, TEST_TOKEN)

	call, err := NewCall(client, TEST_FROM_NUMBER, TO_NUMBER, Callback("http://test.com"))

	if assert.Nil(t, err, "Failed to make call") {
		assert.Equal(t, call.Status, "queued", "Making Call failed, status: "+call.Status)
	}
}

func TestIntegrationCallListNextPage(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(API_KEY, API_TOKEN)

	callList, err := GetCallList(client)

	if assert.Nil(t, err, "Failed to retrieve call list") {
		nextPageCallList, err := callList.NextPage()

		if assert.Nil(t, err, "Failed to retrieve next page") {
			assert.Equal(t, nextPageCallList.Page, 1, "Page incorrect on next page")
		}
	}
}

func TestIntegrationGetCall(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(API_KEY, API_TOKEN)

	callList, err := GetCallList(client)

	if assert.Nil(t, err, "Failed to retrieve call list") {
		callSid := callList.Calls[0].Sid
		call, err := GetCall(client, callSid)

		if assert.Nil(t, err, "Failed to retrieve call") {
			assert.Equal(t, call.Sid, callSid, "Call was invalid")
		}
	}
}
