package main

import (
	"github.com/messagebird/sachet/provider/esendex"
	"io/ioutil"

	"github.com/messagebird/sachet/provider/aliyun"
	"github.com/messagebird/sachet/provider/aspsms"
	"github.com/messagebird/sachet/provider/cm"
	"github.com/messagebird/sachet/provider/exotel"
	"github.com/messagebird/sachet/provider/freemobile"
	"github.com/messagebird/sachet/provider/infobip"
	"github.com/messagebird/sachet/provider/kannel"
	"github.com/messagebird/sachet/provider/kavenegar"
	"github.com/messagebird/sachet/provider/mailruim"
	"github.com/messagebird/sachet/provider/mediaburst"
	"github.com/messagebird/sachet/provider/messagebird"
	"github.com/messagebird/sachet/provider/nexmo"
	"github.com/messagebird/sachet/provider/nowsms"
	"github.com/messagebird/sachet/provider/otc"
	"github.com/messagebird/sachet/provider/ovh"
	"github.com/messagebird/sachet/provider/pushbullet"
	"github.com/messagebird/sachet/provider/sap"
	"github.com/messagebird/sachet/provider/sipgate"
	"github.com/messagebird/sachet/provider/sms77"
	"github.com/messagebird/sachet/provider/smsc"
	"github.com/messagebird/sachet/provider/telegram"
	"github.com/messagebird/sachet/provider/tencentcloud"
	"github.com/messagebird/sachet/provider/turbosms"
	"github.com/messagebird/sachet/provider/twilio"

	"github.com/prometheus/alertmanager/template"
	"gopkg.in/yaml.v2"
)

type ReceiverConf struct {
	Name     string
	Provider string
	To       []string
	From     string
	Text     string
	Type     string
}

var config struct {
	Providers struct {
		MessageBird  messagebird.Config
		Nexmo        nexmo.Config
		Twilio       twilio.Config
		Infobip      infobip.Config
		Kannel       kannel.Config
		KaveNegar    kavenegar.Config
		Exotel       exotel.Config
		CM           cm.Config
		MailruIM     mailruim.Config
		Telegram     telegram.Config
		Turbosms     turbosms.Config
		Smsc         smsc.Config
		OTC          otc.Config
		MediaBurst   mediaburst.Config
		FreeMobile   freemobile.Config
		AspSms       aspsms.Config
		Sipgate      sipgate.Config
		Pushbullet   pushbullet.Config
		NowSms       nowsms.Config
		Aliyun       aliyun.Config
		OVH          ovh.Config
		TencentCloud tencentcloud.Config
		Sap          sap.Config
		Esendex      esendex.Config
		Sms77        sms77.Config
	}

	Receivers []ReceiverConf
	Templates []string
}
var tmpl *template.Template

// LoadConfig loads the specified YAML configuration file.
func LoadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	tmpl, err = template.FromGlobs(config.Templates...)
	return err
}
