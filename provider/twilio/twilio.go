package twilio

import (
	"github.com/carlosdp/twiliogo"
	"github.com/messagebird/sachet"
)

type Config struct {
	AccountSID string `yaml:"account_sid"`
	AuthToken  string `yaml:"auth_token"`
}

var _ (sachet.Provider) = (*Twilio)(nil)

type Twilio struct {
	client twiliogo.Client
}

const (
	accountSIDVaultKey = "TWILIO_ACCOUNT_SID"
	authTokenVaultKey  = "TWILIO_AUTH_TOKEN"
)

func NewTwilio(config Config, secretsProvider sachet.SecureSecretsProvider) (*Twilio, error) {
	var accountSID, authToken string
	if config.AccountSID == accountSIDVaultKey && config.AuthToken == authTokenVaultKey {
		secrets, err := secretsProvider.GetSecrets(accountSIDVaultKey, authTokenVaultKey)
		if err != nil {
			return nil, err
		}
		accountSID, authToken = secrets[accountSIDVaultKey], secrets[authTokenVaultKey]
	} else {
		accountSID, authToken = config.AccountSID, config.AuthToken
	}
	return &Twilio{client: twiliogo.NewClient(accountSID, authToken)}, nil
}

func (tw *Twilio) Send(message sachet.Message) error {
	for _, recipient := range message.To {
		_, err := twiliogo.NewMessage(tw.client, message.From, recipient, twiliogo.Body(message.Text))
		if err != nil {
			return err
		}
	}

	return nil
}
