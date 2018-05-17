package twiliogo

import (
  "encoding/json"
  "net/url"
)

type MessageList struct {
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
  Messages []Message `json:"sms_messages"`
}

func GetMessageList(client Client, optionals ...Optional) (*MessageList, error) {
  var messageList *MessageList

  params := url.Values{}

  for _, optional := range optionals {
    param, value := optional.GetParam()
    params.Set(param, value)
  }

  body, err := client.get(params, client.RootUrl() + "/SMS/Messages.json")

  if err != nil {
    return messageList, err
  }

  messageList = new(MessageList)
  messageList.Client = client
  err = json.Unmarshal(body, messageList)

  return messageList, err
}

func (m *MessageList) GetMessages() []Message {
  return m.Messages
}

func (currentMessageList *MessageList) HasNextPage() bool {
  return currentMessageList.NextPageUri != ""
}

func (currentMessageList *MessageList) NextPage() (*MessageList, error) {
  if !currentMessageList.HasNextPage() {
    return nil, Error {"No next page"}
  }

  return currentMessageList.getPage(currentMessageList.NextPageUri)
}

func (currentMessageList *MessageList) HasPreviousPage() bool {
  return currentMessageList.PreviousPageUri != ""
}

func (currentMessageList *MessageList) PreviousPage() (*MessageList, error) {
  if !currentMessageList.HasPreviousPage() {
    return nil, Error {"No previous page"}
  }

  return currentMessageList.getPage(currentMessageList.NextPageUri)
}

func (currentMessageList *MessageList) FirstPage() (*MessageList, error) {
  return currentMessageList.getPage(currentMessageList.FirstPageUri)
}

func (currentMessageList *MessageList) LastPage() (*MessageList, error) {
  return currentMessageList.getPage(currentMessageList.LastPageUri)
}

func (currentMessageList *MessageList) getPage(uri string) (*MessageList, error) {
  var messageList *MessageList

  client := currentMessageList.Client

  body, err := client.get(nil, uri)

  if err != nil {
    return messageList, err
  }

  messageList = new(MessageList)
  messageList.Client = client
  err = json.Unmarshal(body, messageList)

  return messageList, err
}
