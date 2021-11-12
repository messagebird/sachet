<img src="https://github.com/mail-ru-im/bot-python/blob/master/logo.png" width="100" height="100">

# Golang interface for Mail.ru Instant Messengers bot API
![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
[![CircleCI](https://circleci.com/gh/mail-ru-im/bot-golang.svg?style=svg)](https://circleci.com/gh/mail-ru-im/bot-golang)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/mail-ru-im/bot-golang)

 - *Brand new Bot API!*

 - *Zero-configuration library*

 - *Simple and clear interface*

## API specification:
### [<img src="https://icq.com/cached/img/landing/icon_and_192.png" width="15"> ICQ New ](https://icq.com/botapi/)
### [<img src="https://is3-ssl.mzstatic.com/image/thumb/Purple123/v4/e8/4f/1b/e84f1b57-206f-7750-ac5a-27f93ff4a0d8/icons-bundle.png/460x0w.png" width="16"> Myteam ](https://myteam.mail.ru/botapi/)

### [<img src="https://agent.mail.ru/img/agent2014/common/2x/button_logo.png" width="16"> Agent Mail.ru](https://agent.mail.ru/botapi/) 

## Install
```bash
go get github.com/mail-ru-im/bot-golang
```

## Usage

Create your own bot by sending the /newbot command to Metabot and follow the instructions.

Note a bot can only reply after the user has added it to his contacts list, or if the user was the first to start a dialogue.

### Create your bot

```go
package main

import "github.com/mail-ru-im/bot-golang"

func main() {
    bot, err := botgolang.NewBot(BOT_TOKEN)
    if err != nil {
        log.Println("wrong token")
    }

    message := bot.NewTextMessage("awesomechat@agent.chat", "text")
    message.Send()
}
```

### Send and edit messages

You can create, edit and reply to messages like a piece of cake.

```go
message := bot.NewTextMessage("awesomechat@agent.chat", "text")
message.Send()

fmt.Println(message.MsgID)

message.Text = "new text"

message.Edit()
// AWESOME!

message.Reply("hey, what did you write before???")
// SO HOT!
```

### Subscribe events

Get all updates from the channel. Use context for cancellation.

```go
ctx, finish := context.WithCancel(context.Background())
updates := bot.GetUpdatesChannel(ctx)
for update := range updates {
	// your awesome logic here
}
```

### Passing options

You don't need this.
But if you do, you can override bot's API URL:

```go
bot := botgolang.NewBot(BOT_TOKEN, botgolang.BotApiURL("https://agent.mail.ru/bot/v1"))
```
And debug all api requests and responses:

```go
bot := botgolang.NewBot(BOT_TOKEN, botgolang.BotDebug(true))
```
