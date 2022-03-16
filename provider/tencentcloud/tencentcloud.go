package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/messagebird/sachet"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

type Config struct {
	SecretId     string `yaml:"secret_id"`
	SecretKey    string `yaml:"secret_key"`
	AppId        string `yaml:"app_id"`
	Region       string `yaml:"region"`
	Endpoint     string `yaml:"endpoint"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
	Truncate     bool   `yaml:"truncate"`
}

type TencentCloud struct {
	client *sms.Client
	config *Config
}

func NewTencentCloud(config Config) (*TencentCloud, error) {
	credential := common.NewCredential(
		config.SecretId,
		config.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.Endpoint = config.Endpoint
	cpf.SignMethod = "HmacSHA1"
	client, _ := sms.NewClient(credential, config.Region, cpf)
	return &TencentCloud{
		client: client,
		config: &config,
	}, nil
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "...\ntoo many text truncated, show detail on alert platform"
	}
	return bnoden
}

func (tencentcloud *TencentCloud) Send(message sachet.Message) error {
	switch message.Type {
	case "", "text":
		request := sms.NewSendSmsRequest()
		request.SmsSdkAppid = common.StringPtr(tencentcloud.config.AppId)
		request.Sign = common.StringPtr(tencentcloud.config.SignName)
		sendText := message.Text
		if tencentcloud.config.Truncate {
			sendText = truncateString(sendText, 400)
		}

		request.TemplateParamSet = common.StringPtrs([]string{sendText})
		request.TemplateID = common.StringPtr(tencentcloud.config.TemplateCode)
		request.PhoneNumberSet = common.StringPtrs(message.To)
		response, err := tencentcloud.client.SendSms(request)

		var errTencentCloudSDKError *tcError.TencentCloudSDKError
		if errors.As(err, &errTencentCloudSDKError) {
			fmt.Printf("An API error has returned: %s", err)
			return err
		}

		b, err := json.Marshal(response.Response)
		if err != nil {
			return err
		}
		fmt.Printf("%s", b)
	}

	return nil
}
