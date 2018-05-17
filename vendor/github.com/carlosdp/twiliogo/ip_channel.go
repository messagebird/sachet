package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPChannel is a IP Messaging Channel resource.
type IPChannel struct {
	Sid          string            `json:"sid"`
	AccountSid   string            `json:"account_sid"`
	ServiceSid   string            `json:"service_sid"`
	FriendlyName string            `json:"friendly_name"`
	UniqueName   string            `json:"unique_name"`
	Attributes   string            `json:"attributes"`
	Type         string            `json:"type"`
	DateCreated  string            `json:"date_created"`
	DateUpdated  string            `json:"date_updated"`
	CreatedBy    string            `json:"created_by"`
	URL          string            `json:"url"`
	Links        map[string]string `json:"links"`
}

// IPChannelList gives the results for querying the set of channels. Returns the first page
// by default.
type IPChannelList struct {
	Client   Client
	Channels []IPChannel `json:"channels"`
	Meta     Meta        `json:"meta"`
}

// NewIPChannel creates a new IP Messaging Channel.
func NewIPChannel(client *TwilioIPMessagingClient, serviceSid string, friendlyName string, uniqueName string, public bool, attributes string) (*IPChannel, error) {
	var channel *IPChannel

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("UniqueName", uniqueName)
	kind := "private"
	if public {
		kind = "public"
	}
	params.Set("Type", kind)
	params.Set("Attributes", attributes)

	res, err := client.post(params, "/Services/"+serviceSid+"/Channels.json")

	if err != nil {
		return channel, err
	}

	channel = new(IPChannel)
	err = json.Unmarshal(res, channel)

	return channel, err
}

// UpdateIPChannel updates ane existing IP Messaging Channel.
func UpdateIPChannel(client *TwilioIPMessagingClient, serviceSid string, sid string, friendlyName string, uniqueName string, public bool, attributes string) (*IPChannel, error) {
	var channel *IPChannel

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("UniqueName", uniqueName)
	kind := "private"
	if public {
		kind = "public"
	}
	params.Set("Type", kind)
	params.Set("Attributes", attributes)

	res, err := client.post(params, "/Services/"+serviceSid+"/Channels/"+sid+".json")

	if err != nil {
		return channel, err
	}

	channel = new(IPChannel)
	err = json.Unmarshal(res, channel)

	return channel, err
}

// GetIPChannel returns the specified IP Channel.
func GetIPChannel(client *TwilioIPMessagingClient, serviceSid string, sid string) (*IPChannel, error) {
	var channel *IPChannel

	res, err := client.get(url.Values{}, "/Services/"+serviceSid+"/Channels/"+sid+".json")

	if err != nil {
		return nil, err
	}

	channel = new(IPChannel)
	err = json.Unmarshal(res, channel)

	return channel, err
}

// DeleteIPChannel deletes the given IP Channel.
func DeleteIPChannel(client *TwilioIPMessagingClient, serviceSid, sid string) error {
	return client.delete("/Services/" + serviceSid + "/Channels/" + sid)
}

// ListIPChannels returns the first page of channels.
func ListIPChannels(client *TwilioIPMessagingClient, serviceSid string) (*IPChannelList, error) {
	var channelList *IPChannelList

	body, err := client.get(nil, "/Services/"+serviceSid+"/Channels.json")

	if err != nil {
		return channelList, err
	}

	channelList = new(IPChannelList)
	channelList.Client = client
	err = json.Unmarshal(body, channelList)

	return channelList, err
}

// GetChannels recturns the current page of channels.
func (c *IPChannelList) GetChannels() []IPChannel {
	return c.Channels
}

// GetAllChannels returns all of the channels from all of the pages (from here forward).
func (c *IPChannelList) GetAllChannels() ([]IPChannel, error) {
	channels := c.Channels
	t := c

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		channels = append(channels, t.Channels...)
	}
	return channels, nil
}

// HasNextPage returns whether or not there is a next page of channels.
func (c *IPChannelList) HasNextPage() bool {
	return c.Meta.NextPageUri != ""
}

// NextPage returns the next page of channels.
func (c *IPChannelList) NextPage() (*IPChannelList, error) {
	if !c.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return c.getPage(c.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (c *IPChannelList) HasPreviousPage() bool {
	return c.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of channels.
func (c *IPChannelList) PreviousPage() (*IPChannelList, error) {
	if !c.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return c.getPage(c.Meta.NextPageUri)
}

// FirstPage returns the first page of channels.
func (c *IPChannelList) FirstPage() (*IPChannelList, error) {
	return c.getPage(c.Meta.FirstPageUri)
}

// LastPage returns the last page of channels.
func (c *IPChannelList) LastPage() (*IPChannelList, error) {
	return c.getPage(c.Meta.LastPageUri)
}

func (c *IPChannelList) getPage(uri string) (*IPChannelList, error) {
	var channelList *IPChannelList

	client := c.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return channelList, err
	}

	channelList = new(IPChannelList)
	channelList.Client = client
	err = json.Unmarshal(body, channelList)

	return channelList, err
}
