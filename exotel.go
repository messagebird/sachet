package main

import (
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
	"time"
)

//ExotelBaseURL is the base url for exotel api
const ExotelBaseURL = "https://twilix.exotel.in"
//ExotelRetryInterval is the retry interval for failed requests
const ExotelRetryInterval = 3 //in Seconds

//Exotel is the exte Exotel
type Exotel struct {
	AccountSid string
	Token      string
	requester  *gorequest.SuperAgent
}

//NewExotel creates a new
func NewExotel(sid ,token string) *Exotel {
	Exotel := &Exotel{requester: gorequest.New().SetBasicAuth(sid, token), AccountSid: sid, Token: token}
	return Exotel
}

//SetHeaders Sets needed headers
func (c *Exotel) setHeaders() *Exotel {
	c.requester.
		Set("User-Agent", "sachet-v1").
		Set("Accept-language", "es").
		Set("Content-Type", "application/x-www-form-urlencoded").
		Set("accept", "application/json")
	return c
}

//sendSMS sends sms  ,obviously  
//This API is rate limited to 200 SMSes per minute. Once this limit has been crossed, your requests will be rejected with an HTTP 429 'Too Many Requests' code.
func (c *Exotel) sendSMS(from string, to []string,body string) ( err error) {
	c.requester.Post(fmt.Sprintf("%s/v1/Accounts/%s/Smss/send.json", ExotelBaseURL, c.AccountSid)).
		Param("From", from).
		Param("Body", body).
        Retry(3,time.Second * ExotelRetryInterval)
    //adding to numbers
    for _,number :=  range to {
        c.requester.Param("To", number)
    }
	c.setHeaders()
	resp, body, errs := c.requester.End()
	if errs != nil && len(errs) > 0 {
		return  fmt.Errorf("Failed making sms request using exotel ,  ERRORS : %+v", errs)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("Failed sending sms:Reason: %s , StatusCode : %d", body, resp.StatusCode)
}


//Send send sms to n number of people using bulk sms api  
func (*Exotel) Send(message Message) (err error) {
   client := NewExotel(config.Providers.Exotel.AccountSID, config.Providers.Exotel.Token)
   return client.sendSMS(message.From, message.To, message.Text)
}