# MessageBird `go-rest-api` upgrading guide
This guide documents breaking changes across major versions and should be taken into account when updating your dependencies.

## `v4.2.1` -> `v5.0.0`

### Package structure
The most obvious change in this version is the new package structure. Nearly all code used to live in the `messagebird` package - `v5` changes this.

All resources (balance, voicemessage et al) are now moved to their own packages.
This also means the functions no longer take the client as a receiver, but as their first parameter.

Let's see some code. We'll fetch the balance for your account.

Before:
```go
client := messagebird.New("YOUR_ACCESS_KEY")
balance, err := client.Balance()

// ...
```

After:
```go
client := messagebird.New("YOUR_ACCESS_KEY")
balance, err := balance.Read(client)

// ...
```

### Naming
#### Functions
Some function names have changed. They are now more consistent across the several resources(/packages).

Typically, the HTTP method for the underlaying API calls determine the function's name. Here's an overview.

| HTTP method   | Function name | Note                                                           |
| :------------ |:------------- | :--------------------------------------------------------------|
| GET           | List()        | Lists a resource, e.g. `/some-resource`                        |
| GET           | Read()        | Retrieves a single resource, e.g. `/some-resource/{resourceId}`|
| POST          | Create()      | -                                                              |

#### Resources
Previously, SMS messages were referred to simply as `messages` (e.g. `client.NewMessage()`) .
The package that takes care of this now is named `sms`. Naming it `message` is a [bad idea](https://blog.golang.org/package-names), and `smsmessage` doesn't look very friendly.

MMS messages are handled by the `mms` package.

Many structs have been renamed too. We used to have `messagebird.MMSMessageParams`, but we don't want any stuttering like `mms.MMSMessageParams`. That's now simply `mms.Params`.

### Recipient
The `Recipient.Recipient` field used to be of type `int`.
On i386 architectures, the `int` type has only 32 bits, meaning MSISDNs overflowed this. [This PR](https://github.com/messagebird/go-rest-api/pull/19) fixes this by changing its type to `int64`.
The struct is mainly used for unmarshalling responses so you likely need not worry about this.

### Error handling
The way errors are handled changes significantly in `v5`.

Before, the errors returned by the MessageBird API would be included in the resulting resource. You would do this (ignoring the change to the package structure, as described above):

```go
balance, err := balance.Read(client)
if err != nil && err == messagebird.ErrResponse {
    for _, mbError := range balance.Errors {
        // Do something with mbError.
    }
}
```

Starting with `v5`, errors from the MessageBird API also implement the `error` interface.
They are no longer defined on the resources, but instead returned as you're used to from most Go code (its last return value).

It looks like this:

```go
balance, err := balance.Read(client)
if err != nil {
    // Print a generic error message.
    log.Printf("%s\n", err)

    // Or inspect individual errors:
    switch errResp := err.(type) {
    case messagebird.ErrorResponse:
        for _, mbError := range errResp.Errors {
            log.Printf("%d: %s in %s", mbError.Code, mbError.Description, mbError.Parameter)

            // Implements `error` as well...
            log.Printf("%s\n", mbError)
        }
    }
}
```

### Other improvements
Although those are all breaking changes, it's worth mentioning this new version brings a number of other improvements. An example is that the client now supports contacts and groups - enjoy!
