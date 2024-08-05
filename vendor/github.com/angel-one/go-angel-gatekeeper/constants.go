package gatekeeper

import "errors"

const (
	authentication           = "authentication"
	publicKey                = "publicKey"
	defaultRefreshInterval   = 3600
	refreshIntervalInSec     = "refreshIntervalInSec"
	consulUrl                = "CONSUL_URL"
	consulConfiguration      = "consul-configuration"
	consulToken              = "CONSUL_TOKEN"
	consulTokenHeaderKey     = "X-Consul-Token"
	headers                  = "headers"
	authenticationConfigName = "authenticationConfigName"
	globalLoginConfigName    = "globalLoginConfigName"
	defaultGlobalConfigName  = "global-login-configuration"
	globalLoginSecretName    = "globalLoginSecretName"
	defaultGlobalSecretName  = "global-login-secret"

	defaultConnectionAge         = 90 * 1000
	defaultMinRetryBackOff       = 1
	defaultMaxRetryBackOff       = 10
	defaultIdleConnections       = 1
	defaultIdleConnectionTimeout = 90 * 1000
	defaultPoolTimeout           = 1 * 1000
	defaultContextTimeout        = 1 * 1000
	defaultReadTimeout           = 1 * 1000
	defaultDialTimeout           = 1 * 1000
	defaultRetries               = 3
	defaultDb                    = 0
	defaultCacheType             = "redis"

	devEnv = "dev"

	empty = ""

	authorization = "Authorization"
	Accesstoken   = "Accesstoken"

	bearer     = "Bearer"
	bearerCaps = "BEARER"

	tokenSeparator = " "

	subject    = "subject"
	gmId       = "gmId"
	issuer     = "issuer"
	roles      = "roles"
	dataCenter = "dataCenter"
	audience   = "Audience"
	JTI        = "jti"

	NonTradeAccessToken = "non_trade_access_token"
	TradeAccessToken    = "trade_access_token"
	S2SAccessToken      = "s2s_access_token"

	tradeRefreshToken    = "trade_refresh_token"
	nonTradeRefreshToken = "non_trade_refresh_token"
	accessTokenCookie    = "access_token"
	guest                = "guest"
	session              = "session"

	tokenTypeHeader = "X-tokenType"

	configTypeApplication = "application"

	UserType            = "userType"
	TokenType           = "tokenType"
	UserTypeApplication = "application"
	UserTypeClient      = "client"
	UserTypeAdmin       = "admin"
	authenticated       = "authenticated"
	authenticationError = "authenticationError"

	useConsul = "useConsul"

	applicationName          = "global-login-configuration"
	yaml                     = "yaml"
	configIDKey              = "id"
	configRegionKey          = "region"
	configAppKey             = "app"
	configEnvKey             = "env"
	accessKeyId              = "accessKeyId"
	secretKey                = "secretKey"
	configTypeKey            = "configType"
	secretTypeKey            = "secretType"
	configNamesKey           = "configNames"
	secretNamesKey           = "secretNames"
	secretsDirectoryKey      = "secretsDirectory"
	configsDirectoryKey      = "configsDirectory"
	configCredentialsModeKey = "credentialsMode"
	ExpiresAt                = "exp"
	IssuedAt                 = "iat"

	configurations = "configurations"
	secrets        = "secrets"

	angel                                      = "angel"
	enableOldJWTSupportForNonTradeLoginService = "enableOldJWTSupportForNonTradeLoginService"
	enableOldJWTSupportForTradeLoginService    = "enableOldJWTSupportForTradeLoginService"
	enableS2SSupport                           = "enableS2SSupport"
	angelBroking                               = "AngelBroking"
	angelOne                                   = "AngelOne"
	tradeLoginService                          = "trade_login_service"
	nonTradeLoginService                       = "non_trade_login_service"
	refreshLoginService                        = "refresh_login_service"
	loginService                               = "login_service"
	globalLoginCredentials                     = "global-login-credentials"
	nonTradeTokenJwtSigningKey                 = "non_trade_token_key"
	tradeTokenJwtSigningKey                    = "trade_token_key"
	s2sTokenJwtSigningKey                      = "s2s_token_key"
	mobileNo                                   = "mobileNo"
	source                                     = "source"
	deviceId                                   = "deviceId"

	persistence            = "persistence"
	readTimeout            = "readTimeout"
	dialTimeout            = "dialTimeout"
	poolTimeout            = "poolTimeout"
	cachePassword          = "password"
	address                = "address"
	cacheType              = "cacheType"
	addresses              = "addresses"
	db                     = "db"
	idleConnectionTimeout  = "idleConnectionTimeout"
	contextTimeout         = "contextTimeout"
	clusterModeEnabled     = "clusterModeEnabled"
	maxConnectionAge       = "maxConnectionAge"
	minIdleConnections     = "minIdleConnections"
	poolSize               = "poolSize"
	cacheUsername          = "username"
	maxRetries             = "maxRetries"
	minRetryBackoff        = "minRetryBackoff"
	maxRetryBackoff        = "maxRetryBackoff"
	inmemory               = "inmemory"
	redis                  = "redis"
	sessionValidationError = "sessionValidationError"

	issuedTimeTolerance              = 60
	defaultBufferCacheSeconds        = 2
	s2sApplicationTokenJwtSigningKey = "s2s_token_key_application"
)

var ErrTokenTypeValidationFailed = errors.New("incorrect Token type")

var ErrAudienceValidationFailed = errors.New("incorrect Audience type")

var ErrSessionValidatorNotInitialized = errors.New("token storage validator is not initialised")
var ErrSessionValidationFailed = errors.New("session Validation Failed")

var ConsulConfiguration = map[string]interface{}{
	"method": "GET",
	"headers": map[string]interface{}{
		"content-type": "application/json",
		"accept":       "*/*",
	},
	"timeoutinmillis":             5000,
	"connecttimeoutinmillis":      2000,
	"tlshandshaketimeoutinmillis": 1000,
	"retrycount":                  3,
	"backoffpolicy": map[string]interface{}{
		"constantbackoff": map[string]interface{}{
			"intervalinmillis":          5,
			"maxjitterintervalinmillis": 2,
		},
	},
	"hystrixconfig": map[string]interface{}{
		"hystrixtimeoutinmillis": 5000,
		"maxconcurrentrequests":  100,
		"errorpercentthresold":   50,
		"sleepwindowinmillis":    5000,
		"requestvolumethreshold": 20,
	},
}

const (
	FileBased = iota
	AWSAppConfig
)

var invalidIfPresent = func(val interface{}, claims JWTTokenClaims) bool {
	if val != nil {
		return false
	}
	return true
}

var validIfPresent = func(val interface{}, claims *JWTTokenClaims) bool {
	if val != nil {
		return true
	}
	return false
}
