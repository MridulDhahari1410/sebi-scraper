package gatekeeper

type TypeValidator interface {
	validate(input *JWTTokenClaims) bool
}

type TokenTypeValidator struct {
	UserTypes  []string
	TokenTypes []string
}

func (t TokenTypeValidator) validate(claims *JWTTokenClaims) bool {
	if !isExists(t.UserTypes, claims.UserType) {
		return false
	}
	return isExists(t.TokenTypes, claims.TokenType)
}

type StorageValidator struct {
}

func (s StorageValidator) validate(input *JWTTokenClaims) bool {
	return true
}

func isExists(slice []string, stringToCheck string) bool {
	if slice == nil || len(slice) == 0 {
		return true
	}
	for _, stringInSlice := range slice {
		if stringInSlice == stringToCheck {
			return true
		}
	}
	return false
}
