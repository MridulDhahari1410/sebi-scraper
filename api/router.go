package api

import (
	"strings"

	apiv1 "sebi-scrapper/api/v1"
	"sebi-scrapper/utils/flags"
	"sebi-scrapper/utils/metrics"

	"sebi-scrapper/utils/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"sebi-scrapper/constants"
)

// GetRouter is used to get the instance of the engine powering the various apis.
func GetRouter(mode string) (*gin.Engine, error) {
	if mode == constants.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery(), middlewares.CORS(), middlewares.Security())

	router.GET(constants.ActuatorRoute, actuator)
	router.GET(constants.MetricsRoute, metrics.HTTPMetrics())

	if strings.TrimSpace(strings.ToLower(flags.Env())) != constants.EnvProd {
		router.GET(constants.SwaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1NonTradeRouter := router.Group(constants.V1Route, metrics.GetMetricsMiddleware(),
		middlewares.Logger(middlewares.LoggerMiddlewareOptions{}))
	{
		addSebiNonTradeRoutes(v1NonTradeRouter)
	}

	return router, nil
}

func addSebiNonTradeRoutes(group *gin.RouterGroup) {
	sebi := group.Group(constants.SebiRoute)
	{

		sebi.GET(constants.StartSebiCrawlerRoute, apiv1.HandleCrawlSebiReports)
		sebi.GET(constants.GetSebiPublicReports, apiv1.HandleGetPublicReports)
		sebi.GET(constants.GetSebiReportsDepartmentList, apiv1.HandleGetReportsDepartmentsList)
		sebi.GET(constants.GetSebiPublicReportsCount, apiv1.HandleGetPublicReportsCount)
	}
}
