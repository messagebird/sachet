package twiliogo

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"time"
)

// IPService is a IP Messaging Service resource.
type IPService struct {
	Sid                    string            `json:"sid"`
	AccountSid             string            `json:"account_sid"`
	FriendlyName           string            `json:"friendly_name"`
	DateCreated            string            `json:"date_created"`
	DateUpdated            string            `json:"date_updated"`
	DefaultServiceRoleSid  string            `json:"default_service_role_sid"`
	DefaultChannelRoleSid  string            `json:"default_channel_role_sid"`
	TypingIndicatorTimeout uint              `json:"typing_indicator_timeout"`
	Webhooks               map[string]string `json:"webhooks"`
	URL                    string            `json:"url"`
	Links                  map[string]string `json:"links"`
}

// Meta is a metadata type for the IP messaging services.
type Meta struct {
	Start           int    `json:"start"`
	Total           int    `json:"total"`
	NumPages        int    `json:"num_pages"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	End             int    `json:"end"`
	Uri             string `json:"uri"`
	FirstPageUri    string `json:"first_page_uri"`
	LastPageUri     string `json:"last_page_uri"`
	NextPageUri     string `json:"next_page_uri"`
	PreviousPageUri string `json:"previous_page_uri"`
	Key             string `json:"key"`
}

// IPServiceList gives the results for querying the set of services. Returns the first page
// by default.
type IPServiceList struct {
	Client   Client
	Services []IPService `json:"services"`
	Meta     Meta        `json:"meta"`
}

// Webhooks available for services to specify
const (
	WebhookOnMessageSend    = "Webhooks.OnMessageSend"
	WebhookOnMessageRemove  = "Webhooks.OnMessageRemove"
	WebhookOnMessageUpdate  = "Webhooks.OnMessageUpdate"
	WebhookOnChannelAdd     = "Webhooks.OnChannelAdd"
	WebhookOnChannelUpdate  = "Webhooks.OnChannelUpdate"
	WebhookOnChannelDestroy = "Webhooks.OnChannelDestroy"
	WebhookOnMemberAdd      = "Webhooks.OnMemberAdd"
	WebhookOnMemberRemove   = "Webhooks.OnMemberRemove"
)

// Webhooks are used to define push webhooks for an IP service.
type Webhooks map[string]string

// NewWebhooks creates a new, empty set of web hooks.
func NewWebhooks() Webhooks {
	return Webhooks(make(map[string]string))
}

// Add adds a new webhook. The name should be one of the Webhook* exported values.
// Method is the HTTP method (e.g., "POST"). Format should be "xml" or "json".
func (w Webhooks) Add(name, method, format, url string) {
	w[name+".Method"] = method
	w[name+".Format"] = format
	w[name+".Url"] = url
}

func durationToISO8601(d time.Duration) (string, error) {
	if d > time.Hour {
		return "", fmt.Errorf("Duration is too long: %v", d)
	}
	minutes := int(math.Floor(d.Minutes()))
	seconds := int(math.Floor(d.Minutes()-float64(minutes)) * 60.0)
	return fmt.Sprintf("PT%dM%dS", minutes, seconds), nil
}

// NewIPService creates a new IP Messaging Service.
func NewIPService(client *TwilioIPMessagingClient, friendlyName string, defaultServiceRoleSid string, defaultChannelRoleSid string,
	typingIndicatorTimeout time.Duration, webhooks Webhooks) (*IPService, error) {

	timeout, err := durationToISO8601(typingIndicatorTimeout)
	if err != nil {
		return nil, err
	}

	var service *IPService

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("DefaultServiceRoleSid", defaultServiceRoleSid)
	params.Set("DefaultChannelRoleSid", defaultChannelRoleSid)
	params.Set("TypingIndicatonTimeout", timeout)
	if webhooks != nil {
		for k, v := range webhooks {
			params.Set(k, v)
		}
	}

	res, err := client.post(params, "/Services.json")

	if err != nil {
		return service, err
	}

	service = new(IPService)
	err = json.Unmarshal(res, service)

	return service, err
}

// GetIPService returns information on the specified service.
func GetIPService(client *TwilioIPMessagingClient, sid string) (*IPService, error) {
	var service *IPService

	res, err := client.get(url.Values{}, "/Services/"+sid+".json")

	if err != nil {
		return nil, err
	}

	service = new(IPService)
	err = json.Unmarshal(res, service)

	return service, err
}

// DeleteIPService deletes the given IP Service.
func DeleteIPService(client *TwilioIPMessagingClient, sid string) error {
	return client.delete("/Services/" + sid)
}

// UpdateIPService updates an existing IP Messaging Service.
func UpdateIPService(client *TwilioIPMessagingClient, sid string, friendlyName string, defaultServiceRoleSid string, defaultChannelRoleSid string,
	typingIndicatorTimeout time.Duration, webhooks Webhooks) (*IPService, error) {

	timeout, err := durationToISO8601(typingIndicatorTimeout)
	if err != nil {
		return nil, err
	}

	var service *IPService

	params := url.Values{}
	params.Set("FriendlyName", friendlyName)
	params.Set("DefaultServiceRoleSid", defaultServiceRoleSid)
	params.Set("DefaultChannelRoleSid", defaultChannelRoleSid)
	params.Set("TypingIndicatonTimeout", timeout)
	for k, v := range webhooks {
		params.Set(k, v)
	}

	res, err := client.post(params, "/Services/"+sid+".json")

	if err != nil {
		return service, err
	}

	service = new(IPService)
	err = json.Unmarshal(res, service)

	return service, err
}

// ListIPServices returns the first page of services.
func ListIPServices(client *TwilioIPMessagingClient) (*IPServiceList, error) {
	var serviceList *IPServiceList

	body, err := client.get(nil, "/Services.json")

	if err != nil {
		return serviceList, err
	}

	serviceList = new(IPServiceList)
	serviceList.Client = client
	err = json.Unmarshal(body, serviceList)

	return serviceList, err
}

// GetServices returns the current page of services.
func (s *IPServiceList) GetServices() []IPService {
	return s.Services
}

// GetAllServices returns all of the services from all of the pages (from here forward).
func (s *IPServiceList) GetAllServices() ([]IPService, error) {
	services := s.Services
	t := s

	for t.HasNextPage() {
		var err error
		t, err = t.NextPage()
		if err != nil {
			return nil, err
		}
		services = append(services, t.Services...)
	}
	return services, nil
}

// HasNextPage returns whether or not there is a next page of services.
func (s *IPServiceList) HasNextPage() bool {
	return s.Meta.NextPageUri != ""
}

// NextPage returns the next page of services.
func (s *IPServiceList) NextPage() (*IPServiceList, error) {
	if !s.HasNextPage() {
		return nil, Error{"No next page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// HasPreviousPage indicates whether or not there is a previous page of results.
func (s *IPServiceList) HasPreviousPage() bool {
	return s.Meta.PreviousPageUri != ""
}

// PreviousPage returns the previous page of services.
func (s *IPServiceList) PreviousPage() (*IPServiceList, error) {
	if !s.HasPreviousPage() {
		return nil, Error{"No previous page"}
	}

	return s.getPage(s.Meta.NextPageUri)
}

// FirstPage returns the first page of services.
func (s *IPServiceList) FirstPage() (*IPServiceList, error) {
	return s.getPage(s.Meta.FirstPageUri)
}

// LastPage returns the last page of services.
func (s *IPServiceList) LastPage() (*IPServiceList, error) {
	return s.getPage(s.Meta.LastPageUri)
}

func (s *IPServiceList) getPage(uri string) (*IPServiceList, error) {
	var serviceList *IPServiceList

	client := s.Client

	body, err := client.get(nil, uri)

	if err != nil {
		return serviceList, err
	}

	serviceList = new(IPServiceList)
	serviceList.Client = client
	err = json.Unmarshal(body, serviceList)

	return serviceList, err
}
