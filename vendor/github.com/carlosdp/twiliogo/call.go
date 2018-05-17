package twiliogo

import (
  "net/url"
  "encoding/json"
)

type Call struct {
  Sid string `json:"sid"`
  ParentCallSid string `json:"parent_call_sid"`
  DateCreated string `json:"date_created"`
  DateUpdated string `json:"date_updated"`
  AccountSid string `json:"account_sid"`
  To string `json:"to"`
  From string `json:"from"`
  PhoneNumberSid string `json:"phone_number_sid"`
  Status string `json:"status"`
  StartTime string `json:"start_time"`
  EndTime string `json:"end_time"`
  Duration string `json:"duration"`
  Price string `json:"price"`
  PriceUnit string `json:"price_unit"`
  Direction string `json:"direction"`
  AnsweredBy string `json:"answered_by"`
  ForwardedFrom string `json:"forwarded_from"`
  CallerName string `json:"caller_name"`
  Uri string `json:"uri"`
}

func NewCall(client Client, from, to string, callback Optional, optionals ...Optional) (*Call, error) {
  var call *Call

  params := url.Values{}
  params.Set("From", from)
  params.Set("To", to)

  callbackType, callbackParam := callback.GetParam()
  params.Set(callbackType, callbackParam)

  for _, optional := range optionals {
    param, value := optional.GetParam()
    params.Set(param, value)
  }

  res, err := client.post(params, client.RootUrl() + "/Calls.json")

  if err != nil {
    return nil, err
  }

  call = new(Call)
  err = json.Unmarshal(res, call)

  return call, err
}

func GetCall(client Client, sid string) (*Call, error) {
  var call *Call

  res, err := client.get(url.Values{}, client.RootUrl() + "/Calls/" + sid + ".json")

  if err != nil {
    return nil, err
  }

  call = new(Call)
  err = json.Unmarshal(res, call)

  return call, err
}

func (call *Call) Update(client Client, optionals ...Optional) error {
  var tempCall *Call

  params := url.Values{}

  for _, optional := range optionals {
    param, value := optional.GetParam()
    params.Set(param, value)
  }

  res, err := client.post(params, client.RootUrl() + "/Calls/" + call.Sid + ".json")

  if err != nil {
    return err
  }

  tempCall = new(Call)
  err = json.Unmarshal(res, tempCall)

  if err == nil {
    call = tempCall
  }

  return err
}
