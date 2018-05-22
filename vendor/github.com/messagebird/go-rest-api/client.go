//
// Copyright (c) 2014 MessageBird B.V.
// All rights reserved.
//
// Author: Maurice Nonnekes <maurice@messagebird.com>

// Package messagebird is an official library for interacting with MessageBird.com API.
// The MessageBird API connects your website or application to operators around the world. With our API you can integrate SMS, Chat & Voice.
// More documentation you can find on the MessageBird developers portal: https://developers.messagebird.com/
package messagebird

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strings"
)

const (
	// ClientVersion is used in User-Agent request header to provide server with API level.
	ClientVersion = "4.2.0"

	// Endpoint points you to MessageBird REST API.
	Endpoint = "https://rest.messagebird.com"
)

const (
	// HLRPath represents the path to the HLR resource.
	HLRPath = "hlr"
	// MessagePath represents the path to the Message resource.
	MessagePath = "messages"
	// MMSPath represents the path to the MMS resource.
	MMSPath = "mms"
	// VoiceMessagePath represents the path to the VoiceMessage resource.
	VoiceMessagePath = "voicemessages"
	// VerifyPath represents the path to the Verify resource.
	VerifyPath = "verify"
	// LookupPath represents the path to the Lookup resource.
	LookupPath = "lookup"
)

var (
	// ErrResponse is returned when we were able to contact API but request was not successful and contains error details.
	ErrResponse = errors.New("The MessageBird API returned an error")

	// ErrUnexpectedResponse is used when there was an internal server error and nothing can be done at this point.
	ErrUnexpectedResponse = errors.New("The MessageBird API is currently unavailable")
)

// Client is used to access API with a given key.
// Uses standard lib HTTP client internally, so should be reused instead of created as needed and it is safe for concurrent use.
type Client struct {
	AccessKey  string       // The API access key
	HTTPClient *http.Client // The HTTP client to send requests on
	DebugLog   *log.Logger  // Optional logger for debugging purposes
}

// New creates a new MessageBird client object.
func New(AccessKey string) *Client {
	return &Client{AccessKey: AccessKey, HTTPClient: &http.Client{}}
}

// Request is for internal use only and unstable.
func (c *Client) Request(v interface{}, method, path string, data interface{}) error {
	if !strings.HasPrefix(path, "https://") && !strings.HasPrefix(path, "http://") {
		path = fmt.Sprintf("%s/%s", Endpoint, path)
	}
	uri, err := url.Parse(path)
	if err != nil {
		return err
	}

	var jsonEncoded []byte
	if data != nil {
		jsonEncoded, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

	request, err := http.NewRequest(method, uri.String(), bytes.NewBuffer(jsonEncoded))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "AccessKey "+c.AccessKey)
	request.Header.Set("User-Agent", "MessageBird/ApiClient/"+ClientVersion+" Go/"+runtime.Version())

	if c.DebugLog != nil {
		if data != nil {
			c.DebugLog.Printf("HTTP REQUEST: %s %s %s", method, uri.String(), jsonEncoded)
		} else {
			c.DebugLog.Printf("HTTP REQUEST: %s %s", method, uri.String())
		}
	}

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
		c.DebugLog.Printf("HTTP RESPONSE: %s", string(responseBody))
	}

	// Status code 500 is a server error and means nothing can be done at this
	// point.
	if response.StatusCode == 500 {
		return ErrUnexpectedResponse
	}
	// Status codes 200 and 201 are indicative of being able to convert the
	// response body to the struct that was specified.
	if response.StatusCode == 200 || response.StatusCode == 201 {
		if err := json.Unmarshal(responseBody, &v); err != nil {
			return fmt.Errorf("could not decode response JSON, %s: %v", string(responseBody), err)
		}
		return nil
	}

	// We're dealing with an API error here. try to decode it, but don't do
	// anything with the error. This is because not all values of `v` have
	// `Error` properties and decoding could fail.
	json.Unmarshal(responseBody, &v)

	// Anything else than a 200/201/500 should be a JSON error.
	return ErrResponse
}

// Balance returns the balance information for the account that is associated
// with the access key.
func (c *Client) Balance() (*Balance, error) {
	balance := &Balance{}
	if err := c.Request(balance, http.MethodGet, "balance", nil); err != nil {
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
	if err := c.Request(hlr, http.MethodGet, HLRPath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}

// HLRs lists all HLR objects that were previously created by the NewHLR
// function.
func (c *Client) HLRs() (*HLRList, error) {
	hlrList := &HLRList{}
	if err := c.Request(hlrList, http.MethodGet, HLRPath, nil); err != nil {
		if err == ErrResponse {
			return hlrList, err
		}

		return nil, err
	}

	return hlrList, nil
}

// NewHLR retrieves the information of an existing HLR.
func (c *Client) NewHLR(msisdn string, reference string) (*HLR, error) {
	requestData, err := requestDataForHLR(msisdn, reference)
	if err != nil {
		return nil, err
	}

	hlr := &HLR{}

	if err := c.Request(hlr, http.MethodPost, HLRPath, requestData); err != nil {
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
	if err := c.Request(message, http.MethodGet, MessagePath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// Messages retrieves all messages of the user represented as a MessageList object.
func (c *Client) Messages(msgListParams *MessageListParams) (*MessageList, error) {
	messageList := &MessageList{}
	params, err := paramsForMessageList(msgListParams)
	if err != nil {
		return messageList, err
	}

	if err := c.Request(messageList, http.MethodGet, MessagePath+"?"+params.Encode(), nil); err != nil {
		if err == ErrResponse {
			return messageList, err
		}

		return nil, err
	}

	return messageList, nil
}

// NewMessage creates a new message for one or more recipients.
func (c *Client) NewMessage(originator string, recipients []string, body string, msgParams *MessageParams) (*Message, error) {
	requestData, err := requestDataForMessage(originator, recipients, body, msgParams)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	if err := c.Request(message, http.MethodPost, MessagePath, requestData); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// MMSMessage retrieves the information of an existing MmsMessage.
func (c *Client) MMSMessage(id string) (*MMSMessage, error) {
	mmsMessage := &MMSMessage{}
	if err := c.Request(mmsMessage, http.MethodGet, MMSPath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return mmsMessage, err
		}

		return nil, err
	}

	return mmsMessage, nil
}

// NewMMSMessage creates a new MMS message for one or more recipients.
func (c *Client) NewMMSMessage(originator string, recipients []string, msgParams *MMSMessageParams) (*MMSMessage, error) {
	params, err := paramsForMMSMessage(msgParams)
	if err != nil {
		return nil, err
	}

	params.Set("originator", originator)
	params.Set("recipients", strings.Join(recipients, ","))

	mmsMessage := &MMSMessage{}
	if err := c.Request(mmsMessage, http.MethodPost, MMSPath, params); err != nil {
		if err == ErrResponse {
			return mmsMessage, err
		}

		return nil, err
	}

	return mmsMessage, nil
}

// VoiceMessage retrieves the information of an existing VoiceMessage.
func (c *Client) VoiceMessage(id string) (*VoiceMessage, error) {
	message := &VoiceMessage{}
	if err := c.Request(message, http.MethodGet, VoiceMessagePath+"/"+id, nil); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// VoiceMessages retrieves all VoiceMessages of the user.
func (c *Client) VoiceMessages() (*VoiceMessageList, error) {
	messageList := &VoiceMessageList{}
	if err := c.Request(messageList, http.MethodGet, VoiceMessagePath, nil); err != nil {
		if err == ErrResponse {
			return messageList, err
		}

		return nil, err
	}

	return messageList, nil
}

// NewVoiceMessage creates a new voice message for one or more recipients.
func (c *Client) NewVoiceMessage(recipients []string, body string, params *VoiceMessageParams) (*VoiceMessage, error) {
	requestData, err := requestDataForVoiceMessage(recipients, body, params)
	if err != nil {
		return nil, err
	}

	message := &VoiceMessage{}
	if err := c.Request(message, http.MethodPost, VoiceMessagePath, requestData); err != nil {
		if err == ErrResponse {
			return message, err
		}

		return nil, err
	}

	return message, nil
}

// NewVerify generates a new One-Time-Password for one recipient.
func (c *Client) NewVerify(recipient string, params *VerifyParams) (*Verify, error) {
	requestData, err := requestDataForVerify(recipient, params)
	if err != nil {
		return nil, err
	}

	verify := &Verify{}
	if err := c.Request(verify, http.MethodPost, VerifyPath, requestData); err != nil {
		if err == ErrResponse {
			return verify, err
		}

		return nil, err
	}

	return verify, nil
}

// VerifyToken performs token value check against MessageBird API.
func (c *Client) VerifyToken(id, token string) (*Verify, error) {
	params := &url.Values{}
	params.Set("token", token)

	path := VerifyPath + "/" + id + "?" + params.Encode()

	verify := &Verify{}
	if err := c.Request(verify, http.MethodGet, path, nil); err != nil {
		if err == ErrResponse {
			return verify, err
		}

		return nil, err
	}

	return verify, nil
}

// Lookup performs a new lookup for the specified number.
func (c *Client) Lookup(phoneNumber string, params *LookupParams) (*Lookup, error) {
	urlParams := paramsForLookup(params)
	path := LookupPath + "/" + phoneNumber + "?" + urlParams.Encode()

	lookup := &Lookup{}
	if err := c.Request(lookup, http.MethodPost, path, nil); err != nil {
		if err == ErrResponse {
			return lookup, err
		}

		return nil, err
	}

	return lookup, nil
}

// NewLookupHLR creates a new HLR lookup for the specified number.
func (c *Client) NewLookupHLR(phoneNumber string, params *LookupParams) (*HLR, error) {
	requestData := requestDataForLookup(params)
	path := LookupPath + "/" + phoneNumber + "/" + HLRPath

	hlr := &HLR{}
	if err := c.Request(hlr, http.MethodPost, path, requestData); err != nil {
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
	path := LookupPath + "/" + phoneNumber + "/" + HLRPath + "?" + urlParams.Encode()

	hlr := &HLR{}
	if err := c.Request(hlr, http.MethodGet, path, nil); err != nil {
		if err == ErrResponse {
			return hlr, err
		}

		return nil, err
	}

	return hlr, nil
}
