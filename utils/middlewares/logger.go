package middlewares

import (
	"time"

	"sebi-scrapper/constants"
	"sebi-scrapper/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sinhashubham95/go-utils/log"
)

// LoggerMiddlewareOptions is the set of configurable allowed for log.
type LoggerMiddlewareOptions struct {
	NotLogQueryParams  bool
	NotLogHeaderParams bool
	SkipHeaderParams   []string
	SkipQueryParams    []string
}

// Logger is the middleware to be used for logging the request and response information
// This should be the first middleware to be added, in case the recovery middleware is not being used.
// Otherwise, it should be the second one, just after the recovery middleware.
func Logger(options LoggerMiddlewareOptions) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetHeader(constants.RequestIDHeader)
		if id == "" {
			// get a unique id
			uid, err := uuid.NewUUID()
			if err == nil {
				id = uid.String()
			}
		}
		// apply the id in the context
		ctx.Set(constants.IDLogKey, id)

		path := utils.GetHandlerName(ctx)
		// apply the pah in the context
		ctx.Set(constants.PathLogKey, path)

		// Start timer
		start := time.Now()

		// log the initial request
		event := log.Info(ctx).
			Str(constants.StartTimeLogKey, start.Format(time.RFC3339)).
			Str(constants.MethodLogKey, ctx.Request.Method).
			Str(constants.URILogKey, ctx.Request.URL.Path)
		if !options.NotLogQueryParams {
			if len(options.SkipQueryParams) > 0 {
				queryParamMap := map[string][]string{}
				utils.Copy(queryParamMap, ctx.Request.URL.Query())
				for _, param := range options.SkipQueryParams {
					delete(queryParamMap, param)
				}
				event.Interface(constants.QueryLogKey, queryParamMap)
			} else {
				event.Str(constants.QueryLogKey, ctx.Request.URL.RawQuery)
			}
		}
		if !options.NotLogHeaderParams {
			if len(options.SkipHeaderParams) > 0 {
				headersMap := map[string][]string{}
				utils.Copy(headersMap, ctx.Request.Header)
				for _, header := range options.SkipHeaderParams {
					delete(headersMap, header)
				}
				event.Interface(constants.HeaderLogKeys, headersMap)
			} else {
				event.Interface(constants.HeaderLogKeys, ctx.Request.Header)
			}
		}
		event.Send()

		// Process request
		ctx.Next()

		// stop timer
		end := time.Now()
		latency := end.Sub(start)

		// log the final response details
		log.Info(ctx).
			Int(constants.StatusCodeLogKey, ctx.Writer.Status()).
			Str(constants.MethodLogKey, ctx.Request.Method).
			Str(constants.URILogKey, ctx.Request.URL.Path).
			Str(constants.ClientIPLogKey, ctx.ClientIP()).
			Str(constants.ErrorLogKey, ctx.Errors.ByType(gin.ErrorTypePrivate).String()).
			Str(constants.StartTimeLogKey, start.Format(time.RFC3339)).
			Str(constants.EndTimeLogKey, end.Format(time.RFC3339)).
			Int64(constants.LatencyInMillisLogKey, latency.Milliseconds()).
			Send()
	}
}
