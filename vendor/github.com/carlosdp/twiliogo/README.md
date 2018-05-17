# twilio-go
The unofficial Go helper library for [Twilio](http://twilio.com).

# Installation

``` bash
go get github.com/carlosdp/twilio-go
```

# Usage

## Send a Text

``` go
package main

import (
  "fmt"
  "github.com/carlosdp/twilio-go"
)

func main() {
  client := twilio.NewClient("<ACCOUNT_SID", "<AUTH_TOKEN>")

  message, err := twilio.SendMessage(client, "3334445555", "2223334444", "Hello World!")

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Message sent!")
  }
}
```

## Make a Call

``` go
package main

import (
  "fmt"
  "github.com/carlosdp/twilio-go"
)

func main() {
  client := twilio.NewClient("<ACCOUNT_SID>", "<AUTH_TOKEN>")

  call, err := twilio.MakeCall(client, "8883332222", "3334443333")

  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Call Queued!")
  }
}
```

## Implemented Resources
- Calls
- Messages
- IncomingPhoneNumbers (partial)

## Run Tests
Tests can be run using `go test`, as with most golang projects. This project also contains integration tests (where they can be done non-destructively using the API or the working Test Credential endpoints).

These integration tests can be run by providing the necessary environment variables to access the API, as in this Makefile:

```makefile
test:
	@export API_KEY="<API KEY>";\
	export API_TOKEN="<API TOKEN>";\
	export TEST_KEY="<TEST KEY>";\
	export TEST_TOKEN="<TEST TOKEN>";\
	export TEST_FROM_NUMBER="<DEFAULT TEST CRED NUMBER>";\
	export FROM_NUMBER="<TEST FROM NUMBER>";\
	export TO_NUMBER="<TEST TO NUMBER>";\
	go test -v
```

## Contributing
This is a side project meant to allow for quick adoption of the Twilio API for those programming web applications with it in Go. Feel free to submit pull requests so that we can cover all of the features the Twilio API has to offer!

## To Do
Here are a few things that the project needs in order to reach v1.0:

1. Complete test coverage. Right now, tests cover the bare minimum of usage for each feature implemented.
2. Complete IncomingPhoneNumber functionality.
3. Implement the following resources:
  - AvailablePhoneNumbers
  - OutgoingCallerIds
  - Applications
  - ConnectApps
  - AuthorizedConnectApps
  - Conferences
  - Queues
  - Short Codes
  - Recordings
  - Transcriptions
  - Notifications
  - SIP Domains
  - IpAccessControlLists
  - CredentialLists
  - Usage Records
  - Usage Triggers

## License
This project is licensed under the [MIT License](http://opensource.org/licenses/MIT)

