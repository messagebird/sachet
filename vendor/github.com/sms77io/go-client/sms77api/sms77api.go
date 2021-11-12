package sms77api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HttpMethod string
type Options struct {
	ApiKey   string
	Debug    bool
	SentWith string
}
type resource struct {
	client *Sms77API
}
type StatusCode string
type Sms77API struct {
	Options
	client *http.Client
	base   resource // Instead of allocating a struct for each service we reuse a one

	Analytics        *AnalyticsResource
	Balance          *BalanceResource
	Contacts         *ContactsResource
	Hooks            *HooksResource
	Journal          *JournalResource
	Lookup           *LookupResource
	Pricing          *PricingResource
	Sms              *SmsResource
	Status           *StatusResource
	ValidateForVoice *ValidateForVoiceResource
	Voice            *VoiceResource
}

const (
	defaultOptionSentWith = "go-client"

	sentWithKey = "sentWith"

	HttpMethodGet  HttpMethod = "GET"
	HttpMethodPost HttpMethod = "POST"

	StatusCodeErrorCarrierNotAvailable    StatusCode = "11"
	StatusCodeSuccess                     StatusCode = "100"
	StatusCodeSuccessPartial              StatusCode = "101"
	StatusCodeInvalidSender               StatusCode = "201"
	StatusCodeInvalidRecipient            StatusCode = "202"
	StatusCodeMissingParamTo              StatusCode = "301"
	StatusCodeMissingParamText            StatusCode = "305"
	StatusCodeParamTextExceedsLengthLimit StatusCode = "401"
	StatusCodePreventedByReloadLock       StatusCode = "402"
	StatusCodeReachedDailyLimitForNumber  StatusCode = "403"
	StatusCodeInsufficientCredits         StatusCode = "500"
	StatusCodeErrorCarrierDelivery        StatusCode = "600"
	StatusCodeErrorUnknown                StatusCode = "700"
	StatusCodeErrorAuthentication         StatusCode = "900"
	StatusCodeErrorApiDisabledForKey      StatusCode = "902"
	StatusCodeErrorServerIp               StatusCode = "903"
)

var StatusCodes = map[StatusCode]string{
	StatusCodeErrorCarrierNotAvailable:    "ErrorCarrierNotAvailable",
	StatusCodeSuccess:                     "Success",
	StatusCodeSuccessPartial:              "SuccessPartial",
	StatusCodeInvalidSender:               "InvalidSender",
	StatusCodeInvalidRecipient:            "InvalidRecipient",
	StatusCodeMissingParamTo:              "MissingParamTo",
	StatusCodeMissingParamText:            "MissingParamText",
	StatusCodeParamTextExceedsLengthLimit: "ParamTextExceedsLengthLimit",
	StatusCodePreventedByReloadLock:       "PreventedByReloadLock",
	StatusCodeReachedDailyLimitForNumber:  "ReachedDailyLimitForNumber",
	StatusCodeInsufficientCredits:         "InsufficientCredits",
	StatusCodeErrorCarrierDelivery:        "ErrorCarrierDelivery",
	StatusCodeErrorUnknown:                "ErrorUnknown",
	StatusCodeErrorAuthentication:         "ErrorAuthentication",
	StatusCodeErrorApiDisabledForKey:      "ErrorApiDisabledForKey",
	StatusCodeErrorServerIp:               "ErrorServerIp",
}

func New(options Options) *Sms77API {
	if "" == options.SentWith {
		options.SentWith = defaultOptionSentWith
	}

	c := &Sms77API{client: http.DefaultClient}
	c.Options = options
	c.base.client = c

	c.Analytics = (*AnalyticsResource)(&c.base)
	c.Balance = (*BalanceResource)(&c.base)
	c.Contacts = (*ContactsResource)(&c.base)
	c.Hooks = (*HooksResource)(&c.base)
	c.Journal = (*JournalResource)(&c.base)
	c.Lookup = (*LookupResource)(&c.base)
	c.Pricing = (*PricingResource)(&c.base)
	c.Sms = (*SmsResource)(&c.base)
	c.Status = (*StatusResource)(&c.base)
	c.ValidateForVoice = (*ValidateForVoiceResource)(&c.base)
	c.Voice = (*VoiceResource)(&c.base)

	return c
}

func (api *Sms77API) get(endpoint string, data map[string]interface{}) (string, error) {
	return api.request(endpoint, http.MethodGet, data)
}

func (api *Sms77API) post(endpoint string, data map[string]interface{}) (string, error) {
	return api.request(endpoint, http.MethodPost, data)
}

func (api *Sms77API) request(endpoint string, method string, data interface{}) (string, error) {
	createRequestPayload := func() string {
		params := url.Values{}

		for k, v := range data.(map[string]interface{}) {
			if api.Debug {
				log.Printf("%s: %v", k, v)
			}

			switch v.(type) {
			case nil:
				continue
			case bool:
				if true == v {
					v = "1"
				} else {
					v = "0"
				}
			case int64:
				v = strconv.FormatInt(v.(int64), 10)
			case []interface{}:
				for fileIndex, files := range v.([]interface{}) {
					for fileKey, fileValue := range files.(map[string]interface{}) {
						params.Add(
							fmt.Sprintf("%s[%d][%s]", k, fileIndex, fileKey), fmt.Sprintf("%v", fileValue))
					}
				}

				continue
			}

			params.Add(k, fmt.Sprintf("%v", v))
		}

		return params.Encode()
	}

	initClient := func() (req *http.Request, err error) {
		var uri = fmt.Sprintf("https://gateway.sms77.io/api/%s", endpoint)
		var qs = createRequestPayload()
		var headers = map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", api.ApiKey),
			sentWithKey:     api.SentWith,
		}
		var body = ""

		if http.MethodGet == method {
			if "" != qs {
				uri = fmt.Sprintf("%s?%s", uri, qs)
			}
		} else {
			body = qs

			headers["Content-Type"] = "application/x-www-form-urlencoded"
		}

		if api.Debug {
			log.Printf("%s %s", method, uri)
		}

		req, err = http.NewRequest(method, uri, strings.NewReader(body))

		if nil != err {
			log.Println(err.Error())
			panic(err)
		}

		for k, v := range headers {
			req.Header.Add(k, v)
		}

		return
	}

	if http.MethodGet != method && http.MethodPost != method {
		return "", errors.New(fmt.Sprintf("unsupported http method %s", method))
	}

	if "" == api.Options.ApiKey {
		return "", errors.New("missing required option ApiKey")
	}

	if nil == data {
		data = map[string]interface{}{}
	}
	data, _ = json.Marshal(&data)
	json.Unmarshal(data.([]byte), &data)

	req, err := initClient()
	if err != nil {
		return "", fmt.Errorf("could not execute request! #1 (%s)", err.Error())
	}

	res, err := api.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request! #2 (%s)", err.Error())
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("could not execute request! #3 (%s)", err.Error())
	}

	str := strings.TrimSpace(string(body))

	if api.Debug {
		log.Println(str)
	}

	length := len(str)

	if 2 == length || 3 == length {
		code, msg := pickMapByKey(str, StatusCodes)
		if nil != code {
			return "", errors.New(fmt.Sprintf("%s: %s", code, msg))
		}
	}

	return str, nil
}
