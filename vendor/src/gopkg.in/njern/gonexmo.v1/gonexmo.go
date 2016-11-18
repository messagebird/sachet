/*
Package nexmo implements a simple client library for accessing the Nexmo API.

Usage is simple. Create a nexmo.Client instance with NewClientFromAPI(),
providing your API key and API secret. Compose a new Message and then call
Client.SMS.Send(Message). The API will return a MessageResponse which you can
use to see if your message went through, how much it cost, etc.
*/
package nexmo

const (
	apiRoot   = "https://rest.nexmo.com"
	apiRootv2 = "https://api.nexmo.com"
)
