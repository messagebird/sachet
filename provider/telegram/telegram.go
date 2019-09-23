package telegram

import (
	"strconv"

	"github.com/messagebird/sachet"
	"gopkg.in/telegram-bot-api.v4"
	"math"
)

const (
	messageSize = 4096 
)

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type Telegram struct {
	bot *tgbotapi.BotAPI
}



func NewTelegram(config TelegramConfig) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, err
	}

	return &Telegram{
		bot: bot,
	}, nil
}

func splitByLength(text string, messageSize int) []string {
    runeStr := []rune(text)
    strLength := len(runeStr)
    numChunks := int(math.Ceil(float64(strLength) / float64(messageSize)))

    splited := make([]string, numChunks)

    start,stop := 0,0
    for i := 0; i < numChunks; i += 1 {
        start = i * messageSize
        stop = start + messageSize
        if stop > strLength {
            stop = strLength
        }
        splited[i] = string(runeStr[start : stop])
    }
    return splited
}


func (tg *Telegram) Send(message sachet.Message) error {
	for _, sChatID := range message.To {
		chatID, err := strconv.ParseInt(sChatID, 10, 64)
		if err != nil {
			return err
		}
		
		sendedMsg := splitByLength(message.Text, messageSize)

		for _, v := range sendedMsg {
			msg := tgbotapi.NewMessage(chatID, v)
			_, err = tg.bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
