package sms77api

import (
	"errors"
	"strings"
)

type StatusResource resource

type StatusParams struct {
	MessageId uint64 `json:"msg_id"`
}

type Status struct {
	Code     string
	DateTime string
}

func makeStatus(res *string) (s *Status, e error) {
	if nil == res {
		return nil, errors.New("cannot make status from nil")
	}

	if StatusApiCodeInvalidMessageId == *res {
		return nil, errors.New(StatusApiCodeInvalidMessageId + ": Invalid message ID")
	}

	lines := strings.Split(*res, "\n")

	if 2 != len(lines) {
		return nil, errors.New("need exactly 2 lines to make a status")
	}

	return &Status{
		Code:     lines[0],
		DateTime: lines[1],
	}, nil
}

const StatusApiCodeInvalidMessageId = "901"

func (api *StatusResource) Text(p StatusParams) (*string, error) {
	res, err := api.client.request("status", "POST", p)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (api *StatusResource) Json(p StatusParams) (o *Status, e error) {
	r, e := api.Text(p)

	if nil != e {
		return
	}

	return makeStatus(r)
}
