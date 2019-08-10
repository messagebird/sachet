package pushbullet

import (
	"fmt"
	"strings"

	"github.com/messagebird/sachet"
	"github.com/xconstruct/go-pushbullet"
)

const (
	deviceTargetType  = "device"
	channelTargetType = "channel"
)

// Config is the configuration struct for the Pushbullet provider
type Config struct {
	AccessToken string `yaml:"access_token"`
}

// Pushbullet contains the necessary values for the Pushbullet provider
type Pushbullet struct {
	Config
}

// NewPushbullet creates and returns a new Pushbullet struct
func NewPushbullet(config Config) *Pushbullet {
	return &Pushbullet{config}
}

// Send pushes a note to devices registered in configuration
func (c *Pushbullet) Send(message sachet.Message) error {

	for _, recipient := range message.To {

		// create pushbullet client
		pb := pushbullet.New(c.AccessToken)

		// parse recipient
		targetTypeName := strings.Split(recipient, ":")
		if len(targetTypeName) != 2 {
			return fmt.Errorf("cannot parse recipient %s: expecting targetType:targetName", recipient)
		}
		targetType := targetTypeName[0]
		targetName := targetTypeName[1]

		switch targetType {
		case deviceTargetType:
			// retrieve device
			dev, err := pb.Device(targetName)
			if err != nil {
				return err
			}

			// push note
			err = pb.PushNote(dev.Iden, message.From, message.Text)
			if err != nil {
				return err
			}
		case channelTargetType:
			// retrieve subscription
			sub, err := pb.Subscription(targetName)
			if err != nil {
				return err
			}

			// push note
			err = sub.PushNote(message.From, message.Text)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognised target type: %s", targetType)
		}

	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
