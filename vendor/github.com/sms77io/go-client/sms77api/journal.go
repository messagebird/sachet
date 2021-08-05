package sms77api

import (
	"encoding/json"
)

type JournalResource resource

type JournalType string

const (
	JournalTypeInbound  JournalType = "inbound"
	JournalTypeOutbound JournalType = "outbound"
	JournalTypeReplies  JournalType = "replies"
	JournalTypeVoice    JournalType = "voice"
)

type JournalBase struct {
	From      string `json:"from"`
	Id        string `json:"id"`
	Price     string `json:"price"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
	To        string `json:"to"`
}

type JournalOutbound struct {
	JournalBase
	Connection   string `json:"connection"`
	Dlr          string `json:"dlr"`
	DlrTimestamp string `json:"dlr_timestamp"`
	ForeignId    string `json:"foreign_id"`
	Label        string `json:"label"`
	Latency      string `json:"latency"`
	MccMnc       string `json:"mccmnc"`
	Type         string `json:"type"`
}

type JournalInbound struct {
	JournalBase
}

type JournalReplies struct {
	JournalBase
}

type JournalVoice struct {
	JournalBase
	Duration string `json:"duration"`
	Error    string `json:"error"`
	Status   string `json:"status"`
	Xml      bool   `json:"xml"`
}

type JournalParams struct {
	DateFrom string      `json:"date_from,omitempty"`
	DateTo   string      `json:"date_to,omitempty"`
	Id       int         `json:"id,omitempty"`
	State    string      `json:"state,omitempty"`
	To       string      `json:"to,omitempty"`
	Type     JournalType `json:"type"`
}

func BuildJournalParams(journalType JournalType, p *JournalParams) *JournalParams {
	p.Type = journalType

	return p
}

func (api *JournalResource) Inbound(p *JournalParams) ([]JournalInbound, error) {
	res, err := api.client.request("journal", "GET", BuildJournalParams(JournalTypeInbound, p))

	if err != nil {
		return nil, err
	}

	var js []JournalInbound

	if err := json.Unmarshal([]byte(res), &js); err != nil {
		return nil, err
	}

	return js, nil
}

func (api *JournalResource) Outbound(p *JournalParams) ([]JournalOutbound, error) {
	res, err := api.client.request("journal", "GET", BuildJournalParams(JournalTypeOutbound, p))

	if err != nil {
		return nil, err
	}

	var js []JournalOutbound

	if err := json.Unmarshal([]byte(res), &js); err != nil {
		return nil, err
	}

	return js, nil
}

func (api *JournalResource) Replies(p *JournalParams) ([]JournalReplies, error) {
	res, err := api.client.request("journal", "GET", BuildJournalParams(JournalTypeReplies, p))

	if err != nil {
		return nil, err
	}

	var js []JournalReplies

	if err := json.Unmarshal([]byte(res), &js); err != nil {
		return nil, err
	}

	return js, nil
}

func (api *JournalResource) Voice(p *JournalParams) ([]JournalVoice, error) {
	res, err := api.client.request("journal", "GET", BuildJournalParams(JournalTypeVoice, p))

	if err != nil {
		return nil, err
	}

	var js []JournalVoice

	if err := json.Unmarshal([]byte(res), &js); err != nil {
		return nil, err
	}

	return js, nil
}
