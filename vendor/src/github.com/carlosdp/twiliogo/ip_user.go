package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPUser is a IP Messaging User resource.
type IPUser struct {
	Sid         string `json:"sid"`
	AccountSid  string `json:"account_sid"`
	ServiceSid  string `json:"service_sid"`
	RoleSid     string `json:"role_sid"`
	Identity    string `json:"identity"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	URL         string `json:"url"`
}

// IPUserList gives the results for querying the set of users. Returns the first page
// by default.
type IPUserList struct {
	Client *TwilioIPMessagingClient
	Users  []IPUser `json:"users"`
	Meta   Meta     `json:"meta"`
}

// NewIPUser creates a new IP Messaging User.
func NewIPUser(client *TwilioIPMessagingClient, serviceSid string, identity string, roleSid string) (*IPUser, error) {
	var user *IPUser

	params := url.Values{}
	params.Set("Identity", identity)
	params.Set("RoleSid", roleSid)

	res, err := client.post(params, "/Services/"+serviceSid+"/Users.json")

	if err != nil {
		return user, err
	}

	user = new(IPUser)
	err = json.Unmarshal(res, user)

	return user, err
}

// GetIPUser returns information on the specified user.
func GetIPUser(client *TwilioIPMessagingClient, serviceSid, sid string) (*IPUser, error) {
	var user *IPUser

	res, err := client.get(url.Values{}, "/Services/"+serviceSid+"/Users/"+sid+".json")

	if err != nil {
		return nil, err
	}

	user = new(IPUser)
	err = json.Unmarshal(res, user)

	return user, err
}

// DeleteIPUser deletes the given IP user.
func DeleteIPUser(client *TwilioIPMessagingClient, serviceSid, sid string) error {
	return client.delete("/Services/" + serviceSid + "/Users/" + sid)
}

// UpdateIPUser updates an existing IP Messaging user.
func UpdateIPUser(client *TwilioIPMessagingClient, serviceSid string, sid string, identity string, roleSid string) (*IPUser, error) {
	var user *IPUser

	params := url.Values{}
	params.Set("Identity", identity)
	params.Set("RoleSid", roleSid)

	res, err := client.post(params, "/Services/"+serviceSid+"/Users/"+sid+".json")

	if err != nil {
		return user, err
	}

	user = new(IPUser)
	err = json.Unmarshal(res, user)

	return user, err
}

// ListIPUsers returns the first page of users.
func ListIPUsers(client *TwilioIPMessagingClient, serviceSid string) (*IPUserList, error) {
	var userList *IPUserList

	body, err := client.get(nil, "/Services/"+serviceSid+"/Users.json")

	if err != nil {
		return userList, err
	}

	userList = new(IPUserList)
	userList.Client = client
	err = json.Unmarshal(body, userList)

	return userList, err
}

// GetUsers returns the current page of users.
func (s *IPUserList) GetUsers() []IPUser {
	return s.Users
}

// GetAllUsers returns all of the users from all of the pages (from here forward).
func (s *IPUserList) GetAllUsers() ([]IPUser, error) {
	users := s.Users
	t := s

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		users = append(users, t.Users...)
	}
	return users, nil
}

// HasNextPage returns whether or not there is a next page of users.
func (s *IPUserList) HasNextPage() bool {
	return s.Meta.NextPageUri != ""
}

// NextPage returns the next page of users.
func (s *IPUserList) NextPage() (*IPUserList, error) {
	if !s.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (s *IPUserList) HasPreviousPage() bool {
	return s.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of users.
func (s *IPUserList) PreviousPage() (*IPUserList, error) {
	if !s.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// FirstPage returns the first page of users.
func (s *IPUserList) FirstPage() (*IPUserList, error) {
	return s.getPage(s.Meta.FirstPageUri)
}

// LastPage returns the last page of users.
func (s *IPUserList) LastPage() (*IPUserList, error) {
	return s.getPage(s.Meta.LastPageUri)
}

func (s *IPUserList) getPage(uri string) (*IPUserList, error) {
	var userList *IPUserList

	client := s.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return userList, err
	}

	userList = new(IPUserList)
	userList.Client = client
	err = json.Unmarshal(body, userList)

	return userList, err
}
