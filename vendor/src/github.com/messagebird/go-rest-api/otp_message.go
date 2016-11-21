package messagebird

import (
	"net/url"
	"time"
)

type OtpMessage struct {
	Id                 string
	Recipient          string
	Reference          string
	Status             string
	Href               map[string]string `json:"href"`
	CreatedDatetime    *time.Time
	ValidUntilDatetime *time.Time
	Errors             []Error
}

type OtpParams struct {
	Reference  string
	Originator string
	Type       string
	Language   string
	Voice      string
	Template   string
	DataCoding string
}

// paramsForOtp converts the specified OtpParams struct to a
// url.Values pointer and returns it.
func paramsForOtp(params *OtpParams) *url.Values {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams
	}

	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}
	if params.Originator != "" {
		urlParams.Set("originator", params.Originator)
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

	// Specific params for voice messages
	if params.Language != "" {
		urlParams.Set("language", params.Language)
	}
	if params.Voice != "" {
		urlParams.Set("voice", params.Voice)
	}

	return urlParams
}
