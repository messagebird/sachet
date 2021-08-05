package sms77api

import "encoding/json"

type HookEventType string

const (
	HookEventTypeSmsStatus   HookEventType = "dlr"
	HookEventTypeVoiceStatus HookEventType = "voice_status"
	HookEventTypeInboundSms  HookEventType = "sms_mo"
)

type HookRequestMethod string

const (
	HookRequestMethodGet  HookRequestMethod = "GET"
	HookRequestMethodJson HookRequestMethod = "JSON"
	HookRequestMethodPost HookRequestMethod = "POST"
)

type HooksAction string

const (
	HooksActionRead        HooksAction = "read"
	HooksActionSubscribe   HooksAction = "subscribe"
	HooksActionUnsubscribe HooksAction = "unsubscribe"
)

type Hook struct {
	Created       string            `json:"created"`
	EventType     HookEventType     `json:"event_type"`
	Id            string            `json:"id"`
	RequestMethod HookRequestMethod `json:"request_method"`
	TargetUrl     string            `json:"target_url"`
}

type HooksParams struct {
	Action        HooksAction       `json:"action"`
	EventType     HookEventType     `json:"event_type,omitempty"`
	Id            int               `json:"id,omitempty"`
	RequestMethod HookRequestMethod `json:"request_method,omitempty"`
	TargetUrl     string            `json:"target_url,omitempty"`
}

type HooksReadResponse struct {
	Success bool   `json:"success"`
	Hooks   []Hook `json:"hooks"`
}

type HooksUnsubscribeResponse struct {
	Success bool `json:"success"`
}

type HooksSubscribeResponse struct {
	Id      int  `json:"id"`
	Success bool `json:"success"`
}

type HooksResource resource

func (api *HooksResource) Request(p HooksParams) (interface{}, error) {
	method := "POST"
	if p.Action == HooksActionRead {
		method = "GET"
	}

	res, err := api.client.request("hooks", method, p)

	if err != nil {
		return nil, err
	}

	var js interface{}

	switch p.Action {
	case HooksActionRead:
		js = &HooksReadResponse{}
	case HooksActionSubscribe:
		js = &HooksSubscribeResponse{}
	case HooksActionUnsubscribe:
		js = &HooksUnsubscribeResponse{}
	}

	if err := json.Unmarshal([]byte(res), js); err != nil {
		return nil, err
	}

	return js, nil
}
