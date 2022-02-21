[comment]: <> (HEAD)
# TextMagic Go SDK

This library provides you with an easy way of sending SMS and receiving replies by integrating the TextMagic SMS Gateway into your Go application.

## What Is TextMagic?
TextMagic’s application programming interface (API) provides the communication link between your application and TextMagic’s SMS Gateway, allowing you to send and receive text messages and to check the delivery status of text messages you’ve already sent.

[comment]: <> (/HEAD)
## Installation

With go.mod:
```bash
go get -u github.com/textmagic/textmagic-rest-go-v2/v2@v2.0.1816
```

Without go.mod:
```bash
go get -u github.com/textmagic/textmagic-rest-go-v2
```

## Dependencies:

- https://github.com/antihax/optional - this library is used to provide optional argument types for API calls. See the `GetAllOutboundMessages` call example below.

## Usage Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/antihax/optional"
    // If you're using go.mod use line below to import our module
    // 	tm "github.com/textmagic/textmagic-rest-go-v2/v2"
    tm "github.com/textmagic/textmagic-rest-go-v2"
    "log"
)

func main() {

    cfg := tm.NewConfiguration()
    cfg.BasePath = "https://rest.textmagic.com"
    client := tm.NewAPIClient(cfg)

    // put your Username and API Key from https://my.textmagic.com/online/api/rest-api/keys page.
    auth := context.WithValue(context.Background(), tm.ContextBasicAuth, tm.BasicAuth{
        UserName: "YOUR_USERNAME",
        Password: "YOUR_API_KEY",
    })

    // Simple ping request example
    pingResponse, _, err := client.TextMagicApi.Ping(auth)

    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println(pingResponse.Ping)
    }

    // Send a new message request example
    sendMessageResponse, _, err := client.TextMagicApi.SendMessage(auth, tm.SendMessageInputObject{
        Text: "I love TextMagic!",
        Phones: "+19998887766",
    }, &tm.SendMessageOpts{})

    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println(sendMessageResponse.Id)
    }

    // Get all outgoing messages request example
    getAllOutboundMessageResponse, _, err := client.TextMagicApi.GetAllOutboundMessages(auth, &tm.GetAllOutboundMessagesOpts{
        Page: optional.NewInt32(1),
        Limit: optional.NewInt32(250),
    })

    if err != nil {
        log.Fatal(err)
    } else {
        fmt.Println(getAllOutboundMessageResponse.Resources[0].Id)
    }
}
```

## Limitations
Due to the issue at https://github.com/swagger-api/swagger-codegen/issues/7311, the current version of Go SDK does not support any file uploading API calls.

[comment]: <> (FOOTER)
## License
The library is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).

[comment]: <> (/FOOTER)
