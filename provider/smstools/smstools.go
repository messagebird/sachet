package smstools

import (
	"io/ioutil"

	"github.com/messagebird/sachet"
)

type SmsToolsConfig struct {
	OutgoingDir string `yaml:"outgoing_dir"`
}

type SmsTools struct {
	SmsToolsConfig
}

func NewSmsTools(config SmsToolsConfig) *SmsTools {
	SmsTools := &SmsTools{SmsToolsConfig: config}
	return SmsTools
}

func (smst *SmsTools) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		outfile, err := ioutil.TempFile(smst.OutgoingDir, "sachet-")
		if err != nil {
			return err
		}
		defer outfile.Close()
		outfile.WriteString("To: " + recipient + "\n")
		_, err = outfile.WriteString(message.Text)
		if err != nil {
			return err
		}
	}
	return nil
}
