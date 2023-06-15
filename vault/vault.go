package vault

import (
	"context"
	"fmt"

	vaultApi "github.com/hashicorp/vault/api"
	vaultAuth "github.com/hashicorp/vault/api/auth/kubernetes"
)

type Config struct {
	Address                string
	ServiceAccountTokePath string
	Role                   string
	KVv2MountPath          string
	KVv2SecretPath         string
}

type KVv2SecretsProvider struct {
	Config Config
}

func (p *KVv2SecretsProvider) GetSecrets(secrets ...string) (map[string]string, error) {
	client, err := p.getVaultApiClient()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, secretKey := range secrets {
		secret, err := client.KVv2(p.Config.KVv2MountPath).Get(context.Background(), p.Config.KVv2SecretPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read secret: %w", err)
		}
		value, ok := secret.Data[secretKey]
		if !ok {
			return nil, fmt.Errorf("unable to get secret value for key: %s", secretKey)
		}
		secretValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("value type assertion failed: %T %#v", value, value)
		}
		result[secretKey] = secretValue
	}
	return result, nil
}

func (p *KVv2SecretsProvider) getVaultApiClient() (*vaultApi.Client, error) {
	vaultDefaultConfig := vaultApi.DefaultConfig()
	vaultDefaultConfig.Address = p.Config.Address

	client, err := vaultApi.NewClient(vaultDefaultConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}
	k8sAuth, err := vaultAuth.NewKubernetesAuth(p.Config.Role,
		vaultAuth.WithServiceAccountTokenPath(p.Config.ServiceAccountTokePath))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}
	return client, nil
}
