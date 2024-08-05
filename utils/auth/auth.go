package auth

import (
	"fmt"

	"sebi-scrapper/utils/configs"

	"sebi-scrapper/constants"

	gatekeeper "github.com/angel-one/go-angel-gatekeeper"
)

type gatekeeperKeyProvider struct{}

var authenticator gatekeeper.Authenticator
var keyProvider *gatekeeperKeyProvider

// Init is used to initialise s2s s2sAuth.
func Init() error {
	var err error
	keyProvider = &gatekeeperKeyProvider{}
	configOptions := configs.GetClientOptions()
	authenticator, err = gatekeeper.NewDefaultAuthenticatorWithKeyProvider(keyProvider, gatekeeper.Options{
		Provider: configOptions.Provider,
		Params:   configOptions.Params,
	}, constants.ApplicationName)
	if err != nil {
		return err
	}
	authenticator.SetCtxControl(true)
	return initS2SAuth()
}

// GetPrivateKey is used to provide the private key.
func (g gatekeeperKeyProvider) GetPrivateKey(_ string) (string, error) {
	return constants.Empty, nil
}

// GetPublicKey is used to provide the public key.
// This is only used for s2s.
func (g gatekeeperKeyProvider) GetPublicKey(issuer string, keyName string) (string, error) {
	return configs.Get().GetString(constants.S2SAuthConfig, fmt.Sprintf("%s.%s.%s", issuer,
		constants.S2SAuthPublicKeysConfigKey, keyName))
}
