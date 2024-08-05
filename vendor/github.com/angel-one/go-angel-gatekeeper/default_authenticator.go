package gatekeeper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	configs "github.com/angel-one/go-config-client"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type defaultAuthenticator struct {
	keyProvider           KeyProvider
	commonKeyProvider     KeyProvider
	ctxControl            bool
	applicationName       string
	logger                func(ctx context.Context, logMessage LogMessage)
	tokenStorageValidator tokenStorageValidator
	commonConfigProvider  configProvider
}

func NewDefaultAuthenticator(configClient configs.Client, options configs.Options, applicationName string) (Authenticator, error) {

	commonConfigProvider, err := getCommonConfigProvider(configClient, options.Provider, options.Params)
	if err != nil {
		return nil, err
	}
	authenticationConfigName := getAuthenticationConfigName(options)
	applicationConfigProvider, err := newConfigClientBasedConfigProvider(configClient, authenticationConfigName, empty)
	if err != nil {
		return nil, err
	}

	return &defaultAuthenticator{
		keyProvider:           newConfigAsymmetricKeyProvider(applicationConfigProvider),
		commonKeyProvider:     newCommonKeyProvider(commonConfigProvider),
		applicationName:       applicationName,
		ctxControl:            false,
		commonConfigProvider:  commonConfigProvider,
		tokenStorageValidator: nil,
		logger:                nil,
	}, nil
}

// NewDefaultAuthenticatorWithKeyProvider TODO : this function doesn't support reading the common configuration from option we can add as needed.
func NewDefaultAuthenticatorWithKeyProvider(keyProvider KeyProvider, options Options, applicationName string) (Authenticator, error) {

	commonConfigProvider, err := getCommonConfigProvider(nil, options.Provider, options.Params)
	if err != nil {
		return nil, err
	}
	return &defaultAuthenticator{
		keyProvider:           keyProvider,
		applicationName:       applicationName,
		ctxControl:            false,
		commonKeyProvider:     newCommonKeyProvider(commonConfigProvider),
		commonConfigProvider:  commonConfigProvider,
		tokenStorageValidator: nil,
		logger:                nil,
	}, nil
}

func getAuthenticationConfigName(options configs.Options) string {
	val, ok := options.Params[authenticationConfigName]
	if ok {
		if valueString, isOk := val.(string); isOk {
			return valueString
		}
	}
	return configTypeApplication
}

func (authenticator *defaultAuthenticator) SetCtxControl(ctxControl bool) {
	authenticator.ctxControl = ctxControl
}

func (authenticator *defaultAuthenticator) SetSessionValidatorForUserTokens(validator SessionValidator) error {
	if authenticator.commonConfigProvider == nil {
		return errors.New("common config provider is not initialized properly")
	}
	tokenStorageValidator, err := newTokenStorageValidator(authenticator.commonConfigProvider, authenticator.logger)
	if err != nil {
		return err
	}
	authenticator.tokenStorageValidator = tokenStorageValidator
	authenticator.tokenStorageValidator.setSessionValidator(validator)
	return nil
}

func (authenticator *defaultAuthenticator) SetLogFunction(logger func(ctx context.Context, logMessage LogMessage)) {
	authenticator.logger = logger
}

func (authenticator *defaultAuthenticator) ValidateToken(ctx *gin.Context, opts ...TypeValidator) (*TokenClaims, error) {
	token, readTokenErr := authenticator.getToken(ctx, opts)
	if readTokenErr != nil {
		return nil, readTokenErr
	}
	loginClaims, err := authenticator.verifyToken(token, opts)
	if loginClaims == nil {
		return nil, err
	}
	tokenClaims := &TokenClaims{
		CustomClaims:     loginClaims.CustomClaims,
		RegisteredClaims: loginClaims.RegisteredClaims,
	}
	return tokenClaims, err
}

func (authenticator *defaultAuthenticator) verifyToken(token string, opts []TypeValidator) (*JWTTokenClaims, error) {
	loginClaims, parseTokenErr := authenticator.parseTokenWithKey(token)
	if parseTokenErr != nil {
		return loginClaims, parseTokenErr
	}
	loginClaims.Subject = authenticator.getSubjectFromSubOrUserData(loginClaims)
	loginClaims.GMId = authenticator.getGMId(loginClaims)
	loginClaims.CustomClaims.OmneManagerID = loginClaims.GMId
	loginClaims.MobileNo = authenticator.getMobileNumber(loginClaims)
	loginClaims.Source = authenticator.getSource(loginClaims)
	loginClaims.DeviceId = authenticator.getDeviceId(loginClaims)
	validationErr := loginClaims.Valid()

	if validationErr != nil {
		return loginClaims, validationErr
	}
	loginClaims.Token = token
	if !authenticator.isCorrectTokenType(loginClaims, opts) {
		return loginClaims, ErrTokenTypeValidationFailed
	}

	if !authenticator.isCorrectAudience(loginClaims) {
		return loginClaims, ErrAudienceValidationFailed
	}

	sessionValidationErr := authenticator.isValidSession(loginClaims, opts)
	return loginClaims, sessionValidationErr
}

func (authenticator *defaultAuthenticator) isValidSession(loginClaims *JWTTokenClaims, opts []TypeValidator) error {
	if loginClaims.UserType == UserTypeClient && authenticator.storageValidatorExist(opts) {
		if authenticator.tokenStorageValidator == nil {
			return ErrSessionValidatorNotInitialized
		}
		valid, err := authenticator.tokenStorageValidator.validate(loginClaims)
		if err != nil {
			return err
		}
		if !valid {
			return ErrSessionValidationFailed
		}
	}
	return nil
}

func getCookiesHandler(c *gin.Context, opts []TypeValidator) string {
	cookies := c.Request.Cookies()

	if len(cookies) == 0 {
		return empty
	}

	cookiesMap := make(map[string]string)
	for _, cookie := range cookies {
		cookiesMap[cookie.Name] = cookie.Value
	}

	var tokenTypeValidator TokenTypeValidator
	for _, opt := range opts {
		tokenTypeValidator, _ = opt.(TokenTypeValidator)
	}

	if len(tokenTypeValidator.TokenTypes) == 0 {
		if c.GetHeader(tokenTypeHeader) != empty {
			return cookiesMap[c.GetHeader(tokenTypeHeader)]
		} else {
			if value, ok := cookiesMap[TradeAccessToken]; ok {
				return value
			} else if value, ok = cookiesMap[NonTradeAccessToken]; ok {
				return value
			} else {
				return empty
			}
		}
	} else {
		if len(tokenTypeValidator.TokenTypes) == 1 {
			tokenType := tokenTypeValidator.TokenTypes[0]
			return cookiesMap[tokenType]
		} else if c.GetHeader(tokenTypeHeader) != empty {
			for _, tokenType := range tokenTypeValidator.TokenTypes {
				if tokenType == c.GetHeader(tokenTypeHeader) {
					return cookiesMap[c.GetHeader(tokenTypeHeader)]
				}
			}
			return empty
		}
	}

	return empty
}

func (authenticator *defaultAuthenticator) getToken(ctx *gin.Context, opts []TypeValidator) (string, error) {
	authHeader := ctx.Request.Header.Get(authorization)
	tokenSlice := strings.Split(authHeader, tokenSeparator)
	accessTokenHeader := ctx.Request.Header.Get(Accesstoken)
	token := empty
	if len(tokenSlice) == 2 && (tokenSlice[0] == bearer || tokenSlice[0] == bearerCaps) {
		token = tokenSlice[1]
	} else if accessTokenHeader != empty {
		token = accessTokenHeader
	} else if token = getCookiesHandler(ctx, opts); token != empty {
		return token, nil
	} else {
		return empty, errors.New("incorrect Authorization or Accesstoken")
	}
	return token, nil
}

func (authenticator *defaultAuthenticator) Authenticate(opts ...TypeValidator) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		token, readTokenErr := authenticator.getToken(ctx, opts)
		if readTokenErr != nil {
			authenticator.handleUnauthorized(ctx, readTokenErr, "error incorrect Authorization or Accesstoken", nil)
			return
		}

		loginClaims, parseErr := authenticator.parseTokenWithKey(token)
		if parseErr != nil {
			authenticator.handleUnauthorized(ctx, parseErr, "Invalid signature", loginClaims)
			return
		}
		validationErr := loginClaims.Valid()
		if validationErr != nil {
			authenticator.handleUnauthorized(ctx, validationErr, fmt.Sprintf("Invalid token or expired for subject :%s", loginClaims.Subject), loginClaims)
			return
		}

		if !authenticator.isCorrectTokenType(loginClaims, opts) {
			authenticator.handleUnauthorized(
				ctx,
				ErrTokenTypeValidationFailed,
				fmt.Sprintf("Incorrect token for subject : %s, tokenType : %s, userType : %s", loginClaims.Subject, loginClaims.TokenType, loginClaims.UserType),
				loginClaims,
			)
			return
		}

		if !authenticator.isCorrectAudience(loginClaims) {
			authenticator.handleUnauthorized(
				ctx,
				ErrAudienceValidationFailed,
				fmt.Sprintf("Incorrect audience for subject : %s, audiences : %s", loginClaims.Subject, loginClaims.Audience),
				loginClaims,
			)
			return
		}
		loginClaims.Token = token
		sessionValidationErr := authenticator.isValidSession(loginClaims, opts)
		if sessionValidationErr != nil {
			if sessionValidationErr == ErrSessionValidationFailed {
				authenticator.handleUnauthorized(ctx, ErrSessionValidationFailed, fmt.Sprintf("session validation failed for subject : % s", loginClaims.Subject), loginClaims)
				return
			}
			authenticator.handleError(ctx, sessionValidationErr)
			return
		}

		ctx.Set(authenticated, true)
		ctx.Set(UserType, loginClaims.UserType)
		ctx.Set(TokenType, loginClaims.TokenType)
		ctx.Set(issuer, loginClaims.Issuer)
		ctx.Set(IssuedAt, loginClaims.IssuedAt)
		ctx.Set(ExpiresAt, loginClaims.ExpiresAt)
		ctx.Set(audience, loginClaims.Audience)
		ctx.Set(roles, loginClaims.Scope)
		ctx.Set(dataCenter, loginClaims.DataCenter)
		ctx.Set(subject, authenticator.getSubjectFromSubOrUserData(loginClaims))
		ctx.Set(gmId, authenticator.getGMId(loginClaims))
		ctx.Set(mobileNo, authenticator.getMobileNumber(loginClaims))
		ctx.Set(deviceId, authenticator.getDeviceId(loginClaims))
		ctx.Set(source, authenticator.getSource(loginClaims))
		ctx.Set(Accesstoken, token)
		ctx.Set(JTI, loginClaims.ID)
		switch loginClaims.UserType {
		case UserTypeApplication, UserTypeClient, UserTypeAdmin:
		default:
			authenticator.logMessage(ctx, WarnLevel, fmt.Sprintf("Token is valid but %s is not a valid usertype", loginClaims.UserType), fmt.Errorf("token is valid but %s is not a valid usertype", loginClaims.UserType))
		}
		authenticator.logMessage(ctx,
			InfoLevel,
			fmt.Sprintf(
				"usertype %s with subject %s authorized to access url path %s, issuer: %s, keyId : %s",
				loginClaims.UserType,
				loginClaims.Subject,
				ctx.FullPath(),
				loginClaims.Issuer,
				loginClaims.KeyId,
			), nil,
		)
		if !authenticator.ctxControl {
			ctx.Next()
		}
		return //nil
	}
}

func (authenticator *defaultAuthenticator) getSource(loginClaims *JWTTokenClaims) string {
	source := loginClaims.Source
	if source == empty {
		return loginClaims.SourceID
	}
	return source
}

func (authenticator *defaultAuthenticator) getDeviceId(loginClaims *JWTTokenClaims) string {
	deviceId := loginClaims.DeviceId
	if deviceId == empty {
		return loginClaims.UserData.AppID
	}
	return deviceId
}

func (authenticator *defaultAuthenticator) getGMId(loginClaims *JWTTokenClaims) int16 {
	gmId := loginClaims.GMId
	if gmId == 0 {
		return loginClaims.OmneManagerID
	}
	return gmId
}

func (authenticator *defaultAuthenticator) getMobileNumber(loginClaims *JWTTokenClaims) string {
	mobileNo := loginClaims.MobileNo
	if mobileNo == empty {
		return loginClaims.UserData.MobileNo
	}
	return mobileNo
}

func (authenticator *defaultAuthenticator) handleUnauthorized(ctx *gin.Context, err error, msg string, loginClaims *JWTTokenClaims) {
	authenticator.logMessage(ctx, WarnLevel, msg, err)
	if !authenticator.ctxControl {
		ctx.JSON(http.StatusUnauthorized, failedResponse{
			Status:  "error",
			Message: "Invalid or Expired token",
		})
		ctx.Abort()
	} else {
		ctx.Set(authenticated, false)
		ctx.Set(authenticationError, err)
		if loginClaims != nil {
			ctx.Set(subject, loginClaims.Subject)
		}
	}
}

func (authenticator *defaultAuthenticator) handleError(ctx *gin.Context, err error) {
	authenticator.logMessage(ctx, WarnLevel, err.Error(), err)
	if !authenticator.ctxControl {
		ctx.JSON(http.StatusInternalServerError, failedResponse{
			Status:  "error",
			Message: "Internal Server Error",
		})
		ctx.Abort()
	} else {
		ctx.Set(authenticated, false)
		ctx.Set(sessionValidationError, err)
	}

}

func (authenticator *defaultAuthenticator) isCorrectTokenType(claims *JWTTokenClaims, opts []TypeValidator) bool {
	if opts == nil || len(opts) == 0 {
		return true
	}
	for _, opt := range opts {
		tokenTypeValidator, ok := opt.(TokenTypeValidator)
		if ok {
			return tokenTypeValidator.validate(claims)
		}
	}
	return true
}

func (authenticator *defaultAuthenticator) isCorrectAudience(claims *JWTTokenClaims) bool {
	if claims.UserType == UserTypeApplication {
		// Split applicationName based on comma
		appNames := strings.Split(authenticator.applicationName, ",")
		appNameMap := make(map[string]struct{})
		for _, appName := range appNames {
			appNameMap[appName] = struct{}{}
		}
		for _, aud := range claims.Audience {
			if _, ok := appNameMap[aud]; ok {
				return true
			}
		}
		return false
	}
	return true
}

func (authenticator *defaultAuthenticator) getSubjectFromSubOrUserData(claims *JWTTokenClaims) string {
	userId := claims.Subject
	if userId == empty {
		userId = claims.UserData.UserID
	}
	return userId
}

// VerifyToken verifies the authenticity of the token using the provided public key and expected payload.
func (authenticator *defaultAuthenticator) parseTokenWithKey(token string) (*JWTTokenClaims, error) {

	claims := JWTTokenClaims{}

	verifiedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {

		issuer := claims.Issuer //unverified claims
		// old token support, this should be removed once everyone migrates to new library
		if issuer == angel {
			key, configErr := authenticator.commonKeyProvider.GetPrivateKey(nonTradeTokenJwtSigningKey)
			return []byte(key), configErr
		} else if issuer == empty {
			key, configErr := authenticator.commonKeyProvider.GetPrivateKey(tradeTokenJwtSigningKey)
			return []byte(key), configErr
		} else if issuer == angelBroking || issuer == angelOne {
			key, configErr := authenticator.commonKeyProvider.GetPrivateKey(s2sTokenJwtSigningKey)
			return []byte(key), configErr
		} else {
			// new token verification
			var key string
			var configErr error
			if issuer == tradeLoginService || issuer == nonTradeLoginService {
				key, configErr = authenticator.commonKeyProvider.GetPublicKey(issuer, claims.KeyId)
			} else {
				if authenticator.keyProvider == nil {
					return []byte(empty), errors.New("application key provider is nil")
				}
				if claims.UserType == UserTypeClient && issuer != refreshLoginService && issuer != loginService {
					return []byte(empty), errors.New("usertype client is not supported for s2s calls")
				}
				key, configErr = authenticator.keyProvider.GetPublicKey(issuer, claims.KeyId)
			}
			if configErr != nil {
				return []byte(key), configErr
			}
			return jwt.ParseRSAPublicKeyFromPEM([]byte(key))
		}
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}
	return verifiedToken.Claims.(*JWTTokenClaims), err
}

func (authenticator *defaultAuthenticator) logMessage(ctx context.Context, logLevel Level, logMessage string, err error) {
	if authenticator.logger != nil {
		authenticator.logger(ctx, LogMessage{
			LogLevel: logLevel,
			Message:  logMessage,
			err:      err,
		})
	}
}

func (authenticator *defaultAuthenticator) storageValidatorExist(opts []TypeValidator) bool {
	if opts == nil || len(opts) == 0 {
		return false
	}
	for _, opt := range opts {
		_, ok := opt.(StorageValidator)
		if ok {
			return true
		}
	}
	return false
}
