package twiliogo

import (
	"encoding/json"
	"net/url"
)

// IPCredential is a IP Messaging Credential resource.
type IPCredential struct {
	Sid          string `json:"sid"`
	AccountSid   string `json:"account_sid"`
	FriendlyName string `json:"friendly_name"`
	Type         string `json:"type"` // apns or gcm
	Sandbox      bool   `json:"sandbox"`
	URL          string `json:"url"`
}

// IPCredentialList gives the results for querying the set of credentials. Returns the first page
// by default.
type IPCredentialList struct {
	Client      Client
	Credentials []IPCredential `json:"credentials"`
	Meta        Meta           `json:"meta"`
}

// NewIPCredential creates a new IP Messaging Credential.
// Kind must be apns or gcm.
func NewIPCredential(client *TwilioIPMessagingClient, friendlyName string, kind string, sandbox bool, apnsCert string, apnsPrivateKey string,
	gcmApiKey string) (*IPCredential, error) {
	var credential *IPCredential

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("Type", kind)
	if sandbox {
		params.Set("Sandbox", "true")
	} else {
		params.Set("Sandbox", "false")
	}
	if apnsCert != "" {
		params.Set("Certificate", apnsCert)
	}
	if apnsPrivateKey != "" {
		params.Set("PrivateKey", apnsPrivateKey)
	}
	if gcmApiKey != "" {
		params.Set("ApiKey", gcmApiKey)
	}

	res, err := client.post(params, "/Credentials.json")

	if err != nil {
		return credential, err
	}

	credential = new(IPCredential)
	err = json.Unmarshal(res, credential)

	return credential, err
}

// GetIPCredential returns information on the specified credential.
func GetIPCredential(client *TwilioIPMessagingClient, sid string) (*IPCredential, error) {
	var credential *IPCredential

	res, err := client.get(url.Values{}, "/Credentials/"+sid+".json")

	if err != nil {
		return nil, err
	}

	credential = new(IPCredential)
	err = json.Unmarshal(res, credential)

	return credential, err
}

// DeleteIPCredential deletes the given IP Credential.
func DeleteIPCredential(client *TwilioIPMessagingClient, sid string) error {
	return client.delete("/Credentials/" + sid)
}

// UpdateIPCredential updates an existing IP Messaging Credential.
func UpdateIPCredential(client *TwilioIPMessagingClient, sid string, friendlyName string, kind string, sandbox bool) (*IPCredential, error) {
	var credential *IPCredential

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("Type", kind)
	if sandbox {
		params.Set("Sandbox", "true")
	} else {
		params.Set("Sandbox", "false")
	}

	res, err := client.post(params, "/Credentials/"+sid+".json")

	if err != nil {
		return credential, err
	}

	credential = new(IPCredential)
	err = json.Unmarshal(res, credential)

	return credential, err
}

// ListIPCredentials returns the first page of credentials.
func ListIPCredentials(client *TwilioIPMessagingClient) (*IPCredentialList, error) {
	var credentialList *IPCredentialList

	body, err := client.get(nil, "/Credentials.json")

	if err != nil {
		return credentialList, err
	}

	credentialList = new(IPCredentialList)
	credentialList.Client = client
	err = json.Unmarshal(body, credentialList)

	return credentialList, err
}

// GetCredentials returns the current page of credentials.
func (s *IPCredentialList) GetCredentials() []IPCredential {
	return s.Credentials
}

// GetAllCredentials returns all of the credentials from all of the pages (from here forward).
func (s *IPCredentialList) GetAllCredentials() ([]IPCredential, error) {
	credentials := s.Credentials
	t := s

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, t.Credentials...)
	}
	return credentials, nil
}

// HasNextPage returns whether or not there is a next page of credentials.
func (s *IPCredentialList) HasNextPage() bool {
	return s.Meta.NextPageUri != ""
}

// NextPage returns the next page of credentials.
func (s *IPCredentialList) NextPage() (*IPCredentialList, error) {
	if !s.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (s *IPCredentialList) HasPreviousPage() bool {
	return s.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of credentials.
func (s *IPCredentialList) PreviousPage() (*IPCredentialList, error) {
	if !s.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// FirstPage returns the first page of credentials.
func (s *IPCredentialList) FirstPage() (*IPCredentialList, error) {
	return s.getPage(s.Meta.FirstPageUri)
}

// LastPage returns the last page of credentials.
func (s *IPCredentialList) LastPage() (*IPCredentialList, error) {
	return s.getPage(s.Meta.LastPageUri)
}

func (s *IPCredentialList) getPage(uri string) (*IPCredentialList, error) {
	var credentialList *IPCredentialList

	client := s.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return credentialList, err
	}

	credentialList = new(IPCredentialList)
	credentialList.Client = client
	err = json.Unmarshal(body, credentialList)

	return credentialList, err
}
