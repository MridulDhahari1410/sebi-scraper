package gatekeeper

import "github.com/gin-gonic/gin"

type Authorizer interface {
	Authorize(opt ...interface{}) gin.HandlerFunc
	VerifyAuthorization(ctx *gin.Context, opt ...AccessValidator) (bool, error)
}
