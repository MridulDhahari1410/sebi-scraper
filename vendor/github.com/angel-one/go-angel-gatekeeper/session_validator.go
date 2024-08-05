package gatekeeper

type SessionValidator interface {
	GetSessionKey(claims *JWTTokenClaims) string
	Validate(val interface{}, keyReadErr error, claims *JWTTokenClaims) bool
}
