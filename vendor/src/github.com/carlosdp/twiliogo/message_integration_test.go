package twiliogo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegrationMessageList(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(API_KEY, API_TOKEN)

	messageList, err := GetMessageList(client)

	if assert.Nil(t, err, "Failed to retrieve message list") {
		messages := messageList.Messages
		assert.NotNil(t, messages, "Failed to retrieve messages")
	}
}

func TestIntegrationSendSMS(t *testing.T) {
	/* /Messages endpoint is currently not recognized by Test Credentials */
	/* CheckTestEnv(t) */

	/* client := NewClient(TEST_KEY, TEST_TOKEN) */

	/* message, err := NewMessage(client, TEST_FROM_NUMBER, TO_NUMBER, Body("Test Message")) */

	/* if assert.Nil(t, err, "Failed to Send SMS") { */
	/*   assert.Equal(t, message.Status, "queued", "Sending SMS failed, status: " + message.Status) */
	/* } */
	t.Skip()
}

func TestIntegrationMessageListNextPage(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(API_KEY, API_TOKEN)

	messageList, err := GetMessageList(client)

	if assert.Nil(t, err, "Failed to retrieve message list") {
		nextPageMessageList, err := messageList.NextPage()

		if assert.Nil(t, err, "Failed to retrieve message list") {
			assert.Equal(t, nextPageMessageList.Page, 1, "Page incorrect on next page")
		}
	}
}

func TestIntegrationGetMessage(t *testing.T) {
	CheckTestEnv(t)

	client := NewClient(API_KEY, API_TOKEN)

	messageList, err := GetMessageList(client)

	if assert.Nil(t, err, "Failed to retrieve message list") {
		messageSid := messageList.Messages[0].Sid
		message, err := GetMessage(client, messageSid)

		if assert.Nil(t, err, "Failed to retrieve message") {
			assert.Equal(t, message.Sid, messageSid, "Message was invalid")
		}
	}
}
