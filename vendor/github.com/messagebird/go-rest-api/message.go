package messagebird

import (
	"errors"
	"net/url"
	"strconv"
	"time"
)

// TypeDetails is a hash with extra information.
// Is only used when a binary or premium message is sent.
type TypeDetails map[string]interface{}

// Message struct represents a message at MessageBird.com
type Message struct {
	ID                string
	HRef              string
	Direction         string
	Type              string
	Originator        string
	Body              string
	Reference         string
	Validity          *int
	Gateway           int
	TypeDetails       TypeDetails
	DataCoding        string
	MClass            int
	ScheduledDatetime *time.Time
	CreatedDatetime   *time.Time
	Recipients        Recipients
	Errors            []Error
}

// MessageList represents a list of Messages.
type MessageList struct {
	Offset     int
	Limit      int
	Count      int
	TotalCount int
	Links      map[string]*string
	Items      []Message
}

// MessageParams provide additional message send options and used in URL as params.
type MessageParams struct {
	Type              string
	Reference         string
	Validity          int
	Gateway           int
	TypeDetails       TypeDetails
	DataCoding        string
	ScheduledDatetime time.Time
}

// MessageListParams provides additional message list options.
type MessageListParams struct {
	Originator string
	Direction  string
	Type       string
	Limit      int
	Offset     int
}

type messageRequest struct {
	Originator        string      `json:"originator"`
	Body              string      `json:"body"`
	Recipients        []string    `json:"recipients"`
	Type              string      `json:"type,omitempty"`
	Reference         string      `json:"reference,omitempty"`
	Validity          int         `json:"validity,omitempty"`
	Gateway           int         `json:"gateway,omitempty"`
	TypeDetails       TypeDetails `json:"typeDetails,omitempty"`
	DataCoding        string      `json:"datacoding,omitempty"`
	MClass            int         `json:"mclass,omitempty"`
	ScheduledDatetime string      `json:"scheduledDatetime,omitempty"`
}

func requestDataForMessage(originator string, recipients []string, body string, params *MessageParams) (*messageRequest, error) {
	if originator == "" {
		return nil, errors.New("originator is required")
	}
	if len(recipients) == 0 {
		return nil, errors.New("at least 1 recipient is required")
	}
	if body == "" {
		return nil, errors.New("body is required")
	}

	request := &messageRequest{
		Originator: originator,
		Recipients: recipients,
		Body:       body,
	}

	if params == nil {
		return request, nil
	}

	request.Type = params.Type
	if request.Type == "flash" {
		request.MClass = 0
	} else {
		request.MClass = 1
	}

	if !params.ScheduledDatetime.IsZero() {
		request.ScheduledDatetime = params.ScheduledDatetime.Format(time.RFC3339)
	}

	request.Reference = params.Reference
	request.Validity = params.Validity
	request.Gateway = params.Gateway
	request.TypeDetails = params.TypeDetails
	request.DataCoding = params.DataCoding

	return request, nil
}

// paramsForMessageList converts the specified MessageListParams struct to a
// url.Values pointer and returns it.
func paramsForMessageList(params *MessageListParams) (*url.Values, error) {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams, nil
	}

	if params.Direction != "" {
		urlParams.Set("direction", params.Direction)
	}
	if params.Originator != "" {
		urlParams.Set("originator", params.Originator)
	}
	if params.Limit != 0 {
		urlParams.Set("limit", strconv.Itoa(params.Limit))
	}
	urlParams.Set("offset", strconv.Itoa(params.Offset))

	return urlParams, nil
}
