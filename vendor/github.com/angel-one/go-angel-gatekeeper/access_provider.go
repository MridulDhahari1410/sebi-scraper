package gatekeeper

import (
	"github.com/gin-gonic/gin"
)


type AccessProvider interface {
	ValidateAccess(ctx *gin.Context) (bool, error)
}

