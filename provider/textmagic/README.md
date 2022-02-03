# TextMagic
## Provider
To configure the TextMagic provider, you need to specify a username and an API key.

```yaml
providers:
  textmagic:
    username: 'donaldduck'
    api_key: 'JNV7NgCKNzQVXNOEpZxwU4c9blGEfF'
```

You can find your username and API key by going to https://my.textmagic.com/online/api/rest-api/keys.

The TextMagic provider supports SMS text messages only.

A TextMagic trial account can be created at https://www.textmagic.com/.

You can test the TextMagic API without spending SMS credits by using a mobile number that
beings with 999. For example, 999742033616, 999742033617, and so on.

## Receivers
To configure a TextMagic receiver you must specify a list of targets. The `from` field is optional.

```yaml
receivers:
  - name: 'team1'
    provider: textmagic
    to:
      - '+999742033616'
      - '+999742033617'
  - name: 'team2'
    provider: textmagic
    to:
      - '+999742033616'
      - '+999742033617'
    from: '08039591643'
```
