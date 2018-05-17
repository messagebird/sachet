package twiliogo

import (
  "encoding/json"
  "net/url"
)

type Message struct {
  Sid string `json:"sid"`
  DateCreated string `json:"date_created"`
  DateUpdated string `json:"date_updated"`
  DateSent string `json:"date_sent"`
  AccountSid string `json:"account_sid"`
  From string `json:"from"`
  To string `json:"to"`
  Body string `json:"body"`
  NumSegments string `json:"num_segments"`
  Status string `json:"status"`
  Direction string `json:"direction"`
  Price string `json:"price"`
  PriceUnit string `json:"price_unit"`
  ApiVersion string `json:"api_version"`
  Uri string `json:"uri"`
}

func NewMessage(client Client, from string, to string, content ...Optional) (*Message, error) {
  var message *Message

  params := url.Values{}
  params.Set("From", from)
  params.Set("To", to)

  if len(content) < 1 {
    return nil, Error { "Must have at least a Body or MediaUrl" }
  }

  for _, optional := range content {
    param, value := optional.GetParam()

    if param != "Body" && param != "MediaUrl" {
      return nil, Error { "Only Body or MediaUrl allowed" }
    }

    params.Set(param, value)
  }

  res, err := client.post(params, client.RootUrl() + "/Messages.json")

  if err != nil {
    return message, err
  }

  message = new(Message)
  err = json.Unmarshal(res, message)

  return message, err
}

func GetMessage(client Client, sid string) (*Message, error) {
  var message *Message

  res, err := client.get(url.Values{}, client.RootUrl() + "/Messages/" + sid + ".json")

  if err != nil {
    return nil, err
  }

  message = new(Message)
  err = json.Unmarshal(res, message)

  return message, err
}
