package textmagic

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/messagebird/sachet"
	textmagic "github.com/textmagic/textmagic-rest-go-v2/v2"
)

type Config struct {
	Username string `yaml:"username"`
	APIKey   string `yaml:"api_key"`
}

type TextMagic struct {
	client *textmagic.APIClient
	auth   context.Context
}

func NewTextMagic(config Config) *TextMagic {
	cfg := textmagic.NewConfiguration()
	cfg.BasePath = "https://rest.textmagic.com"
	client := textmagic.NewAPIClient(cfg)
	auth := context.WithValue(context.Background(), textmagic.ContextBasicAuth, textmagic.BasicAuth{
		UserName: config.Username,
		Password: config.APIKey,
	})
	return &TextMagic{
		client: client,
		auth:   auth,
	}
}

func (tm *TextMagic) Send(message sachet.Message) error {
	var err error = nil
	switch message.Type {
	case "", "text":
		joinedPhones := strings.Join(message.To[:], ",")
		_, _, err = tm.client.TextMagicApi.SendMessage(tm.auth, textmagic.SendMessageInputObject{
			Text:   message.Text,
			Phones: joinedPhones,
			From:   message.From,
		})
	default:
		return fmt.Errorf("unknown message type %s", message.Type)
	}
	return err
}
