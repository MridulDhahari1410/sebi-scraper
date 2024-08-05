package constants

// api names.
const (
	ScripAndHolidayMasterAccessTokenRequestName          = "scrip-holiday-master-access-token"
	ScripMasterSecInfoRequestName                        = "scrip-master-sec-info"
	ConditionalOrdersCreateScreenerRequestName           = "conditionalOrdersCreateScreener"
	ConditionalOrdersCreateBacktestRequestName           = "conditionalOrdersCreateBacktest"
	ConditionalOrdersGetBacktestSummaryRequestName       = "conditionalOrdersGetBacktestSummary"
	NSESecurityArchivesForDeliveriesAndTradesRequestName = "nseSecurityArchivesForDeliveriesAndTrades"
	Nifty50Top10HoldingsRequestName                      = "nifty50Top10Holdings"
	NSEHistoricData                                      = "nseHistoricData"
)

// api statuses.
const (
	SuccessAPIStatus = "success"
	ErrorAPIStatus   = "error"
)

const (
	DepartmentQuery = "department"
	OrderQuery      = "order"
	AllDepartment   = "All Departments"
	Uncategorised   = "uncategorised"
)

// IDParamKey for id usage.
const IDParamKey = "id"

// header constants.
const (
	AccessTokenHeaderKey                = "AccessToken"
	AccessControlAllowOriginHeader      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentialsHeader = "Access-Control-Allow-Credentials"
	AccessControlAllowHeadersHeader     = "Access-Control-Allow-Headers"
	AccessControlAllowMethodsHeader     = "Access-Control-Allow-Methods"
	CacheControlHeader                  = "Cache-Control"
	CacheControlHeaderValue             = "must-revalidate, no-cache, no-store"
	ExpiresHeader                       = "Expires"
	ExpiresHeaderValue                  = "0"
	PragmaHeader                        = "Pragma"
	PragmaHeaderValue                   = "no-cache"
	XFrameOptionsHeader                 = "X-Frame-Options"
	XFrameOptionsHeaderValue            = "SAMEORIGIN"
	XXSSProtectionHeader                = "X-XSS-Protection"
	XXSSProtectionHeaderValue           = "1; mode=block"
	XContentTypeHeader                  = "X-Content-Type-Options"
	XContentTypeHeaderValue             = "nosniff"
	ContentSecurityPolicyHeader         = "Content-Security-Policy"
	ContentSecurityPolicyHeaderValue    = "default-src 'self' style-src 'self' 'unsafe-inline';"
	StrictTransportSecurityHeader       = "Strict-Transport-Security"
	StrictTransportSecurityHeaderValue  = "max-age=31536000; includeSubDomains; preload"
	RequestIDHeader                     = "X-requestId"
)

// query constants.
