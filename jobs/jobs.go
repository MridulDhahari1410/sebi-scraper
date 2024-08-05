package jobs

import (
	"context"
)

// Init is used to initialise the jobs.
func Init(ctx context.Context) {
	// crons.New().At(configs.Get().GetStringD(constants.ApplicationConfig, constants.ApplicationJobsSebiPublicReportsAtConfigKey,
	// 	"10:00")).Start(ctx, SebiPublicReports)
}
