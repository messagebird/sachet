package sms77api

import (
	"strconv"
	"strings"
)

type Voice struct {
	Code int
	Cost float64
	Id   int
}

type VoiceParams struct {
	Debug bool   `json:"debug,omitempty"`
	To    string `json:"to"`
	Text  string `json:"text"`
	Xml   bool   `json:"xml,omitempty"`
	From  string `json:"from,omitempty"`
}

type VoiceResource resource

func makeVoice(res string) Voice {
	lines := strings.Split(res, "\n")

	code, _ := strconv.Atoi(lines[0])
	id, _ := strconv.Atoi(lines[1])
	cost, _ := strconv.ParseFloat(lines[2], 64)

	return Voice{
		Code: code,
		Cost: cost,
		Id:   id,
	}
}

func (api *VoiceResource) Text(p VoiceParams) (*string, error) {
	res, err := api.client.request("voice", "POST", p)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (api *VoiceResource) Json(p VoiceParams) (o Voice, e error) {
	r, e := api.Text(p)

	if nil != e {
		return
	}

	return makeVoice(*r), nil
}
