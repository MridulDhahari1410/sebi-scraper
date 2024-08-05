package api

import (
	"sebi-scrapper/constants"
	"sebi-scrapper/utils/flags"

	goActuator "github.com/angel-one/go-actuator"
	"github.com/gin-gonic/gin"
)

var actuatorHandler = goActuator.GetActuatorHandler(&goActuator.Config{
	Env:  flags.Env(),
	Name: constants.ApplicationName,
	Port: flags.Port(),
})

// actuator is used to handle the actuator requests.
func actuator(ctx *gin.Context) {
	actuatorHandler(ctx.Writer, ctx.Request)
}
