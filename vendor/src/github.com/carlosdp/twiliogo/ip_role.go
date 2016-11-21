package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPRole is a IP Messaging Role resource.
type IPRole struct {
	Sid          string   `json:"sid"`
	AccountSid   string   `json:"account_sid"`
	ServiceSid   string   `json:"service_sid"`
	FriendlyName string   `json:"friendly_name"`
	Type         string   `json:"type"`
	Permissions  []string `json:"permissions"`
	DateCreated  string   `json:"date_created"`
	DateUpdated  string   `json:"date_updated"`
	URL          string   `json:"url"`
}

// IPRoleList gives the results for querying the set of roles. Returns the first page
// by default.
type IPRoleList struct {
	Client *TwilioIPMessagingClient
	Roles  []IPRole `json:"roles"`
	Meta   Meta     `json:"meta"`
}

// Permissions allowed for IP Roles.
const (
	PermissionCreateChannel         = "createChannel"
	PermissionJoinChannel           = "joinChannel"
	PermissionDestroyChannel        = "destroyChannel"
	PermissionInviteMember          = "inviteMember"
	PermissionRemoveMember          = "removeMember"
	PermissionEditChannelName       = "editChannelName"
	PermissionEditChannelAttributes = "editChannelAttributes"
	PermissionAddMember             = "addMember"
	PermissionEditAnyMessage        = "editAnyMessage"
	PermissionDeleteAnyMessage      = "deleteAnyMessage"
	PermissionSendMessage           = "sendMessage"
	PermissionLeaveChannel          = "leaveChannel"
	PermissionEditOwnMessage        = "editOwnMessage"
	PermissionDeleteOwnMessage      = "deleteOwnMessage"
)

// NewIPRole creates a new IP Messaging Role.
// kind should be "channel" or "service".
// permissions should be a subset of the permissions consts above.
func NewIPRole(client *TwilioIPMessagingClient, serviceSid string, friendlyName string, kind string, permissions []string) (*IPRole, error) {
	var role *IPRole

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("Type", kind)
	if permissions != nil {
		for _, p := range permissions {
			params.Add("Permission", p)
		}
	}

	res, err := client.post(params, "/Services/"+serviceSid+"/Roles.json")

	if err != nil {
		return role, err
	}

	role = new(IPRole)
	err = json.Unmarshal(res, role)

	return role, err
}

// GetIPRole returns information on the specified role.
func GetIPRole(client *TwilioIPMessagingClient, serviceSid, sid string) (*IPRole, error) {
	var role *IPRole

	res, err := client.get(url.Values{}, "/Services/"+serviceSid+"/Roles/"+sid+".json")

	if err != nil {
		return nil, err
	}

	role = new(IPRole)
	err = json.Unmarshal(res, role)

	return role, err
}

// DeleteIPRole deletes the given IP Role.
func DeleteIPRole(client *TwilioIPMessagingClient, serviceSid, sid string) error {
	return client.delete("/Services/" + serviceSid + "/Roles/" + sid)
}

// UpdateIPRole updates an existing IP Messaging Role.
func UpdateIPRole(client *TwilioIPMessagingClient, serviceSid string, sid string, friendlyName string, kind string, permissions []string) (*IPRole, error) {
	var role *IPRole

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("Type", kind)
	if permissions != nil {
		for _, p := range permissions {
			params.Add("Permission", p)
		}
	}

	res, err := client.post(params, "/Services/"+serviceSid+"/Roles/"+sid+".json")

	if err != nil {
		return role, err
	}

	role = new(IPRole)
	err = json.Unmarshal(res, role)

	return role, err
}

// ListIPRoles returns the first page of roles.
func ListIPRoles(client *TwilioIPMessagingClient, serviceSid string) (*IPRoleList, error) {
	var roleList *IPRoleList

	body, err := client.get(nil, "/Services/"+serviceSid+"/Roles.json")

	if err != nil {
		return roleList, err
	}

	roleList = new(IPRoleList)
	roleList.Client = client
	err = json.Unmarshal(body, roleList)

	return roleList, err
}

// GetRoles returns the current page of roles.
func (s *IPRoleList) GetRoles() []IPRole {
	return s.Roles
}

// GetAllRoles returns all of the roles from all of the pages (from here forward).
func (s *IPRoleList) GetAllRoles() ([]IPRole, error) {
	roles := s.Roles
	t := s

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		roles = append(roles, t.Roles...)
	}
	return roles, nil
}

// HasNextPage returns whether or not there is a next page of roles.
func (s *IPRoleList) HasNextPage() bool {
	return s.Meta.NextPageUri != ""
}

// NextPage returns the next page of roles.
func (s *IPRoleList) NextPage() (*IPRoleList, error) {
	if !s.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (s *IPRoleList) HasPreviousPage() bool {
	return s.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of roles.
func (s *IPRoleList) PreviousPage() (*IPRoleList, error) {
	if !s.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// FirstPage returns the first page of roles.
func (s *IPRoleList) FirstPage() (*IPRoleList, error) {
	return s.getPage(s.Meta.FirstPageUri)
}

// LastPage returns the last page of roles.
func (s *IPRoleList) LastPage() (*IPRoleList, error) {
	return s.getPage(s.Meta.LastPageUri)
}

func (s *IPRoleList) getPage(uri string) (*IPRoleList, error) {
	var roleList *IPRoleList

	client := s.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return roleList, err
	}

	roleList = new(IPRoleList)
	roleList.Client = client
	err = json.Unmarshal(body, roleList)

	return roleList, err
}
