package telegram

import (
	"strconv"
	"strings"

	"github.com/messagebird/sachet"
	"gopkg.in/telegram-bot-api.v4"
)

type TelegramConfig struct {
	Token     string `yaml:"token"`
	ParseMode string `yaml:"parse_mode"`
}

type Telegram struct {
	bot       *tgbotapi.BotAPI
	ParseMode string
}

func NewTelegram(config TelegramConfig) (*Telegram, error) {
	ParseMode := strings.ToLower(config.ParseMode)
	switch ParseMode {
	case "md", "markdown":
		ParseMode = "Markdown"
	case "html", "h":
		ParseMode = "HTML"
	default:
		ParseMode = ""
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		bot:       bot,
		ParseMode: ParseMode,
	}, nil
}

func (tg *Telegram) Send(message sachet.Message) error {
	for _, sChatID := range message.To {
		chatID, err := strconv.ParseInt(sChatID, 10, 64)
		if err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(chatID, message.Text)
		msg.ParseMode = tg.ParseMode
		_, err = tg.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
