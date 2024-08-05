package apiv1

import (
	"strconv"

	"sebi-scrapper/constants"

	"github.com/gin-gonic/gin"
)

func getIDFromContext(ctx *gin.Context, key string) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param(key), 10, 64)
	if err != nil {
		return 0, constants.ErrInvalidStrategyID.WithDetails(err.Error())
	}
	if id < 0 {
		return 0, constants.ErrInvalidStrategyID.WithDetails("id cannot be less than 0")
	}
	return id, nil
}
