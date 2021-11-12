package sms77api

import "encoding/json"

type Carrier struct {
	Country     string `json:"country"`
	Name        string `json:"name"`
	NetworkCode string `json:"network_code"`
	NetworkType string `json:"network_type"`
}

type LookupParams struct {
	Type   string `json:"type"`
	Number string `json:"number,omitempty"`
	Json   bool   `json:"json,omitempty"`
}

type LookupCnamResponse struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Number  string `json:"number"`
	Success string `json:"success"`
}

type LookupFormatResponse struct {
	National               string `json:"national"`
	Carrier                string `json:"carrier"`
	CountryCode            string `json:"country_code"`
	CountryIso             string `json:"country_iso"`
	CountryName            string `json:"country_name"`
	International          string `json:"international"`
	InternationalFormatted string `json:"international_formatted"`
	NetworkType            string `json:"network_type"`
	Success                bool   `json:"success"`
}

type LookupHlrResponse struct {
	CountryCode               string  `json:"country_code"`
	CountryCodeIso3           *string `json:"country_code_iso3,omitempty"`
	CountryName               string  `json:"country_name"`
	CountryPrefix             string  `json:"country_prefix"`
	CurrentCarrier            Carrier `json:"current_carrier"`
	GsmCode                   string  `json:"gsm_code"`
	GsmMessage                string  `json:"gsm_message"`
	InternationalFormatNumber string  `json:"international_format_number"`
	InternationalFormatted    string  `json:"international_formatted"`
	LookupOutcome             bool    `json:"lookup_outcome"`
	LookupOutcomeMessage      string  `json:"lookup_outcome_message"`
	NationalFormatNumber      string  `json:"national_format_number"`
	OriginalCarrier           Carrier `json:"original_carrier"`
	Ported                    string  `json:"ported"`
	Reachable                 string  `json:"reachable"`
	Roaming                   string  `json:"roaming"`
	Status                    bool    `json:"status"`
	StatusMessage             string  `json:"status_message"`
	ValidNumber               string  `json:"valid_number"`
}

type LookupMnpResponse struct {
	Code    int64   `json:"code"`
	Mnp     Mnp     `json:"mnp"`
	Price   float64 `json:"price"`
	Success bool    `json:"success"`
}

type Mnp struct {
	Country                string `json:"country"`
	InternationalFormatted string `json:"international_formatted"`
	IsPorted               bool   `json:"isPorted"`
	Mccmnc                 string `json:"mccmnc"`
	NationalFormat         string `json:"national_format"`
	Network                string `json:"network"`
	Number                 string `json:"number"`
}

type LookupResource resource

func (api *LookupResource) Post(p LookupParams) (interface{}, error) {
	res, err := api.client.request("lookup", "GET", p)

	if err != nil {
		return nil, err
	}

	var js interface{}

	switch p.Type {
	case "mnp":
		if !p.Json {
			return res, nil
		}

		js = &LookupMnpResponse{}
	case "cnam":
		js = &LookupCnamResponse{}
	case "format":
		js = &LookupFormatResponse{}
	case "hlr":
		js = &LookupHlrResponse{}
	}

	if err := json.Unmarshal([]byte(res), js); err != nil {
		return nil, err
	}

	return js, nil
}
