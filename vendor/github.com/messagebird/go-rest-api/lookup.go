package messagebird

import "net/url"

// Formats represents phone number in multiple formats.
type Formats struct {
	E164          string
	International string
	National      string
	Rfc3966       string
}

// Lookup is used to validate and look up a mobile number.
type Lookup struct {
	Href          string
	CountryCode   string
	CountryPrefix int
	PhoneNumber   int64
	Type          string
	Formats       Formats
	HLR           *HLR
}

// LookupParams provide additional lookup information.
type LookupParams struct {
	CountryCode string
	Reference   string
}

type lookupRequest struct {
	CountryCode string `json:"countryCode,omitempty"`
	Reference   string `json:"reference,omitempty"`
}

func requestDataForLookup(params *LookupParams) *lookupRequest {
	request := &lookupRequest{}

	if params == nil {
		return request
	}

	request.CountryCode = params.CountryCode
	request.Reference = params.Reference

	return request
}

func paramsForLookup(params *LookupParams) *url.Values {
	urlParams := &url.Values{}

	if params == nil {
		return urlParams
	}

	if params.CountryCode != "" {
		urlParams.Set("countryCode", params.CountryCode)
	}
	if params.Reference != "" {
		urlParams.Set("reference", params.Reference)
	}

	return urlParams
}
