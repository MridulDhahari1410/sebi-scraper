package gatekeeper

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Level int8

const (
	// InfoLevel defines info log level.
	InfoLevel Level = iota
	// WarnLevel defines warn log level.
	WarnLevel
)

type LogMessage struct {
	LogLevel Level
	Message  string
	err      error
}

type Options struct {
	Provider int
	Params   map[string]interface{}
}

type JWTTokenClaims struct {
	UserData      TokenUserData `json:"userData,omitempty"`
	OmneManagerID int16         `json:"omnemanagerid,omitempty"`
	Token         string        `json:"token,omitempty"`
	SourceID      string        `json:"sourceid,omitempty"`
	CustomClaims
	jwt.RegisteredClaims
	jwt.Claims
}

type TokenClaims struct {
	CustomClaims
	jwt.RegisteredClaims
}

type TokenUserData struct {
	CountryCode string    `json:"country_code,omitempty"`
	MobileNo    string    `json:"mob_no,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Source      string    `json:"source,omitempty"`
	AppID       string    `json:"app_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	DataCenter  string    `json:"dataCenter,omitempty"`
}

type AngelOneClaims struct {
	UserType   string                            `json:"user_type,omitempty"` // accepted value application, client, admin
	MobileNo   string                            `json:"mobile_no,omitempty"`
	TokenType  string                            `json:"token_type,omitempty"` // accepted value trade_access_token, trade_refresh_token, non_trade_access_token, non_trade_refresh_token
	DataCenter string                            `json:"data_center,omitempty"`
	GMId       int16                             `json:"gm_id,omitempty"`
	Source     string                            `json:"source,omitempty"`
	DeviceId   string                            `json:"device_id,omitempty"`
	Issuer     string                            `json:"issuer"`
	Subject    string                            `json:"subject,omitempty"`
	Audience   jwt.ClaimStrings                  `json:"audience,omitempty"`
	Scope      jwt.ClaimStrings                  `json:"scope,omitempty"` // will be used for authorization roles
	KeyId      string                            `json:"key_id"`
	Products   map[string]map[string]interface{} `json:"products,omitempty"`
}

type CustomClaims struct {
	UserType      string                            `json:"user_type,omitempty"` // accepted value application, client, admin
	MobileNo      string                            `json:"mobile_no,omitempty"`
	TokenType     string                            `json:"token_type,omitempty"` // accepted value trade_access_token, trade_refresh_token, non_trade_access_token, non_trade_refresh_token
	Scope         jwt.ClaimStrings                  `json:"scope,omitempty"`      // will be used for authorization roles
	DataCenter    string                            `json:"data_center,omitempty"`
	GMId          int16                             `json:"gm_id"`
	Source        string                            `json:"source,omitempty"`
	DeviceId      string                            `json:"device_id,omitempty"`
	KeyId         string                            `json:"kid,omitempty"`
	OmneManagerID int16                             `json:"omnemanagerid"`
	Products      map[string]map[string]interface{} `json:"products,omitempty"`
}

type Actor struct {
	Subject string `json:"sub,omitempty"`
}

func (c JWTTokenClaims) Valid() error {
	return c.RegisteredClaims.Valid()
}

type failedResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
