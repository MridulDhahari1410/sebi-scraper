package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sinhashubham95/go-utils/errors"

	"time"

	"sebi-scrapper/constants"
	"sebi-scrapper/entities"

	"github.com/sinhashubham95/go-utils/log"
)

// Config is the set of configurable parameters for database.
type Config struct {
	DriverName            string        `json:"driverName"`
	URL                   string        `json:"-"`
	MaxOpenConnections    int           `json:"maxOpenConnections"`
	MaxIdleConnections    int           `json:"maxIdleConnections"`
	ConnectionMaxLifetime time.Duration `json:"connectionMaxLifetime"`
	ConnectionMaxIdleTime time.Duration `json:"connectionMaxIdleTime"`
}

// Database is the wrapper around the native database.
type Database interface {
	GetByID(ctx context.Context, entity entities.Entity) error
	GetAll(ctx context.Context, entity entities.Entity) ([]entities.Entity, error)
	FreshSave(ctx context.Context, source string, entities ...entities.Entity) error
	Save(ctx context.Context, source string, entities ...entities.Entity) error
	Delete(ctx context.Context, entities ...entities.Entity) error
	QueryRaw(ctx context.Context, entity entities.RawEntity, code int) error
	QueryMultiRaw(ctx context.Context, entity entities.RawEntity, code int) ([]entities.RawEntity, error)
	ExecRaws(ctx context.Context, source string, execs ...entities.RawExec) error
}

// dbClient is the default client.
type dbClient struct {
	*sql.DB
}

// dbTxClient is the transactional client.
type dbTxClient struct {
	*sql.Tx
}

var db *dbClient

// InitDatabase is used to initialise the database.
func InitDatabase(config Config) error {
	// open the database
	conn, err := sql.Open(
		config.DriverName,
		config.URL,
	)

	if err != nil {
		return err
	}

	// try to ping
	err = conn.Ping()
	if err != nil {
		return err
	}

	// now set the configurations
	conn.SetMaxOpenConns(config.MaxOpenConnections)
	conn.SetMaxIdleConns(config.MaxIdleConnections)
	conn.SetConnMaxIdleTime(config.ConnectionMaxIdleTime)
	conn.SetConnMaxLifetime(config.ConnectionMaxLifetime)

	// now create the database
	db = &dbClient{DB: conn}
	query := `CREATE TABLE sebi_reports (
		id SERIAL PRIMARY KEY,
		date DATE NOT NULL,
		title TEXT NOT NULL,
		content BYTEA NOT NULL,
		department TEXT NOT NULL,
		status VARCHAR(30) NOT NULL
	);`
	_, err = db.Exec(query)
	fmt.Println("Error DB : ", err)
	return nil
}

// Get is used to get the database instance client.
func Get() Database {
	return db
}

// GetTx is used to get the transactional instance client.
func GetTx(ctx context.Context, options *sql.TxOptions) Database {
	t, err := db.DB.BeginTx(ctx, options)
	if err != nil {
		log.Error(ctx).Err(err).Msg("Error while begin transaction")
	}
	return &dbTxClient{Tx: t}
}

// CommitTx is used to commit the transaction.
func CommitTx(db Database) error {
	if t, ok := db.(*dbTxClient); ok {
		err := t.Commit()
		if err != nil {
			return constants.ErrDatabaseCommit.WithDetails(err.Error())
		}
	}
	return nil
}

// RollbackTx is used to roll back the transaction.
func RollbackTx(_ context.Context, db Database) {
	if t, ok := db.(*dbTxClient); ok {
		_ = t.Rollback()
	}
}

// Close is used to close the database instance.
func Close() error {
	return db.Close()
}

func handleGetAllResponse(ctx context.Context, rows *sql.Rows, err error, entity entities.Entity) ([]entities.Entity, error) {
	if errors.Is(err, sql.ErrNoRows) {
		return nil, constants.ErrNoRecords
	}
	if err != nil {
		return nil, constants.ErrDatabase.WithDetails(err.Error())
	}
	defer closeRows(ctx, rows)
	result := make([]entities.Entity, 0)
	for rows.Next() {
		err = entity.BindRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
		entity = entity.GetNext()
	}
	if len(result) == 0 {
		return nil, constants.ErrNoRecords
	}
	return result, nil
}

func handleQueryMultiRawResponse(ctx context.Context, rows *sql.Rows, err error, entity entities.RawEntity, code int) ([]entities.RawEntity, error) {
	if errors.Is(err, sql.ErrNoRows) {
		return nil, constants.ErrNoRecords
	}
	if err != nil {
		return nil, constants.ErrDatabase.WithDetails(err.Error())
	}
	defer closeRows(ctx, rows)
	result := make([]entities.RawEntity, 0)
	for rows.Next() {
		err = entity.BindRawRow(code, rows)
		if err != nil {
			return nil, constants.ErrDatabase.WithDetails(err.Error())
		}
		result = append(result, entity)
		entity = entity.GetNextRaw()
	}
	if len(result) == 0 {
		return nil, constants.ErrNoRecords
	}
	return result, nil
}

func closeRows(ctx context.Context, rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Error(ctx).Err(err).Msg("error closing rows")
	}
}
