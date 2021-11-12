package telegram

import (
	"github.com/messagebird/sachet"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
)

type Config struct {
	Token     string `yaml:"token"`
	ParseMode string `yaml:"parse_mode"`
}

type Telegram struct {
	bot       *tgbotapi.BotAPI
	config    *Config
}

func NewTelegram(config Config) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		bot:    bot,
		config: &config,
	}, nil
}

func (tg *Telegram) Send(message sachet.Message) error {
	for _, sChatID := range message.To {
		chatID, err := strconv.ParseInt(sChatID, 10, 64)
		if err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(chatID, message.Text)
		msg.ParseMode = tg.config.ParseMode

		_, err = tg.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
