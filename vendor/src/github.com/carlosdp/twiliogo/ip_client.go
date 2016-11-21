package twiliogo

// Constants for the IP Messaging service.
const (
	IP_MESSAGING_ROOT     = "https://ip-messaging.twilio.com"
	IP_MESSAGING_VERSION  = "v1"
	IP_MESSAGING_ROOT_URL = IP_MESSAGING_ROOT + "/" + IP_MESSAGING_VERSION
)

// TwilioIPMessagingClient is used for accessing the Twilio IP Messaging API.
type TwilioIPMessagingClient struct {
	TwilioClient
}

var _ Client = &TwilioIPMessagingClient{}

// NewIPMessagingClient creates a new Twilio IP Messaging client.
func NewIPMessagingClient(accountSid, authToken string) *TwilioIPMessagingClient {
	rootUrl := IP_MESSAGING_ROOT + "/" + IP_MESSAGING_VERSION
	return &TwilioIPMessagingClient{TwilioClient{accountSid, authToken, rootUrl}}
}
