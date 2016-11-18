//
// Copyright (c) 2014 MessageBird B.V.
// All rights reserved.
//
// Author: Maurice Nonnekes <maurice@messagebird.com>

package messagebird

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const (
	ClientVersion = "2.3.0"
	Endpoint      = "https://rest.messagebird.com"
)

var (
	ErrResponse           = errors.New("The MessageBird API returned an error")
	ErrUnexpectedResponse = errors.New("The MessageBird API is currently unavailable")
)

type Client struct {
	AccessKey  string       // The API access key
	HTTPClient *http.Client // The HTTP client to send requests on
	DebugLog   *log.Logger  // Optional logger for debugging purposes
}

// New creates a new MessageBird client object.
func New(AccessKey string) *Client {
	return &Client{AccessKey: AccessKey, HTTPClient: &http.Client{}}
}

func (c *Client) request(v interface{}, path string, params *url.Values) error {
	uri, err := url.Parse(Endpoint + "/" + path)
	if err != nil {
		return err
	}

	var request *http.Request
	if params != nil {
		body := params.Encode()
		if request, err = http.NewRequest("POST", uri.String(), strings.NewReader(body)); err != nil {
			return err
		}

		if c.DebugLog != nil {
			if unescapedBody, err := url.QueryUnescape(body); err == nil {
				log.Printf("HTTP REQUEST: POST %s %s", uri.String(), unescapedBody)
			} else {
				log.Printf("HTTP REQUEST: POST %s %s", uri.String(), body)
			}
		}

		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		if request, err = http.NewRequest("GET", uri.String(), nil); err != nil {
			return err
		}

		if c.DebugLog != nil {
			log.Printf("HTTP REQUEST: GET %s", uri.String())
		}
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", "AccessKey "+c.AccessKey)
	request.Header.Add("User-Agent", "MessageBird/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if c.DebugLog != nil {
		log.Printf("HTTP RESPONSE: %s", string(responseBody))
	}

	// Status code 500 is a server error and means nothing can be done at this
	// point.
	if response.StatusCode == 500 {
		return ErrUnexpectedResponse
	}

	if err = json.Unmarshal(responseBody, &v); err != nil {
		return err
	}

	// Status codes 200 and 201 are indicative of being able to convert the
	// response body to the struct that was specified.
	if response.StatusCode == 200 || response.StatusCode == 201 {
		return nil
	}

	// Anything else than a 200/201/500 should be a JSON error.
	return ErrResponse
}

// Balance returns the balance information for the account that is associated
// with the access key.
func (c *Client) Balance() (*Balance, error) {
	balance := &Balance{}
	if err := c.request(balance, "balance", nil); err != nil {
		if err == ErrResponse {
			return balance, err
		}

		return nil, err
	}

	return balance, nil
}

// HLR looks up an existing HLR object for the specified id that was previously
// created by the NewHLR function.
func (c *Client) HLR(id string) (*HLR, error) {
	hlr := &HLR{}
	if err := c.request(hlr, "hlr/"+id, nil); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// NewHLR retrieves the information of an existing HLR.
func (c *Client) NewHLR(msisdn, reference string) (*HLR, error) {
	params := &url.Values{
		"msisdn":    {msisdn},
		"reference": {reference}}

	hlr := &HLR{}
	if err := c.request(hlr, "hlr", params); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// Message retrieves the information of an existing Message.
func (c *Client) Message(id string) (*Message, error) {
	message := &Message{}
	if err := c.request(message, "messages/"+id, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// NewMessage creates a new message for one or more recipients.
func (c *Client) NewMessage(originator string, recipients []string, body string, msgParams *MessageParams) (*Message, error) {
	params, err := paramsForMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("body", body)
	params.Set("recipients", strings.Join(recipients, ","))

	message := &Message{}
	if err := c.request(message, "messages", params); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// VoiceMessage retrieves the information of an existing VoiceMessage.
func (c *Client) VoiceMessage(id string) (*VoiceMessage, error) {
	message := &VoiceMessage{}
	if err := c.request(message, "voicemessages/"+id, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// NewVoiceMessage creates a new voice message for one or more recipients.
func (c *Client) NewVoiceMessage(recipients []string, body string, params *VoiceMessageParams) (*VoiceMessage, error) {
	urlParams := paramsForVoiceMessage(params)
	urlParams.Set("body", body)
	urlParams.Set("recipients", strings.Join(recipients, ","))

	message := &VoiceMessage{}
	if err := c.request(message, "voicemessages", urlParams); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// OtpGenerate generates a new One-Time-Password for one recipient.
func (c *Client) OtpGenerate(recipient string, params *OtpParams) (*OtpMessage, error) {
	urlParams := paramsForOtp(params)
	urlParams.Set("recipient", recipient)

	message := &OtpMessage{}
	if err := c.request(message, "otp/generate", urlParams); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// OtpVerify verifies the token that was generated with OtpGenerate.
func (c *Client) OtpVerify(recipient string, token string, params *OtpParams) (*OtpMessage, error) {
	urlParams := paramsForOtp(params)
	urlParams.Set("recipient", recipient)
	urlParams.Set("token", token)

	path := "otp/verify?" + urlParams.Encode()

	message := &OtpMessage{}
	if err := c.request(message, path, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// Lookup performs a new lookup for the specified number.
func (c *Client) Lookup(phoneNumber string, params *LookupParams) (*Lookup, error) {
	urlParams := paramsForLookup(params)
	path := "lookup/" + phoneNumber + "?" + urlParams.Encode()

	lookup := &Lookup{}
	if err := c.request(lookup, path, nil); err != nil {
		if err == ErrResponse {
			return lookup, err
		}

		return nil, err
	}

	return lookup, nil
}

// NewLookupHLR creates a new HLR lookup for the specified number.
func (c *Client) NewLookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	urlParams := paramsForLookup(params)
	path := "lookup/" + phoneNumber + "/hlr"

	hlr := &HLR{}
	if err := c.request(hlr, path, urlParams); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// LookupHLR performs a HLR lookup for the specified number.
func (c *Client) LookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	urlParams := paramsForLookup(params)
	path := "lookup/" + phoneNumber + "/hlr?" + urlParams.Encode()

	hlr := &HLR{}
	if err := c.request(hlr, path, nil); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}
