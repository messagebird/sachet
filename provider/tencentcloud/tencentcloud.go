package tencentcloud

import (
	"encoding/json"
	"fmt"
	"github.com/messagebird/sachet"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
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

func (tencentcloud *TencentCloud) Send(message sachet.Message) error {
	var err error = nil
	switch message.Type {
	case "", "text":
		request := sms.NewSendSmsRequest()
		request.SmsSdkAppid = common.StringPtr(tencentcloud.config.AppId)
		request.Sign = common.StringPtr(tencentcloud.config.SignName)
		request.TemplateParamSet = common.StringPtrs([]string{message.Text})
		request.TemplateID = common.StringPtr(tencentcloud.config.TemplateCode)
		request.PhoneNumberSet = common.StringPtrs(message.To)
		response, err := tencentcloud.client.SendSms(request)
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			fmt.Printf("An API error has returned: %s", err)
			return err
		}
		b, _ := json.Marshal(response.Response)
		fmt.Printf("%s", b)
	}
	return err

}
