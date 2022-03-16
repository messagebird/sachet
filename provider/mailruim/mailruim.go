package mailruim

import (
	botgolang "github.com/mail-ru-im/bot-golang"

	"github.com/messagebird/sachet"
)

type Config struct {
	Token string `yaml:"token"`
	Url   string `yaml:"url"`
}

type MailruIM struct {
	bot *botgolang.Bot
}

func NewMailruIM(config Config) (*MailruIM, error) {
	bot, err := botgolang.NewBot(config.Token, botgolang.BotApiURL(config.Url))
	if err != nil {
		return nil, err
	}

	return &MailruIM{
		bot: bot,
	}, nil
}

func (mr *MailruIM) Send(message sachet.Message) error {
	for _, ChatID := range message.To {
		msg := mr.bot.NewTextMessage(ChatID, message.Text)
		if err := msg.Send(); err != nil {
			// TODO: handle the error
		}
	}
	return nil
}
