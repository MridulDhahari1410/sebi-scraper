package middlewares

import (
	"sebi-scrapper/constants"

	"github.com/gin-gonic/gin"
)

// Security is used to set the various headers and response parameters for security vulnerabilities.
// https://beaglesecurity.com/blog/article/hardening-server-security-by-implementing-security-headers.html
// https://caniuse.com
func Security() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// browser cache and history
		ctx.Header(constants.CacheControlHeader, constants.CacheControlHeaderValue)
		ctx.Header(constants.ExpiresHeader, constants.ExpiresHeaderValue)
		ctx.Header(constants.PragmaHeader, constants.PragmaHeaderValue)

		// frame options
		ctx.Header(constants.XFrameOptionsHeader, constants.XFrameOptionsHeaderValue)

		// xss protection
		ctx.Header(constants.XXSSProtectionHeader, constants.XXSSProtectionHeaderValue)

		// content type
		ctx.Header(constants.XContentTypeHeader, constants.XContentTypeHeaderValue)

		// strict transport security
		ctx.Header(constants.StrictTransportSecurityHeader, constants.StrictTransportSecurityHeaderValue)

		// content security policy
		ctx.Header(constants.ContentSecurityPolicyHeader, constants.ContentSecurityPolicyHeaderValue)
	}
}
