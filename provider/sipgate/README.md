# Sipgate REST API v2

To configure the Sipgate provider, you need to specify the SMS ID which will be used to lookup
SMS sender name and number. If you do not know the SMS ID, you can find it out via the Sipgate API.

First, obtain your user ID:

```
$ curl -X GET --user <username>:<password> https://api.sipgate.com/v2/authorization/userinfo
```

The relevant user ID is in the form "wNNN" where "NNN" are numerical digits. If your user is not
the main account owner, it may be in the "sub" attribute of the response.

Next, get a list of SMS IDs, replacing the `<wNNN>` with your user ID.

```
$ curl -X GET --user <username>:<password> https://api.sipgate.com/v2/<wNNN>/sms
```

The response should contain an array of one or more items, with IDs in the form "sNNN", where the
"NNN" are numerical digits (often the same as the "wNNN" user ID).

Use the "sNNN" SMS ID as the "from" attribute in the "receivers" section of the configuration file.
