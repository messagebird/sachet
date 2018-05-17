package twiliogo

import (
  "encoding/json"
  "net/url"
)

type IncomingPhoneNumberList struct {
  Client Client
  Start int `json:"start"`
  Total int `json:"total"`
  NumPages int `json:"num_pages"`
  Page int `json:"page"`
  PageSize int `json:"page_size"`
  End int `json:"end"`
  Uri string `json:"uri"`
  FirstPageUri string `json:"first_page_uri"`
  LastPageUri string `json:"last_page_uri"`
  NextPageUri string `json:"next_page_uri"`
  PreviousPageUri string `json"previous_page_uri"`
  IncomingPhoneNumbers []IncomingPhoneNumber `json:"sms_messages"`
}

func GetIncomingPhoneNumberList(client Client, optionals ...Optional) (*IncomingPhoneNumberList, error) {
  var incomingPhoneNumberList *IncomingPhoneNumberList

  params := url.Values{}

  for _, optional := range optionals {
    param, value := optional.GetParam()
    params.Set(param, value)
  }

  body, err := client.get(params, client.RootUrl() + "/IncomingPhoneNumbers.json")

  if err != nil {
    return nil, err
  }

  incomingPhoneNumberList = new(IncomingPhoneNumberList)
  incomingPhoneNumberList.Client = client
  err = json.Unmarshal(body, incomingPhoneNumberList)

  return incomingPhoneNumberList, err
}

