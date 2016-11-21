# gonexmo [![GoDoc](https://godoc.org/github.com/njern/gonexmo?status.png)](https://godoc.org/gopkg.in/njern/gonexmo.v1)

gonexmo is a [Go](http://golang.org/) library tailored for sending SMS's with
[Nexmo](https://www.nexmo.com/).


## Installation

Assuming you have a working Go environment, installation is simple:

    go get "gopkg.in/njern/gonexmo.v1"
    
You can take a look at the documentation locally with:

	godoc github.com/njern/gonexmo

The included tests in `gonexmo_test.go` also illustrate usage of the package.

**Note:** You must enter valid API credentials and a valid phone number in
`gonexmo_test.go` or the tests will fail! I didn't feel like draining my own
Nexmo account or receiving thousands of test SMS's - sorry :)


## Usage
    import "gopkg.in/njern/gonexmo.v1"

    nexmoClient, _ := nexmo.NewClientFromAPI("API_KEY_GOES_HERE", "API_SECRET_GOES_HERE")

    // Test if it works by retrieving your account balance
    balance, err := nexmoClient.Account.GetBalance()

    // Send an SMS
    // See https://docs.nexmo.com/index.php/sms-api/send-message for details.
	message := &nexmo.SMSMessage{
		From:           "go-nexmo",
        To:              "00358123412345",
		Type:            nexmo.Text,
		Text:            "Gonexmo test SMS message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           nexmo.Standard,
	}

	messageResponse, err := nexmoClient.SMS.Send(message)

## Receiving inbound messages

    import (
        "gopkg.in/njern/gonexmo.v1"
        "log"
        "net/http"
    )

    func main() {
        messages := make(chan *nexmo.RecvdMessage)
        h := nexmo.NewMessageHandler(messages,false)

        go func() {
            for {
                msg := <-messages
                log.Printf("%v\n",msg)
            }
        }()

        // Set your Nexmo callback url to http://<domain or ip>:8080/get/
        http.HandleFunc("/get/", h)
        if err := http.ListenAndServe(":8080", nil); err != nil {
            log.Fatal("ListenAndServe: ", err)
        }

    }


## Future plans

* Implement the rest of the Nexmo API

## How can you help?

* Let me know if you're using gonexmo by dropping me a line at
  [github user name] at walkbase.com
* Let me know about any bugs / annoyances the same way
