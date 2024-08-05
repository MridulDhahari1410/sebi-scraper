package gatekeeper

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	configs "github.com/angel-one/go-config-client"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type TokenGenerator interface {
	GenerateToken(customClaims AngelOneClaims, expiresInSec time.Duration) (string, error)
	GenerateTokenWithIssueTimeAndExpiryTime(customClaims AngelOneClaims, iat time.Time, exp time.Time) (string, error)
}

type DefaultTokenGenerator struct {
	keyProvider KeyProvider
}

func NewTokenGenerator(configClient configs.Client, secretName string) (TokenGenerator, error) {
	//accept the config client only and the use private key
	configProvider, err := newConfigClientBasedConfigProvider(configClient, empty, secretName)
	if err != nil {
		return nil, err
	}
	provider := newConfigAsymmetricKeyProvider(configProvider)
	return &DefaultTokenGenerator{
		keyProvider: provider,
	}, nil
}

func NewTokenGeneratorWithKeyProvider(keyProvider KeyProvider) (TokenGenerator, error) {
	return &DefaultTokenGenerator{
		keyProvider: keyProvider,
	}, nil
}

func (d *DefaultTokenGenerator) GenerateToken(customClaims AngelOneClaims, expiresIn time.Duration) (string, error) {
	if d.invalidUserType(customClaims) || customClaims.Issuer == empty || customClaims.KeyId == empty {
		return empty, errors.New("mandatory params are empty")
	}
	key, err := d.getPrivateKey(customClaims)
	if err != nil {
		return empty, err
	}
	currentTime := time.Now()
	iat := currentTime.Add(-time.Second * issuedTimeTolerance)
	exp := currentTime.Add(expiresIn)

	return d.getToken(customClaims, exp, iat, key)
}

func (d *DefaultTokenGenerator) getToken(customClaims AngelOneClaims, exp time.Time, iat time.Time, key *rsa.PrivateKey) (string, error) {
	claims := TokenClaims{
		CustomClaims: CustomClaims{
			UserType:      customClaims.UserType,
			MobileNo:      customClaims.MobileNo,
			TokenType:     customClaims.TokenType,
			Scope:         customClaims.Scope,
			DataCenter:    customClaims.DataCenter,
			GMId:          customClaims.GMId,
			OmneManagerID: customClaims.GMId,
			Source:        customClaims.Source,
			DeviceId:      customClaims.DeviceId,
			KeyId:         customClaims.KeyId,
			Products:      customClaims.Products,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    customClaims.Issuer,
			Subject:   customClaims.Subject,
			Audience:  customClaims.Audience,
			ExpiresAt: jwt.NewNumericDate(exp),
			NotBefore: jwt.NewNumericDate(iat),
			IssuedAt:  jwt.NewNumericDate(iat),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(key)
}

func (d *DefaultTokenGenerator) GenerateTokenWithIssueTimeAndExpiryTime(customClaims AngelOneClaims, iat time.Time, exp time.Time) (string, error) {
	if d.invalidUserType(customClaims) || customClaims.Issuer == empty || customClaims.KeyId == empty {
		return empty, errors.New("mandatory params are empty")
	}

	key, err := d.getPrivateKey(customClaims)
	if err != nil {
		return empty, err
	}
	return d.getToken(customClaims, exp, iat, key)
}

func (d *DefaultTokenGenerator) getPrivateKey(customClaims AngelOneClaims) (*rsa.PrivateKey, error) {
	base64Data, _ := d.keyProvider.GetPrivateKey(customClaims.KeyId)
	actualPrivateKey, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(actualPrivateKey)

	if err != nil {
		return nil, err
	}
	return key, nil
}

func (d *DefaultTokenGenerator) invalidUserType(customClaims AngelOneClaims) bool {
	return customClaims.UserType == empty ||
		(customClaims.UserType != UserTypeApplication &&
			customClaims.UserType != UserTypeAdmin &&
			customClaims.UserType != UserTypeClient)
}
