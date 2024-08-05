package constants

const (
	ActuatorRoute = "/actuator/*any"
	MetricsRoute  = "/metrics"
	SwaggerRoute  = "/swagger/*any"

	V1Route                      = "v1"
	V1InternalRoute              = "v1/internal"
	SebiRoute                    = "/sebi"
	StartSebiCrawlerRoute        = "/crawlPublicReports"
	GetStrategyDetailsRoute      = "/:id/details"
	GetScripsForStrategyRoute    = "/:id/scrips"
	TradeInStrategyRoute         = "/:id/trade"
	StrategyScreenerCallback     = "/screener/callback"
	NewsRoute                    = "/news"
	GetSebiPublicReports         = "/publicReports"
	GetSebiReportsDepartmentList = "/departments"
	GetSebiPublicReportsCount    = "/publicReportsCount"
)
