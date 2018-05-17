package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPMember is a IP Messaging Member resource.
type IPMember struct {
	Sid         string  `json:"sid"`
	AccountSid  string  `json:"account_sid"`
	ChannelSid  string  `json:"channel_sid"`
	ServiceSid  string  `json:"service_sid"`
	Identity    string  `json:"identity"`
	RoleSid     *string `json:"role_sid"`
	DateCreated string  `json:"date_created"`
	DateUpdated string  `json:"date_updated"`
	URL         string  `json:"url"`
}

// IPMemberList gives the results for querying the set of members. Returns the first page
// by default.
type IPMemberList struct {
	Client  Client
	Members []IPMember `json:"members"`
	Meta    Meta       `json:"meta"`
}

// AddIPMemberToChannel adds a member to a channel.
func AddIPMemberToChannel(client *TwilioIPMessagingClient, serviceSid string, channelSid string, identity string, roleSid string) (*IPMember, error) {
	var member *IPMember

	params := url.Values{}
	params.Set("Identity", identity)
	if roleSid != "" {
		params.Set("RoleSid", roleSid)
	}

	res, err := client.post(params, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Members.json")

	if err != nil {
		return member, err
	}

	member = new(IPMember)
	err = json.Unmarshal(res, member)

	return member, err
}

// GetIPChannelMember returns the specified IP Member in the channel.
func GetIPChannelMember(client *TwilioIPMessagingClient, serviceSid, channelSid, sid string) (*IPMember, error) {
	var member *IPMember

	res, err := client.get(url.Values{}, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Members/"+sid+".json")

	if err != nil {
		return nil, err
	}

	member = new(IPMember)
	err = json.Unmarshal(res, member)

	return member, err
}

// RemoveIPMemberFromChannel removes the given member from the channel.
func RemoveIPMemberFromChannel(client *TwilioIPMessagingClient, serviceSid, channelSid, sid string) error {
	return client.delete("/Services/" + serviceSid + "/Channels/" + channelSid + "/Members/" + sid)
}

// ListIPMembers returns the first page of members.
func ListIPMembers(client *TwilioIPMessagingClient, serviceSid, channelSid string) (*IPMemberList, error) {
	var memberList *IPMemberList

	body, err := client.get(nil, "/Services/"+serviceSid+"/Channels/"+channelSid+"/Members.json")

	if err != nil {
		return memberList, err
	}

	memberList = new(IPMemberList)
	memberList.Client = client
	err = json.Unmarshal(body, memberList)

	return memberList, err
}

// GetMembers recturns the current page of members.
func (c *IPMemberList) GetMembers() []IPMember {
	return c.Members
}

// GetAllMembers returns all of the members from all of the pages (from here forward).
func (c *IPMemberList) GetAllMembers() ([]IPMember, error) {
	members := c.Members
	t := c

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		members = append(members, t.Members...)
	}
	return members, nil
}

// HasNextPage returns whether or not there is a next page of members.
func (c *IPMemberList) HasNextPage() bool {
	return c.Meta.NextPageUri != ""
}

// NextPage returns the next page of members.
func (c *IPMemberList) NextPage() (*IPMemberList, error) {
	if !c.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return c.getPage(c.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (c *IPMemberList) HasPreviousPage() bool {
	return c.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of members.
func (c *IPMemberList) PreviousPage() (*IPMemberList, error) {
	if !c.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return c.getPage(c.Meta.NextPageUri)
}

// FirstPage returns the first page of members.
func (c *IPMemberList) FirstPage() (*IPMemberList, error) {
	return c.getPage(c.Meta.FirstPageUri)
}

// LastPage returns the last page of members.
func (c *IPMemberList) LastPage() (*IPMemberList, error) {
	return c.getPage(c.Meta.LastPageUri)
}

func (c *IPMemberList) getPage(uri string) (*IPMemberList, error) {
	var memberList *IPMemberList

	client := c.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return memberList, err
	}

	memberList = new(IPMemberList)
	memberList.Client = client
	err = json.Unmarshal(body, memberList)

	return memberList, err
}
