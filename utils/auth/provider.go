package auth

import "context"

// KeyProvider is the provider for s2s s2sAuth.
type KeyProvider struct {
	service string
}

// NewKeyProvider is used to get key provider.
func NewKeyProvider(service string) *KeyProvider {
	return &KeyProvider{service: service}
}

// GenerateS2SToken is used to generate s2s token.
func (k *KeyProvider) GenerateS2SToken(_ context.Context) (string, error) {
	return GenerateS2SToken(k.service)
}
