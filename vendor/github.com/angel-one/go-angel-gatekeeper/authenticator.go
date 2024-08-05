package gatekeeper

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Authenticator interface {
	Authenticate(opt ...TypeValidator) gin.HandlerFunc                          //gin middleware
	ValidateToken(ctx *gin.Context, opt ...TypeValidator) (*TokenClaims, error) //direct function
	SetCtxControl(ctxControl bool)
	SetLogFunction(logger func(ctx context.Context, logMessage LogMessage))
	SetSessionValidatorForUserTokens(validator SessionValidator) error
}
