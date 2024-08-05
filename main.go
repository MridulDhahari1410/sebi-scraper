package main

import (
	"context"
	"fmt"
	"sebi-scrapper/constants"

	"sebi-scrapper/utils/flags"

	"github.com/sinhashubham95/go-utils/log"

	"sebi-scrapper/jobs"

	"sebi-scrapper/utils/auth"

	"sebi-scrapper/api"

	"sebi-scrapper/utils/configs"
	"sebi-scrapper/utils/database"
	"sebi-scrapper/utils/http"
	"sebi-scrapper/utils/metrics"

	timeUtils "sebi-scrapper/utils/time"
	"time"

	_ "github.com/lib/pq"
)

// @title			Trading Signals
// @version		1.0
// @description	Trading Signals
// @termsOfService	https://swagger.io/terms/
// @BasePath
func main() {
	ctx := context.Background()
	initConfigs(ctx)
	initLogger()
	initMetrics()
	initTime(ctx)
	initHTTP(ctx)
	initAuth(ctx)
	initJobs(ctx)
	initDatabase(ctx)
	defer closeDatabase(ctx)
	startRouter(ctx)
}

func initConfigs(ctx context.Context) {
	var err error
	if flags.Mode() == constants.TestMode {
		err = configs.InitTestModeConfigs(flags.BaseConfigPath(), constants.LoggerConfig, constants.JWTConfig,
			constants.ApplicationConfig, constants.APIConfig, constants.DatabaseConfig, constants.S2SAuthConfig)
	} else if flags.Mode() == constants.ReleaseMode {
		err = configs.InitReleaseModeConfigs(constants.LoggerConfig, constants.JWTConfig,
			constants.ApplicationConfig, constants.APIConfig, constants.DatabaseConfig, constants.S2SAuthConfig)
	}
	if err != nil {
		log.Fatal(ctx).Err(err).Msg("error loading configs")
	}
}

func initLogger() {
	log.InitLogger(log.Level(configs.Get().GetStringD(constants.LoggerConfig, constants.LogLevelConfigKey,
		constants.DefaultLogLevel)), configs.Get().GetStringSliceD(constants.LoggerConfig, constants.LogParamsConfigKey, nil))
}

func initMetrics() {
	metrics.Init(metrics.Config{Bucket: metrics.BucketConfig{
		Start: configs.Get().GetFloatD(constants.ApplicationConfig, constants.ApplicationMetricsLinearBucketStartConfigKey, 0),
		Width: configs.Get().GetFloatD(constants.ApplicationConfig, constants.ApplicationMetricsLinearBucketWidthConfigKey, 0),
		Count: int(configs.Get().GetIntD(constants.ApplicationConfig, constants.ApplicationMetricsLinearBucketCountConfigKey, 0)),
	}})
}

func initTime(ctx context.Context) {
	err := timeUtils.Init()
	if err != nil {
		log.Fatal(ctx).Err(err).Msg("error initialising time utils")
		return
	}
	log.Info(ctx).Msg("successfully initialised time")
}

func initHTTP(ctx context.Context) {
	http.InitHTTPClient(
		http.NewRequestConfig(constants.ConditionalOrdersCreateScreenerRequestName, configs.Get().
			GetMapD(constants.APIConfig, constants.APIConditionalOrdersCreateScreenerConfigKey, nil)).
			SetURL(configs.Get().GetStringWithEnvD(constants.APIConfig,
				constants.APIConditionalOrdersCreateScreenerURLConfigKey, constants.Empty)).
			SetHeaderParams(configs.Get().GetMapD(constants.APIConfig,
				constants.APIConditionalOrdersCreateScreenerHeadersConfigKey, nil)),
		http.NewRequestConfig(constants.ConditionalOrdersCreateBacktestRequestName, configs.Get().
			GetMapD(constants.APIConfig, constants.APIConditionalOrdersCreateBacktestConfigKey, nil)).
			SetURL(configs.Get().GetStringWithEnvD(constants.APIConfig,
				constants.APIConditionalOrdersCreateBacktestURLConfigKey, constants.Empty)).
			SetHeaderParams(configs.Get().GetMapD(constants.APIConfig,
				constants.APIConditionalOrdersCreateBacktestHeadersConfigKey, nil)),
		http.NewRequestConfig(constants.ConditionalOrdersGetBacktestSummaryRequestName, configs.Get().
			GetMapD(constants.APIConfig, constants.APIConditionalOrderGetBacktestSummaryConfigKey, nil)).
			SetHeaderParams(configs.Get().GetMapD(constants.APIConfig,
				constants.APIConditionalOrderGetBacktestSummaryHeadersLConfigKey, nil)),
		http.NewRequestConfig(constants.NSESecurityArchivesForDeliveriesAndTradesRequestName, configs.Get().
			GetMapD(constants.APIConfig, constants.APINSESecurityArchivesForDeliveriesAndTradesConfigKey, nil)).
			SetURL(configs.Get().GetStringWithEnvD(constants.APIConfig,
				constants.APINSESecurityArchivesForDeliveriesAndTradesURLConfigKey, constants.Empty)).
			SetHeaderParams(configs.Get().GetMapD(constants.APIConfig,
				constants.APINSESecurityArchivesForDeliveriesAndTradesHeadersConfigKey, nil)),
		http.NewRequestConfig(constants.Nifty50Top10HoldingsRequestName, configs.Get().
			GetMapD(constants.APIConfig, constants.APINifty50Top10HoldingsConfigKey, nil)).
			SetURL(configs.Get().GetStringWithEnvD(constants.APIConfig,
				constants.APINifty50Top10HoldingsURLConfigKey, constants.Empty)).
			SetHeaderParams(configs.Get().GetMapD(constants.APIConfig,
				constants.APINifty50Top10HoldingsHeadersConfigKey, nil)),
		http.NewRequestConfig(constants.NSEHistoricData, configs.Get().
			GetMapD(constants.APIConfig, constants.APINSEHistoricDataConfigKey, nil)).
			SetURL(configs.Get().GetStringWithEnvD(constants.APIConfig,
				constants.APINSEHistoricDataURLConfigKey, constants.Empty)).
			SetHeaderParams(configs.Get().GetMapD(constants.APIConfig,
				constants.APINSEHistoricDataHeadersConfigKey, nil)),
	)
	log.Info(ctx).Msg("successfully initialised http")
}

func initAuth(ctx context.Context) {
	err := auth.Init()
	if err != nil {
		log.Fatal(ctx).Err(err).Msg("error initialising auth")
	}
}

func initDatabase(ctx context.Context) {
	driverName, err := configs.Get().GetString(constants.DatabaseConfig, constants.DatabaseDriverNameConfigKey)
	if err != nil {
		log.Fatal(ctx).Stack().Err(err).Msg("error getting database driver name")
		return
	}
	url, err := configs.Get().GetStringWithEnv(constants.DatabaseConfig, constants.DatabaseURLConfigKey)
	if err != nil {
		log.Fatal(ctx).Stack().Err(err).Msg("error getting database url")
		return
	}
	err = database.InitDatabase(database.Config{
		DriverName: driverName,
		URL:        url,
		MaxOpenConnections: int(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseMaxOpenConnectionsConfigKey, constants.DatabaseDefaultMaxOpenConnections)),
		MaxIdleConnections: int(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseMaxIdleConnectionsConfigKey, constants.DatabaseDefaultMaxIdleConnections)),
		ConnectionMaxLifetime: time.Second * time.Duration(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseConnectionMaxLifetimeInSecondsConfigKey,
			constants.DatabaseDefaultConnectionMaxLifetimeInSeconds)),
		ConnectionMaxIdleTime: time.Second * time.Duration(configs.Get().GetIntD(constants.DatabaseConfig,
			constants.DatabaseConnectionMaxIdleTimeInSecondsConfigKey,
			constants.DatabaseDefaultConnectionMaxIdleTimeInSeconds)),
	})
	if err != nil {
		log.Fatal(ctx).Stack().Err(err).Msg("error initialising database")
	}
}

func closeDatabase(ctx context.Context) {
	err := database.Close()
	if err != nil {
		log.Fatal(ctx).Stack().Err(err).Msg("error closing database")
	}
}

func initJobs(ctx context.Context) {
	jobs.Init(ctx)
}

func startRouter(ctx context.Context) {
	// get the router
	router, err := api.GetRouter(flags.Mode())
	if err != nil {
		log.Fatal(ctx).Err(err).Msg("error preparing router")
		return
	}
	// now start router
	err = router.Run(fmt.Sprintf(":%d", flags.Port()))
	if err != nil {
		log.Fatal(ctx).Stack().Err(err).Msg("error starting router")
	}
}
