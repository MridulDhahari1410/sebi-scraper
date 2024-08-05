package middlewares

import (
	"net/http"

	"sebi-scrapper/constants"

	"github.com/gin-gonic/gin"
)

// CORS is used as a cors middleware.
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header(constants.AccessControlAllowOriginHeader, constants.Star)
		ctx.Header(constants.AccessControlAllowCredentialsHeader, constants.True)
		ctx.Header(constants.AccessControlAllowHeadersHeader, constants.Star)
		ctx.Header(constants.AccessControlAllowMethodsHeader, "POST,OPTIONS,GET,PUT")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
