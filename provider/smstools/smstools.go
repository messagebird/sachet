package smstools

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path/filepath"
	"time"

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

func (smst *SmsTools) createFile(recipient, message string) (*os.File, error) {
	h := fnv.New32()
	if _, err := io.WriteString(h, recipient+message); err != nil {
		return nil, err
	}
	if err := binary.Write(h, binary.LittleEndian, time.Now().UnixNano()); err != nil {
		return nil, err
	}

	hsum := h.Sum32()
	prefix := filepath.Join(smst.OutgoingDir, "sachet-")

	try := 0
	for {
		name := fmt.Sprintf("%s%d", prefix, hsum)
		f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
		if os.IsExist(err) {
			if try++; try < 10000 {
				hsum = hsum + 1
				continue
			}
			return nil, fmt.Errorf("unable to create conflict-free filename")
		}
		return f, err
	}
}

func (smst *SmsTools) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		outfile, err := smst.createFile(recipient, message.Text)
		if err != nil {
			return err
		}
		defer outfile.Close()
		outfile.WriteString("To: " + recipient + "\n")
		outfile.WriteString("Alphabet: UTF-8\n")
		outfile.WriteString("\n")
		_, err = outfile.WriteString(message.Text)
		if err != nil {
			return err
		}
	}
	return nil
}
