package twiliogo

import (
  "net/http"
  "strings"
  "net/url"
  "io/ioutil"
  "encoding/json"
)

const ROOT = "https://api.twilio.com"
const VERSION = "2010-04-01"

type Client interface {
  AccountSid() string
  AuthToken() string
  RootUrl() string
  get(url.Values, string) ([]byte, error)
  post(url.Values, string) ([]byte, error)
}

type TwilioClient struct {
  accountSid string
  authToken string
  rootUrl string
}

func NewClient(accountSid, authToken string) *TwilioClient {
  rootUrl := "/" + VERSION + "/Accounts/" + accountSid
  return &TwilioClient{accountSid, authToken, rootUrl}
}

func (client *TwilioClient) post(values url.Values, uri string) ([]byte, error) {
  req, err := http.NewRequest("POST", ROOT + uri, strings.NewReader(values.Encode()))

  if err != nil {
    return nil, err
  }

  req.SetBasicAuth(client.AccountSid(), client.AuthToken())
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  httpClient := &http.Client{}

  res, err := httpClient.Do(req)

  if err != nil {
    return nil, err
  }

  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)

  if err != nil {
    return body, err
  }

  if res.StatusCode != 200 && res.StatusCode != 201 {
    if res.StatusCode == 500 {
      return body, Error{"Server Error"}
    } else {
      twilioError := new(TwilioError)
      json.Unmarshal(body, twilioError)
      return body, twilioError
    }
  }

  return body, err
}

func (client *TwilioClient) get(queryParams url.Values, uri string) ([]byte, error) {
  var params *strings.Reader

  if queryParams == nil {
    queryParams = url.Values{}
  }

  params = strings.NewReader(queryParams.Encode())
  req, err := http.NewRequest("GET", ROOT + uri, params)

  if err != nil {
    return nil, err
  }

  req.SetBasicAuth(client.AccountSid(), client.AuthToken())
  httpClient := &http.Client{}

  res, err := httpClient.Do(req)

  if err != nil {
    return nil, err
  }

  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)

  if err != nil {
    return body, err
  }

  if res.StatusCode != 200 && res.StatusCode != 201 {
    if res.StatusCode == 500 {
      return body, Error{"Server Error"}
    } else {
      twilioError := new(TwilioError)
      json.Unmarshal(body, twilioError)
      return body, twilioError
    }
  }

  return body, err
}

func (client *TwilioClient) AccountSid() string {
  return client.accountSid
}

func (client *TwilioClient) AuthToken() string {
  return client.authToken
}

func (client *TwilioClient) RootUrl() string {
  return client.rootUrl
}
