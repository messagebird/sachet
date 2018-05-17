package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPMessage is a IP Messaging Message resource.
type IPMessage struct {
	Sid         string `json:"sid"`
	AccountSid  string `json:"account_sid"`
	ServiceSid  string `json:"service_sid"`
	To          string `json:"to"` // channel sid
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	WasEdited   bool   `json:"was_edited"`
	From        string `json:"from"` // identity
	Body        string `json:"body"`
	URL         string `json:"url"`
}

// IPMessageList gives the results for querying the set of messages. Returns the first page
// by default.
type IPMessageList struct {
	Client   Client
	Messages []IPMessage `json:"messages"`
	Meta     Meta        `json:"meta"`
}

// SendIPMessageToChannel sends a message to a channel.
func SendIPMessageToChannel(client *TwilioIPMessagingClient, serviceSid string, channelSid string, from string, body string) (*IPMessage, error) {
	var message *IPMessage

	params := url.Values{}
	params.Set("Body", body)
	if from != "" {
		params.Set("From", from)
	}

	res, err := client.post(params, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages.json")

	if err != nil {
		return message, err
	}

	message = new(IPMessage)
	err = json.Unmarshal(res, message)

	return message, err
}

// GetIPChannelMessage returns the specified IP Message in the channel.
func GetIPChannelMessage(client *TwilioIPMessagingClient, serviceSid, channelSid, sid string) (*IPMessage, error) {
	var message *IPMessage

	res, err := client.get(url.Values{}, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages/"+sid+".json")

	if err != nil {
		return nil, err
	}

	message = new(IPMessage)
	err = json.Unmarshal(res, message)

	return message, err
}

// ListIPMessages returns the first page of messages for a channel.
func ListIPMessages(client *TwilioIPMessagingClient, serviceSid, channelSid string) (*IPMessageList, error) {
	var messageList *IPMessageList

	body, err := client.get(nil, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Messages.json")

	if err != nil {
		return messageList, err
	}

	messageList = new(IPMessageList)
	messageList.Client = client
	err = json.Unmarshal(body, messageList)

	return messageList, err
}

// GetMessages recturns the current page of messages.
func (c *IPMessageList) GetMessages() []IPMessage {
	return c.Messages
}

// GetAllMessages returns all of the messages from all of the pages (from here forward).
func (c *IPMessageList) GetAllMessages() ([]IPMessage, error) {
	messages := c.Messages
	t := c

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		messages = append(messages, t.Messages...)
	}
	return messages, nil
}

// HasNextPage returns whether or not there is a next page of messages.
func (c *IPMessageList) HasNextPage() bool {
	return c.Meta.NextPageUri != ""
}

// NextPage returns the next page of messages.
func (c *IPMessageList) NextPage() (*IPMessageList, error) {
	if !c.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return c.getPage(c.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (c *IPMessageList) HasPreviousPage() bool {
	return c.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of messages.
func (c *IPMessageList) PreviousPage() (*IPMessageList, error) {
	if !c.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return c.getPage(c.Meta.NextPageUri)
}

// FirstPage returns the first page of messages.
func (c *IPMessageList) FirstPage() (*IPMessageList, error) {
	return c.getPage(c.Meta.FirstPageUri)
}

// LastPage returns the last page of messages.
func (c *IPMessageList) LastPage() (*IPMessageList, error) {
	return c.getPage(c.Meta.LastPageUri)
}

func (c *IPMessageList) getPage(uri string) (*IPMessageList, error) {
	var messageList *IPMessageList

	client := c.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return messageList, err
	}

	messageList = new(IPMessageList)
	messageList.Client = client
	err = json.Unmarshal(body, messageList)

	return messageList, err
}
