package messagebird

import (
	"errors"
	"time"
)

// Verify object represents MessageBird server response.
type Verify struct {
	ID                 string
	HRef               string
	Reference          string
	Status             string
	Messages           map[string]string
	CreatedDatetime    *time.Time
	ValidUntilDatetime *time.Time
	Recipient          int
	Errors             []Error
}

// VerifyParams handles optional verification parameters.
type VerifyParams struct {
	Originator  string
	Reference   string
	Type        string
	Template    string
	DataCoding  string
	Voice       string
	Language    string
	Timeout     int
	TokenLength int
}

type verifyRequest struct {
	Recipient   string `json:"recipient"`
	Originator  string `json:"originator,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Type        string `json:"type,omitempty"`
	Template    string `json:"template,omitempty"`
	DataCoding  string `json:"dataCoding,omitempty"`
	Voice       string `json:"voice,omitempty"`
	Language    string `json:"language,omitempty"`
	Timeout     int    `json:"timeout,omitempty"`
	TokenLength int    `json:"tokenLength,omitempty"`
}

func requestDataForVerify(recipient string, params *VerifyParams) (*verifyRequest, error) {
	if recipient == "" {
		return nil, errors.New("recipient is required")
	}

	request := &verifyRequest{
		Recipient: recipient,
	}

	if params == nil {
		return request, nil
	}

	request.Originator = params.Originator
	request.Reference = params.Reference
	request.Type = params.Type
	request.Template = params.Template
	request.DataCoding = params.DataCoding
	request.Voice = params.Voice
	request.Language = params.Language
	request.Timeout = params.Timeout
	request.TokenLength = params.TokenLength

	return request, nil
}
