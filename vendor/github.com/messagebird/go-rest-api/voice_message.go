package messagebird

import (
	"errors"
	"time"
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
	Recipients        Recipients
	Errors            []Error
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

// VoiceMessageParams struct provides additional VoiceMessage details.
type VoiceMessageParams struct {
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

func requestDataForVoiceMessage(recipients []string, body string, params *VoiceMessageParams) (*voiceMessageRequest, error) {
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
