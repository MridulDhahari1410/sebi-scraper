package gatekeeper

import (
	"errors"
	"fmt"
)

type commonKeyProvider struct {
	configProvider configProvider
}

func newCommonKeyProvider(configProvider configProvider) KeyProvider {

	keyProvider := &commonKeyProvider{
		configProvider: configProvider,
	}
	return keyProvider
}

func (receiver *commonKeyProvider) GetPrivateKey(keyName string) (string, error) {
	if keyName == nonTradeTokenJwtSigningKey {
		configValue, err := receiver.configProvider.getBoolConfig(enableOldJWTSupportForNonTradeLoginService)
		if configValue {
			return receiver.configProvider.getStringSecret(keyName)
		}
		return "", fmt.Errorf("old non trade jwt token support is disabled or error: %w", err)
	} else if keyName == tradeTokenJwtSigningKey {
		configValue, err := receiver.configProvider.getBoolConfig(enableOldJWTSupportForTradeLoginService)
		if configValue {
			return receiver.configProvider.getStringSecret(keyName)
		}
		return "", fmt.Errorf("old trade jwt token support is disabled or error: %w", err)
	} else if keyName == s2sTokenJwtSigningKey {
		configValue, err := receiver.configProvider.getBoolConfig(enableS2SSupport)
		if configValue {
			return receiver.configProvider.getStringSecret(keyName)
		}
		return "", fmt.Errorf("old s2s jwt token support is disabled or error: %w", err)
	}
	return "", errors.New("common key configProvider called for incorrect secret key")
}
func (receiver *commonKeyProvider) GetPublicKey(issuer string, keyName string) (string, error) {
	key, err := receiver.configProvider.getStringConfig(fmt.Sprintf("%s.%s.%s.%s", authentication, publicKey, issuer, keyName))
	if err != nil {
		receiver.configProvider.refreshConfigs()
		return receiver.configProvider.getStringConfig(fmt.Sprintf("%s.%s.%s.%s", authentication, publicKey, issuer, keyName))
	}
	return key, err
}
