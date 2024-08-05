package entities

import (
	"context"

	timeUtils "sebi-scrapper/utils/time"
)

// Scanner is used to scan the data to the respective types.
type Scanner interface {
	Scan(dest ...any) error
}

// Entity is the set of common methods in an entity.
type Entity interface {
	GetIDQuery() string
	IsIDQueryPermitted(ctx context.Context) bool
	GetIDValues() []any
	GetAllQuery(ctx context.Context) (string, bool)
	GetNext() Entity
	BindRow(row Scanner) error
	GetFreshSaveQuery() string
	IsFreshSavePermitted(ctx context.Context) bool
	GetFreshFieldValues(source string) []any
	GetSaveQuery() string
	IsSavePermitted(ctx context.Context) bool
	GetFieldValues(source string) []any
	GetDeleteQuery() string
	IsDeletePermitted(ctx context.Context) bool
	GetDeleteValues() []any
}

// RawEntity is the set of common methods for doing raw queries with an entity.
type RawEntity interface {
	GetQuery(ctx context.Context, code int) (string, bool)
	GetQueryValues(code int) []any
	GetMultiQuery(ctx context.Context, code int) (string, bool)
	GetMultiQueryValues(code int) []any
	GetNextRaw() RawEntity
	BindRawRow(code int, row Scanner) error
	GetExec(code int) string
	IsExecPermitted(ctx context.Context, code int) bool
	GetExecValues(code int, source string) []any
}

// RawExec is the structure for the entity and the code.
type RawExec struct {
	Entity RawEntity
	Code   int
}

func getAuditValues(source string) []any {
	return []any{source, timeUtils.GetCurrentTime(), source, timeUtils.GetCurrentTime()}
}
