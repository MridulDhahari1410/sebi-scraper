package auth

import (
	"fmt"
	"time"

	"sebi-scrapper/constants"
	"sebi-scrapper/utils/configs"

	gatekeeper "github.com/angel-one/go-angel-gatekeeper"
)

type s2sConfig struct {
	targetApplicationNames []string
	issuer                 string
	keyID                  string
	subject                string
	expiry                 time.Duration
}

type s2sAuth struct {
	cfg       *s2sConfig
	generator gatekeeper.TokenGenerator
}

type s2sAuthTokenProvider struct {
	privateKey string
}

var s2sAuthTokenGenerator map[string]*s2sAuth

// GenerateS2SToken is used to generate token.
func GenerateS2SToken(service string) (string, error) {
	if g, ok := s2sAuthTokenGenerator[service]; ok {
		return g.generator.GenerateToken(gatekeeper.AngelOneClaims{
			UserType:  gatekeeper.UserTypeApplication,
			Issuer:    g.cfg.issuer,
			Audience:  g.cfg.targetApplicationNames,
			KeyId:     g.cfg.keyID,
			TokenType: gatekeeper.TradeAccessToken,
			Subject:   g.cfg.subject,
		}, g.cfg.expiry)
	}
	return constants.Empty, constants.ErrInvalidS2SAuthService
}

func initS2SAuth() error {
	s2sAuthTokenGenerator = make(map[string]*s2sAuth)
	services := configs.Get().GetStringSliceD(constants.S2SAuthConfig,
		constants.S2SAuthSupportedServicesConfigKey, nil)
	for _, service := range services {
		t, err := gatekeeper.NewTokenGeneratorWithKeyProvider(
			&s2sAuthTokenProvider{privateKey: configs.Get().GetStringWithEnvD(constants.S2SAuthConfig,
				fmt.Sprintf("%s.%s", service, constants.S2SAuthPrivateKeyConfigKey), constants.Empty)})
		if err != nil {
			return err
		}
		s2sAuthTokenGenerator[service] = &s2sAuth{
			cfg: &s2sConfig{
				targetApplicationNames: configs.Get().GetStringSliceD(constants.S2SAuthConfig,
					fmt.Sprintf("%s.%s", service, constants.S2SAuthTargetApplicationNamesConfigKey), nil),
				issuer: configs.Get().GetStringD(constants.S2SAuthConfig,
					fmt.Sprintf("%s.%s", service, constants.S2SAuthIssuerConfigKey), constants.Empty),
				keyID: configs.Get().GetStringD(constants.S2SAuthConfig,
					fmt.Sprintf("%s.%s", service, constants.S2SAuthKeyIDConfigKey), constants.Empty),
				subject: configs.Get().GetStringD(constants.S2SAuthConfig,
					fmt.Sprintf("%s.%s", service, constants.S2SAuthSubjectConfigKey), constants.Empty),
				expiry: time.Second * time.Duration(configs.Get().GetIntD(constants.S2SAuthConfig,
					fmt.Sprintf("%s.%s", service, constants.S2SAuthExpiryInSecondsConfigKey), 0)),
			},
			generator: t,
		}
	}
	return nil
}

// GetPrivateKey is used to provide the private key.
func (a *s2sAuthTokenProvider) GetPrivateKey(_ string) (string, error) {
	return a.privateKey, nil
}

// GetPublicKey is used to get the public key.
func (a *s2sAuthTokenProvider) GetPublicKey(_ string, _ string) (string, error) {
	return constants.Empty, nil
}
