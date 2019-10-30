# MessageBird
## Provider
To configure the MessageBird provider, you need to specify an access key.

```yaml
providers:
    messagebird:
        access_key: 'live_qKwVZ02ULV70GqabBYxLU8d5r'
        debug: true
        gateway: 240
        language: en-us
        voice: female
        repeat: 2
```

The MessageBird provider supports text and voice messages.
For text messages you can configure the following parameters (see https://developers.messagebird.com/api/sms-messaging):
 * `gateway`: an integer specifying the SMS gateway

For voice messages you can configure the following parameters (see https://developers.messagebird.com/api/voice-messaging/):
 * `language`: specifies the language in which the message needs to be read to the recipient
 * `voice`: specifies the voice in which the message needs to be read to the recipient
 * `repeat`: specifies the number of times the message needs to be repeated

## Receivers
To configure a MessageBird receiver you must specify a list of targets. By default text messages are send, but you can specify a voice message by setting `type: voice` for a receiver. The `type` defaults to `text`. The `from` field is not used for voice messages.

```yaml
receivers:
- name: 'team1'
    provider: messagebird
    to:
      - '+919742033616'
      - '+919742033617'
    type: voice
- name: 'team2'
    provider: messagebird
    to:
      - '+919742033616'
      - '+919742033617'
    from: '08039591643'
    type: text
```
