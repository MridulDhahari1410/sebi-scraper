package entities

import (
	"context"
	"fmt"
	"sebi-scrapper/constants"
	"time"
)

// CreateStrategy is used to create a new strategy.
const (
	SaveReport = iota
	CheckReportExists
	GetPublicReportsWithDepartmentCount
	GetAllPublicReportsCount
	GetPublicReportsWithDepartmentAsc
	GetAllPublicReportsAsc
	GetPublicReportsWithDepartmentDesc
	GetAllPublicReportsDesc
)

// Strategy is used for the operations on the strategy table.
type Reports struct {
	ID         int64     `json:"id"`
	Title      string    `json:"title"`
	Department string    `json:"department"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
	Content    []byte    `json:"content"`
	Count      int       `json:"count"`
	Exists     bool      `json:"exists"`
	Offset     int       `json:"offset"`
	Size       int       `json:"size"`
}

// GetIDQuery is used as the id query.
func (r *Reports) GetIDQuery() string {
	return constants.Empty
}

// IsIDQueryPermitted is used to check if id query is permitted.
func (r *Reports) IsIDQueryPermitted(_ context.Context) bool {
	return true
}

// GetIDValues is used to check for id values.
func (r *Reports) GetIDValues() []any {
	return []any{r.ID}
}

// GetAllQuery is used as the all query.
func (r *Reports) GetAllQuery(_ context.Context) (string, bool) {
	return constants.Empty, false
}

// GetNext is used to create a new instance.
func (r *Reports) GetNext() Entity {
	return new(Reports)
}

// BindRow is used to bind row to the entity.
func (r *Reports) BindRow(row Scanner) error {
	return row.Scan(&r.ID)
}

// GetFreshSaveQuery is used to have the fresh save query.
func (r *Reports) GetFreshSaveQuery() string {
	return constants.Empty
}

// IsFreshSavePermitted is used to check if fresh save is permitted.
func (r *Reports) IsFreshSavePermitted(_ context.Context) bool {
	return false
}

// GetFreshFieldValues is used as the fresh field values.
func (r *Reports) GetFreshFieldValues(_ string) []any {
	return nil
}

// GetSaveQuery is used as the save query.
func (r *Reports) GetSaveQuery() string {
	return constants.Empty
}

// IsSavePermitted is used to check if save query is permitted.
func (r *Reports) IsSavePermitted(_ context.Context) bool {
	return false
}

// GetFieldValues is used to get the field values for save query.
func (r *Reports) GetFieldValues(_ string) []any {
	return nil
}

// GetDeleteQuery is used to get the delete query.
func (r *Reports) GetDeleteQuery() string {
	return constants.Empty
}

// IsDeletePermitted is used to check if the delete query is permitted.
func (r *Reports) IsDeletePermitted(_ context.Context) bool {
	return false
}

// GetDeleteValues is used to get the delete values.
func (r *Reports) GetDeleteValues() []any {
	return nil
}

// GetQuery is used to query a single row.
func (r *Reports) GetQuery(ctx context.Context, code int) (string, bool) {
	switch code {
	case SaveReport:
		return `INSERT INTO sebi_reports (date, title, department, content, status) VALUES ($1, $2, $3, $4, 'pending') RETURNING id;`, true
	case CheckReportExists:
		return `SELECT EXISTS (
			SELECT 1
			FROM sebi_reports
			WHERE date = $1 AND title = $2
		);`, true
	case GetPublicReportsWithDepartmentCount:
		return `SELECT count(id) AS count FROM sebi_reports WHERE department = $1`, true
	case GetAllPublicReportsCount:
		return `SELECT count(id) AS count FROM sebi_reports`, true
	default:
		return constants.Empty, false
	}
}

// GetQueryValues is used to get the set of query values for strategy.
func (r *Reports) GetQueryValues(code int) []any {
	switch code {
	case SaveReport:
		fmt.Println([]interface{}{&r.Date, &r.Title, &r.Department})
		return []interface{}{&r.Date, &r.Title, &r.Department, &r.Content}
	case CheckReportExists:
		return []interface{}{&r.Date, &r.Title}
	case GetPublicReportsWithDepartmentCount:
		return []interface{}{&r.Department}
	default:
		return nil
	}
}

// GetMultiQuery is used to get multiple rows.
func (r *Reports) GetMultiQuery(_ context.Context, code int) (string, bool) {
	switch code {
	case GetPublicReportsWithDepartmentDesc:
		return `SELECT id, department, title, status, date FROM sebi_reports WHERE department = $1 ORDER BY date DESC LIMIT $2 OFFSET $3`, true
	case GetAllPublicReportsDesc:
		return `SELECT id, department, title, status, date FROM sebi_reports ORDER BY date DESC LIMIT $1 OFFSET $2`, true
	case GetPublicReportsWithDepartmentAsc:
		return `SELECT id, department, title, status, date FROM sebi_reports WHERE department = $1 ORDER BY date ASC LIMIT $2 OFFSET $3`, true
	case GetAllPublicReportsAsc:
		return `SELECT id, department, title, status, date FROM sebi_reports ORDER BY date ASC LIMIT $1 OFFSET $2`, true
	default:
		return constants.Empty, false
	}
}

// GetMultiQueryValues is used to get the set of multi query values.
func (r *Reports) GetMultiQueryValues(code int) []any {
	switch code {
	case GetPublicReportsWithDepartmentDesc:
		return []interface{}{&r.Department, &r.Size, &r.Offset}
	case GetAllPublicReportsDesc:
		return []interface{}{&r.Size, &r.Offset}
	case GetPublicReportsWithDepartmentAsc:
		return []interface{}{&r.Department, &r.Size, &r.Offset}
	case GetAllPublicReportsAsc:
		return []interface{}{&r.Size, &r.Offset}
	default:
		return nil
	}

}

// GetNextRaw is used to get the next entity instance.
func (r *Reports) GetNextRaw() RawEntity {
	return new(Reports)
}

// BindRawRow is used to bind raw row.
func (r *Reports) BindRawRow(code int, row Scanner) error {
	switch code {
	case SaveReport:
		return row.Scan(&r.ID)
	case CheckReportExists:
		return row.Scan(&r.Exists)
	case GetPublicReportsWithDepartmentDesc:
		return row.Scan(&r.ID, &r.Department, &r.Title, &r.Status, &r.Date)
	case GetAllPublicReportsDesc:
		return row.Scan(&r.ID, &r.Department, &r.Title, &r.Status, &r.Date)
	case GetPublicReportsWithDepartmentAsc:
		return row.Scan(&r.ID, &r.Department, &r.Title, &r.Status, &r.Date)
	case GetAllPublicReportsAsc:
		return row.Scan(&r.ID, &r.Department, &r.Title, &r.Status, &r.Date)
	case GetAllPublicReportsCount:
		return row.Scan(&r.Count)
	case GetPublicReportsWithDepartmentCount:
		return row.Scan(&r.Count)
	default:
		return nil
	}
}

// GetExec is used to get the exec query.
func (r *Reports) GetExec(_ int) string {
	return constants.Empty
}

// IsExecPermitted is used to check if execution is permitted.
func (r *Reports) IsExecPermitted(_ context.Context, _ int) bool {
	return false
}

// GetExecValues is used to get the execution values.
func (r *Reports) GetExecValues(_ int, _ string) []any {
	return nil
}
