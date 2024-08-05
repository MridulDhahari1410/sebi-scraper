package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/sinhashubham95/go-utils/errors"

	"sebi-scrapper/constants"
	"sebi-scrapper/entities"
	"sebi-scrapper/utils/metrics"

	"github.com/sinhashubham95/go-utils/log"
)

// GetByID is used to get the entity by its primary key.
func (d *dbClient) GetByID(ctx context.Context, entity entities.Entity) error {
	if entity.IsIDQueryPermitted(ctx) {
		q := entity.GetIDQuery()
		timer := metrics.GetDBQueryTimer("DAO_GET_BY_ID", reflect.TypeOf(entity).String(), q)
		row := d.QueryRowContext(ctx, q, entity.GetIDValues()...)
		timer.ObserveDuration()
		err := entity.BindRow(row)
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrNoRecords
		}
		if err != nil {
			return constants.ErrDatabase.WithDetails(err.Error())
		}
		return nil
	}
	return constants.ErrForbidden.Value()
}

// GetAll is used to get all the entries for a particular entity.
func (d *dbClient) GetAll(ctx context.Context, entity entities.Entity) ([]entities.Entity, error) {
	q, ok := entity.GetAllQuery(ctx)
	if ok {
		timer := metrics.GetDBQueryTimer("DAO_GET_ALL", reflect.TypeOf(entity).String(), q)
		rows, err := d.QueryContext(ctx, q)
		timer.ObserveDuration()
		return handleGetAllResponse(ctx, rows, err, entity)
	}
	return nil, constants.ErrForbidden.Value()
}

// FreshSave is used to freshly save the information of the provided entities.
func (d *dbClient) FreshSave(ctx context.Context, source string, entities ...entities.Entity) error {
	// begin the transaction for saving information for these entities
	tx, err := d.Begin()
	if err != nil {
		return constants.ErrDatabase.WithDetails(err.Error())
	}

	// now try to perform the queries sequentially
	for _, entity := range entities {
		if entity.IsFreshSavePermitted(ctx) {
			q := entity.GetFreshSaveQuery()
			timer := metrics.GetDBQueryTimer("DAO_FRESH_SAVE", reflect.TypeOf(entity).String(), q)
			_, err = tx.ExecContext(ctx, q, entity.GetFreshFieldValues(source)...)
			timer.ObserveDuration()
			if err != nil {
				return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
			}
		} else {
			return constants.ErrForbidden.Value()
		}
	}

	// if all good, commit
	err = tx.Commit()
	if err != nil {
		return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
	}

	// success
	return nil
}

// Save is used to save the information of the provided entities.
func (d *dbClient) Save(ctx context.Context, source string, entities ...entities.Entity) error {
	// begin the transaction for saving information for these entities
	tx, err := d.Begin()
	if err != nil {
		return constants.ErrDatabase.WithDetails(err.Error())
	}

	// now try to perform the queries sequentially
	for _, entity := range entities {
		if entity.IsSavePermitted(ctx) {
			q := entity.GetSaveQuery()
			timer := metrics.GetDBQueryTimer("DAO_SAVE", reflect.TypeOf(entity).String(), q)
			result, err := tx.ExecContext(ctx, q, entity.GetFieldValues(source)...)
			timer.ObserveDuration()
			if err != nil {
				return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
			}
			if rows == 0 {
				err = d.rollbackBecauseOfError(tx, constants.ErrNoRowsAffected)
				if err != nil {
					log.Warn(ctx).Err(err).Msg("error rolling back for no rows affected")
				}
				return constants.ErrNoRowsAffected
			}
		} else {
			return constants.ErrForbidden.Value()
		}
	}

	// if all good, commit
	err = tx.Commit()
	if err != nil {
		return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
	}

	// success
	return nil
}

// Delete is used to delete the information of the provided entities.
func (d *dbClient) Delete(ctx context.Context, entities ...entities.Entity) error {
	// begin the transaction for deleting information for these entities
	tx, err := d.Begin()
	if err != nil {
		return constants.ErrDatabase.WithDetails(err.Error())
	}

	// now try to perform the queries sequentially
	for _, entity := range entities {
		if entity.IsDeletePermitted(ctx) {
			q := entity.GetDeleteQuery()
			timer := metrics.GetDBQueryTimer("DAO_DELETE", reflect.TypeOf(entity).String(), q)
			_, err = tx.ExecContext(ctx, q, entity.GetDeleteValues()...)
			timer.ObserveDuration()
			if err != nil {
				return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
			}
		} else {
			return constants.ErrForbidden.Value()
		}
	}

	// if all good, commit
	err = tx.Commit()
	if err != nil {
		return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
	}

	// success
	return nil
}

// QueryRaw is used to execute raw query.
func (d *dbClient) QueryRaw(ctx context.Context, entity entities.RawEntity, code int) error {

	q, ok := entity.GetQuery(ctx, code)
	fmt.Println("q : ", q)
	if ok {
		timer := metrics.GetDBQueryTimer("DAO_QUERY_RAW", reflect.TypeOf(entity).String(), q)
		// fmt.Println("Values : ", entity.GetQueryValues(code))
		row := d.QueryRowContext(ctx, q, entity.GetQueryValues(code)...)
		fmt.Println("Row :", row)
		timer.ObserveDuration()
		err := entity.BindRawRow(code, row)
		fmt.Println("Err 1 : ", err)
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrNoRecords
		}
		if err != nil {
			return constants.ErrDatabase.WithDetails(err.Error())
		}
		fmt.Println("Success")
		return nil
	}
	return constants.ErrForbidden.Value()
}

// QueryMultiRaw is used to query multiple rows according to the given query.
func (d *dbClient) QueryMultiRaw(ctx context.Context, entity entities.RawEntity, code int) ([]entities.RawEntity, error) {
	q, ok := entity.GetMultiQuery(ctx, code)
	if ok {
		timer := metrics.GetDBQueryTimer("DAO_QUERY_MULTI_RAW", reflect.TypeOf(entity).String(), q)
		rows, err := d.QueryContext(ctx, q, entity.GetMultiQueryValues(code)...)
		timer.ObserveDuration()
		return handleQueryMultiRawResponse(ctx, rows, err, entity, code)
	}
	return nil, constants.ErrForbidden.Value()
}

// ExecRaws is used to execute raw queries.
func (d *dbClient) ExecRaws(ctx context.Context, source string, execs ...entities.RawExec) error {
	// begin the transaction for execution to happen
	tx, err := d.Begin()
	if err != nil {
		return constants.ErrDatabase.WithDetails(err.Error())
	}

	// now try to perform the executions sequentially
	for _, exec := range execs {
		if exec.Entity.IsExecPermitted(ctx, exec.Code) {
			q := exec.Entity.GetExec(exec.Code)
			timer := metrics.GetDBQueryTimer("DAO_EXEC_RAWS", reflect.TypeOf(execs).String(), q)
			result, err := tx.ExecContext(ctx, q, exec.Entity.GetExecValues(exec.Code, source)...)
			timer.ObserveDuration()
			if err != nil {
				return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
			}
			if rows == 0 {
				err = d.rollbackBecauseOfError(tx, constants.ErrNoRowsAffected)
				if err != nil {
					log.Warn(ctx).Err(err).Msg("error rolling back for no rows affected")
				}
				return constants.ErrNoRowsAffected
			}
		} else {
			return constants.ErrForbidden.Value()
		}
	}

	// if all good, commit
	err = tx.Commit()
	if err != nil {
		return constants.ErrDatabase.WithDetails(d.rollbackBecauseOfError(tx, err).Error())
	}

	// success
	return nil
}

func (d *dbClient) rollbackBecauseOfError(tx *sql.Tx, err error) error {
	return fmt.Errorf("rolling back (%w) because of error %v", tx.Rollback(), err)
}
