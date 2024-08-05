package constants

import (
	"time"
)

// ApplicationName is the identifier for a lot of external resources for this application.
const ApplicationName = "sebi-scrapper"

// run modes of the application.
const (
	TestMode    = "test"
	ReleaseMode = "release"
)

const Page = "page"
const DefaultPage = "1"

// constants.
const (
	Empty  = ""
	Star   = "*"
	True   = "true"
	Active = "active"
)

// DefaultLogLevel is used as a fallback to initialise logger.
const DefaultLogLevel = "debug"

// database init constants.
const (
	DatabaseDefaultMaxOpenConnections             = 20
	DatabaseDefaultMaxIdleConnections             = 10
	DatabaseDefaultConnectionMaxLifetimeInSeconds = 90
	DatabaseDefaultConnectionMaxIdleTimeInSeconds = 30
)

// ClientCodeKey is used for the client code in context.
const ClientCodeKey = "user_id"

// IndianTimeLocation for time utils.
const IndianTimeLocation = "Asia/Kolkata"

// ConditionalOrdersTimestampFormat is the timestamp format for conditional orders.
const ConditionalOrdersTimestampFormat = time.RFC3339

// context keys.
const (
	AuthenticatedKey = "authenticated"
	SubjectKey       = "subject"
)

// NSETimestampFormat is used for timestamp as date from NSE.
const NSETimestampFormat = "2006-01-02"

// NSERequestTimestampFormat is used for timestamp as date from NSE.
const NSERequestTimestampFormat = "02-01-2006"

const (
	Stock       = "stock"
	Ascending   = "asc"
	Descending  = "desc"
	DefaultSize = "5"
	Terms       = "terms"
)

var DepartmentToValue = map[string]string{
	"Alternative Investment Fund and Foreign Portfolio Investors Department": "75",
	"Corporation Finance Department":                                         "1",
	"Department Economic and Policy Analysis":                                "3",
	"Department of Debt and Hybrid Securities":                               "64",
	"Enforcement Department - 1":                                             "6",
	"Information Technology Department":                                      "35",
	"Integrated Surveillance Department":                                     "10",
	"Investment Management Department":                                       "9",
	"Legal Affairs Department 1":                                             "12",
	"Market Intermediaries Regulation and Supervision Department":            "14",
	"Market Regulation Department":                                           "15",
	"Office of Investor Assistance and Education":                            "19",
}

var Departments = []string{
	"Alternative Investment Fund and Foreign Portfolio Investors Department",
	"Corporation Finance Department",
	"Department Economic and Policy Analysis",
	"Department of Debt and Hybrid Securities",
	"Enforcement Department - 1",
	"Information Technology Department",
	"Integrated Surveillance Department",
	"Investment Management Department",
	"Legal Affairs Department 1",
	"Market Intermediaries Regulation and Supervision Department",
	"Market Regulation Department",
	"Office of Investor Assistance and Education",
}

const AllReports = "AllReports"

const AllReprtsValue = "-1"
