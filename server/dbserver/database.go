package dbserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nordluma/go-bookstore/config"
	"github.com/nordluma/go-bookstore/values"
)

var (
	// Initialize database
	InitializeDb = initializeDb

	// Create db runner fore master database and put it into context
	PrepareDbRunner = prepareDbRunner

	dbHandler *sql.DB
)

func initializeDb() (err error) {
	connctionString := config.GetDatabaseConnectionString()
	maxIdleConnections := config.GetDatabaseMaxIdleConnections()
	maxOpenConnections := config.GetDatabaseMaxOpenConnections()
	connectionLifetime := config.GetDatabaseConnectonLifetime()

	dbHandler, err = initDbHandle(
		"master",
		"postgres",
		connctionString,
		maxIdleConnections,
		maxOpenConnections,
		connectionLifetime,
	)
	if err != nil {
		return
	}

	return
}

func createRunner(db *sql.DB) Runner {
	run := new(dbRunner)
	run.db = db

	return run
}

func prepareDbRunner(ctx context.Context) context.Context {
	return context.WithValue(
		ctx,
		values.ContextKeyDbRunner,
		createRunner(dbHandler),
	)
}

func initDbHandle(
	name, dbtype, connectionString string,
	maxIdleConnections, maxOpenConnections int,
	connectionLifetime time.Duration,
) (*sql.DB, error) {
	if dbtype == "" {
		return nil, errors.New("Database type is empty")
	}

	if connectionString == "" {
		return nil, errors.New("Connection string is empty")
	}

	dbHandler, err := sql.Open(dbtype, connectionString)
	if err != nil {
		return nil, err
	}

	dbHandler.SetMaxIdleConns(maxIdleConnections)
	dbHandler.SetMaxOpenConns(maxOpenConnections)
	dbHandler.SetConnMaxLifetime(connectionLifetime)

	err = validateDb(dbHandler)
	if err != nil {
		dbHandler.Close()
	}

	return dbHandler, nil
}

func validateDb(dbHandler *sql.DB) error {
	err := dbHandler.Ping()
	if err != nil {
		return err
	}

	timeZone, err := readDatabaseTimeZone(context.Background(), dbHandler)
	if err != nil {
		return err
	}

	if timeZone != "UTC" {
		err = fmt.Errorf(
			"Database 'timezone' must be set to 'UTC'. Currently it's '%v'",
			timeZone,
		)

		return err
	}

	return nil
}

func readDatabaseTimeZone(
	ctx context.Context,
	dbHandler *sql.DB,
) (timeZone string, err error) {
	rowsTimeZone, err := dbHandler.QueryContext(ctx, "show timezone")
	if err != nil {
		return
	}

	defer rowsTimeZone.Close()
	if !rowsTimeZone.Next() {
		fmt.Errorf("No time zone")
		return
	}

	err = rowsTimeZone.Scan(&timeZone)

	return
}
