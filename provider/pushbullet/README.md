# Pushbullet

To configure the Pushbullet provider, you need to specify an access token, which can be generated in your [account settings](https://www.pushbullet.com/#settings/account).

```
providers:
    pushbullet:
        access_token: 'o.AbCdEfGhIjKlMnOpQrStUvWxYz012345'
```

To configure a pushbullet receiver you must specify a list of targets in the form `targetType:targetName` where the format of the name depends on the type. Currently pushing to a device or to subscribers of a channel are supported.

```
receivers:
- name: 'pushbullet'
    provider: 'pushbullet'
    to:
      - device:My Nickname
      - channel:mytag
```