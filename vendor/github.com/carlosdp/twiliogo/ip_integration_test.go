package twiliogo

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationIPMessaging(t *testing.T) {
	CheckTestEnv(t)

	client := NewIPMessagingClient(API_KEY, API_TOKEN)

	service, err := NewIPService(client, "integration_test", "", "", 60*time.Second, nil)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to create service: %v", err)) {
		return
	}
	ssid := service.Sid
	if !assert.NotEqual(t, "", ssid, "Service SID was empty") {
		return
	}

	defer DeleteIPService(client, ssid)
	serviceList, err := ListIPServices(client)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to retrieve service list: %v", err)) {
		return
	}
	services, err := serviceList.GetAllServices()
	if !assert.Nil(t, err, fmt.Sprintf("Failed to get all services: %v", err)) {
		return
	}

	found := false
	for _, s := range services {
		if s.FriendlyName == "integration_test" && s.Sid == ssid {
			found = true
			break
		}
	}
	if !assert.True(t, found, "Could not find service") {
		return
	}

	channel, err := NewIPChannel(client, ssid, "integration-channel", "", false, "")
	if !assert.Nil(t, err, fmt.Sprintf("Failed to create channel: %v", err)) {
		return
	}
	csid := channel.Sid
	if !assert.NotEqual(t, "", csid, "Channel SID was empty") {
		return
	}

	channelList, err := ListIPChannels(client, ssid)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to list channels: %v", err)) {
		return
	}
	channels, err := channelList.GetAllChannels()
	if !assert.Nil(t, err, fmt.Sprintf("Failed to get all channels: %v", err)) {
		return
	}

	found = false
	for _, c := range channels {
		if c.Sid == csid {
			found = true
			break
		}
	}
	if !assert.True(t, found, "Could not find channel") {
		return
	}

	// disables for now because I don't have keys to test with
	/*
		credential, err := NewIPCredential(client, "integration_cred", "gcm", false, "", "", "")
		if !assert.Nil(t, err, fmt.Sprintf("Failed to create credential: %v", err)) {
			return
		}
		credSid := credential.Sid
		if !assert.NotEqual(t, "", credSid, "Credential SID was empty") {
			return
		}
	*/

	role, err := NewIPRole(client, ssid, "integration_role", "channel",
		[]string{PermissionSendMessage, PermissionEditOwnMessage})
	if !assert.Nil(t, err, fmt.Sprintf("Failed to create role: %v", err)) {
		return
	}

	roleSid := role.Sid
	if !assert.NotEqual(t, "", roleSid, "Role SID was empty") {
		return
	}

	user, err := NewIPUser(client, ssid, "integration_user", "")
	if !assert.Nil(t, err, fmt.Sprintf("Failed to create user: %v", err)) {
		return
	}

	userSid := user.Sid
	if !assert.NotEqual(t, "", userSid, "User SID was empty") {
		return
	}

	member, err := AddIPMemberToChannel(client, ssid, csid, "integration_user", roleSid)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to add member: %v", err)) {
		return
	}

	memberSid := member.Sid
	if !assert.NotEqual(t, "", memberSid, "Member SID was empty") {
		return
	}

	_, err = NewIPUser(client, ssid, "integration_user2", "")
	if !assert.Nil(t, err, fmt.Sprintf("Failed to create user: %v", err)) {
		return
	}
	_, err = AddIPMemberToChannel(client, ssid, csid, "integration_user2", roleSid)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to add member: %v", err)) {
		return
	}

	memberList, err := ListIPMembers(client, ssid, csid)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to add member: %v", err)) {
		return
	}

	members, err := memberList.GetAllMembers()
	if !assert.Nil(t, err, fmt.Sprintf("Failed to list members: %v", err)) {
		return
	}

	found = false
	for _, m := range members {
		if m.Sid == memberSid && m.Identity == member.Identity {
			found = true
			break
		}
	}

	if !assert.True(t, found, "Could not find member") {
		return
	}

	message, err := SendIPMessageToChannel(client, ssid, csid, "integration_user2", "testing integration")
	if !assert.Nil(t, err, fmt.Sprintf("Failed to send message: %v", err)) {
		return
	}
	if !assert.NotEqual(t, "", message.Sid, "Message SID was empty") {
		return
	}

	found = false
	messageList, err := ListIPMessages(client, ssid, csid)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to list messages: %v", err)) {
		return
	}
	messages, err := messageList.GetAllMessages()
	if !assert.Nil(t, err, fmt.Sprintf("Failed to get all messages: %v", err)) {
		return
	}
	found = false
	for _, m := range messages {
		if m.Sid == message.Sid && m.Body == "testing integration" {
			found = true
			break
		}
	}
	if !assert.True(t, found, "Could not find message") {
		return
	}
}
