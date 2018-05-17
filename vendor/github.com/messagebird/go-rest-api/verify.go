package messagebird

import (
	"net/url"
	"strconv"
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

func paramsForVerify(params *VerifyParams) *url.Values {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams
	}

	if params.Originator != "" {
		urlParams.Set("originator", params.Originator)
	}

	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}

	if params.Type != "" {
		urlParams.Set("type", params.Type)
	}

	if params.Template != "" {
		urlParams.Set("template", params.Template)
	}

	if params.DataCoding != "" {
		urlParams.Set("datacoding", params.DataCoding)
	}

	if params.Timeout != 0 {
		urlParams.Set("timeout", strconv.Itoa(params.Timeout))
	}

	if params.TokenLength != 0 {
		urlParams.Set("tokenLength", strconv.Itoa(params.TokenLength))
	}

	// Specific params for voice messages
	if params.Language != "" {
		urlParams.Set("language", params.Language)
	}
	if params.Voice != "" {
		urlParams.Set("voice", params.Voice)
	}

	return urlParams
}
