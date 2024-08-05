package gatekeeper

import (
	"fmt"
)

type configAsymmetricKeyProvider struct {
	configProvider configProvider
}

func newConfigAsymmetricKeyProvider(configProvider configProvider) KeyProvider {
	return &configAsymmetricKeyProvider{
		configProvider: configProvider,
	}
}

func (receiver *configAsymmetricKeyProvider) GetPrivateKey(keyName string) (string, error) {
	return receiver.configProvider.getStringSecret(keyName)
}
func (receiver *configAsymmetricKeyProvider) GetPublicKey(issuer string, keyName string) (string, error) {
	return receiver.configProvider.getStringConfig(fmt.Sprintf("%s.%s.%s.%s", authentication, publicKey, issuer, keyName))
}
