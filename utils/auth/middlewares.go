package auth

import (
	"sebi-scrapper/constants"
	"sebi-scrapper/models"
	"sebi-scrapper/utils/flags"

	gatekeeper "github.com/angel-one/go-angel-gatekeeper"
	"github.com/gin-gonic/gin"
)

// InternalAuthenticator is used to authenticate internal requests.
func InternalAuthenticator() gin.HandlerFunc {
	internalAccessToken := flags.InternalAccessToken()
	return func(ctx *gin.Context) {
		if ctx.GetHeader(constants.AccessTokenHeaderKey) != internalAccessToken {
			ctx.AbortWithStatusJSON(models.GetErrorResponse(constants.ErrUnauthorized.Value()))
			return
		}
		ctx.Next()
	}
}

// TTAuthenticator is used as the authenticator for trade token.
func TTAuthenticator() gin.HandlerFunc {
	return authenticator.Authenticate(gatekeeper.TokenTypeValidator{
		UserTypes:  []string{gatekeeper.UserTypeClient},
		TokenTypes: []string{gatekeeper.TradeAccessToken},
	})
}

// NTTAuthenticator is used as the authenticator for non-trade token.
func NTTAuthenticator() gin.HandlerFunc {
	return authenticator.Authenticate(gatekeeper.TokenTypeValidator{
		UserTypes:  []string{gatekeeper.UserTypeClient},
		TokenTypes: []string{gatekeeper.NonTradeAccessToken},
	})
}

// S2SAuthenticator is used as the authenticator.
func S2SAuthenticator() gin.HandlerFunc {
	return authenticator.Authenticate(gatekeeper.TokenTypeValidator{
		UserTypes: []string{gatekeeper.UserTypeApplication},
	})
}

// TTAuthenticationIdentity is used as the identity for tt authentication.
func TTAuthenticationIdentity() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ctx.GetBool(constants.AuthenticatedKey) {
			ctx.AbortWithStatusJSON(models.GetErrorResponse(constants.ErrUnauthorized.Value()))
			return
		}
		ctx.Set(constants.ClientCodeKey, ctx.GetString(constants.SubjectKey))
		ctx.Next()
	}
}

// NTTAuthenticationIdentity is used as the identity for ntt authentication.
func NTTAuthenticationIdentity() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ctx.GetBool(constants.AuthenticatedKey) {
			ctx.AbortWithStatusJSON(models.GetErrorResponse(constants.ErrUnauthorized.Value()))
			return
		}
		ctx.Set(constants.ClientCodeKey, ctx.GetString(constants.SubjectKey))
		ctx.Next()
	}
}

// S2SAuthenticationIdentity is used as the identity for s2s authentication.
func S2SAuthenticationIdentity() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ctx.GetBool(constants.AuthenticatedKey) {
			ctx.AbortWithStatusJSON(models.GetErrorResponse(constants.ErrUnauthorized.Value()))
			return
		}
		ctx.Next()
	}
}
