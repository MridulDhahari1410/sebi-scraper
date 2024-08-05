package gatekeeper

import "context"

type tokenStorageValidator interface {
	validate(loginClaims *JWTTokenClaims) (bool, error)
	setSessionValidator(keyValidator SessionValidator)
}

func newTokenStorageValidator(provider configProvider, logger func(ctx context.Context, logMessage LogMessage)) (tokenStorageValidator, error) {
	return newCacheBasedTokenStorageValidator(provider, logger)
}
