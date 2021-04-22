package voicemessage

import (
	"errors"
	"net/http"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
)

// VoiceMessage wraps data needed to transform text messages into voice messages.
// Voice messages are identified by a unique random ID. With this ID you can always check the status of the voice message through the provided endpoint.
type VoiceMessage struct {
	ID                string
	HRef              string
	Originator        string
	Body              string
	Reference         string
	Language          string
	Voice             string
	Repeat            int
	IfMachine         string
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        messagebird.Recipients
}

// VoiceMessageList represents a list of VoiceMessages.
type VoiceMessageList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      map[string]*string
	Items      []VoiceMessage
}

// Params struct provides additional VoiceMessage details.
type Params struct {
	Originator        string
	Reference         string
	Language          string
	Voice             string
	Repeat            int
	IfMachine         string
	ScheduledDatetime time.Time
}

type voiceMessageRequest struct {
	Recipients        []string `json:"recipients"`
	Body              string   `json:"body"`
	Originator        string   `json:"originator,omitempty"`
	Reference         string   `json:"reference,omitempty"`
	Language          string   `json:"language,omitempty"`
	Voice             string   `json:"voice,omitempty"`
	Repeat            int      `json:"repeat,omitempty"`
	IfMachine         string   `json:"ifMachine,omitempty"`
	ScheduledDatetime string   `json:"scheduledDatetime,omitempty"`
}

// path represents the path to the VoiceMessage resource.
const path = "voicemessages"

// Read retrieves the information of an existing VoiceMessage.
func Read(c *messagebird.Client, id string) (*VoiceMessage, error) {
	message := &VoiceMessage{}
	if err := c.Request(message, http.MethodGet, path+"/"+id, nil); err != nil {
		return nil, err
	}

	return message, nil
}

// List retrieves all VoiceMessages of the user.
func List(c *messagebird.Client) (*VoiceMessageList, error) {
	messageList := &VoiceMessageList{}
	if err := c.Request(messageList, http.MethodGet, path, nil); err != nil {
		return nil, err
	}

	return messageList, nil
}

// Create a new voice message for one or more recipients.
func Create(c *messagebird.Client, recipients []string, body string, params *Params) (*VoiceMessage, error) {
	requestData, err := requestDataForVoiceMessage(recipients, body, params)
	if err != nil {
		return nil, err
	}

	message := &VoiceMessage{}
	if err := c.Request(message, http.MethodPost, path, requestData); err != nil {
		return nil, err
	}

	return message, nil
}

func requestDataForVoiceMessage(recipients []string, body string, params *Params) (*voiceMessageRequest, error) {
	if len(recipients) == 0 {
		return nil, errors.New("at least 1 recipient is required")
	}
	if body == "" {
		return nil, errors.New("body is required")
	}

	request := &voiceMessageRequest{
		Recipients: recipients,
		Body:       body,
	}

	if params == nil {
		return request, nil
	}

	request.Originator = params.Originator
	request.Reference = params.Reference
	request.Language = params.Language
	request.Voice = params.Voice
	request.Repeat = params.Repeat
	request.IfMachine = params.IfMachine
	request.ScheduledDatetime = params.ScheduledDatetime.Format(time.RFC3339)

	return request, nil
}
