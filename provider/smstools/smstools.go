package smstools

import (
	"io/ioutil"

	"github.com/messagebird/sachet"
)

type Config struct {
	OutgoingDir string `yaml:"outgoing_dir"`
}

var _ (sachet.Provider) = (*SmsTools)(nil)

type SmsTools struct {
	Config
}

func NewSmsTools(config Config) *SmsTools {
	SmsTools := &SmsTools{Config: config}
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
		outfile.WriteString("\n")
		_, err = outfile.WriteString(message.Text)
		if err != nil {
			return err
		}
	}
	return nil
}
