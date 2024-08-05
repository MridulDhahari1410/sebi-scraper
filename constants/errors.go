package constants

import (
	"net/http"

	"github.com/sinhashubham95/go-utils/errors"
)

// error constants.
var (
	ErrNoRecords              = errors.New("no records found in database")
	ErrNoRowsAffected         = errors.New("no records updated")
	ErrInvalidS2SAuthService  = errors.New("invalid s2s auth service")
	ErrTokenNotFoundForSymbol = errors.New("token not found for symbol")
)

// api errors.
var (
	ErrForbidden                           = errors.Error{StatusCode: http.StatusForbidden, Code: "400-1", Message: "Forbidden"}
	ErrParseJSON                           = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-2", Message: "Unexpected error occurred. Please contact customer support."}
	ErrServiceNoContentFailureInformation  = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-3", Message: "no content failure from external service"}
	ErrServiceBadRequestFailureInformation = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-4", Message: "bad request failure from external service"}
	ErrInvalidStrategyID                   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-5", Message: "invalid strategy id"}
	ErrInvalidStrategyName                 = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-6", Message: "invalid strategy name"}
	ErrInvalidStrategyMode                 = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-7", Message: "invalid strategy mode"}
	ErrInvalidStrategyDescription          = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-8", Message: "invalid strategy description"}
	ErrInvalidTradeExchange                = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-9", Message: "invalid trade exchange"}
	ErrInvalidTradeToken                   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-10", Message: "invalid trade token"}
	ErrInvalidTradeTokenType               = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-11", Message: "invalid trade token type"}
	ErrInvalidTradeQuantity                = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-12", Message: "invalid trade quantity"}
	ErrInvalidTradeTransactionType         = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-13", Message: "invalid trade transaction type"}
	ErrInvalidConditionTimeframe           = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-15", Message: "invalid condition timeframe"}
	ErrInvalidConditionEntry               = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-16", Message: "invalid condition entry"}
	ErrInvalidConditionExit                = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-17", Message: "invalid condition exit"}
	ErrInvalidStrategyType                 = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-18", Message: "invalid strategy type"}
	ErrUnauthorized                        = errors.Error{StatusCode: http.StatusUnauthorized, Code: "400-19", Message: "Unauthorized"}
	ErrInvalidConditionalOrdersID          = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-20", Message: "invalid conditional orders id"}
	ErrInvalidConditionalOrdersTimestamp   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-21", Message: "invalid conditional orders timestamp"}
	ErrInvalidConditionalOrdersToken       = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-22", Message: "invalid conditional orders token"}
	ErrInvalidStrategyQuery                = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-23", Message: "invalid strategy query"}
	ErrInvalidStrategyTimeFrame            = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-24", Message: "invalid strategy time frame"}
	ErrStrategyNotFound                    = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-25", Message: "strategy not found"}
	ErrInvalidNews                         = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-26", Message: "invalid input news"}
	ErrInvalidParamSize                    = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-27", Message: "Invalid size param"}
	ErrInvalidParamOrder                   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-28", Message: "Invalid order param"}
	ErrInvalidParamQuery                   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-29", Message: "Invalid query param"}
	ErrInvalidParamTokens                  = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-30", Message: "Invalid tokens param"}
	ErrInvalidParamtradingSymbols          = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-31", Message: "Invalid trading-symbols param"}
	ErrMissingParamSize                    = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-27", Message: "missing size param"}
	ErrMissingParamOrder                   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-28", Message: "missing order param"}
	ErrMissingParamQuery                   = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-29", Message: "missing query param"}
	ErrMissingParamTokens                  = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-30", Message: "missing tokens param"}
	ErrMissingParamtradingSymbols          = errors.Error{StatusCode: http.StatusBadRequest, Code: "400-31", Message: "missing trading-symbols param"}

	ErrDatabase                     = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-1", Message: "Unexpected error occurred. Please contact customer support."}
	ErrDatabaseCommit               = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-2", Message: "Unexpected error occurred. Please contact customer support."}
	ErrServiceUnavailable           = errors.Error{StatusCode: http.StatusServiceUnavailable, Code: "500-3", Message: "request timed out"}
	ErrServiceUnexpectedResponse    = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-4", Message: "unexpected response"}
	ErrServiceUnexpectedRequestBody = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-4", Message: "unexpected request body"}
	ErrServiceUnexpectedError       = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-5", Message: "unexpected error"}
	ErrUnexpectedMarshalError       = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-6", Message: "unexpected error"}
	ErrUnexpectedUnmarshalError     = errors.Error{StatusCode: http.StatusInternalServerError, Code: "500-7", Message: "unexpected error"}
)
