package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DBConnectionString           = "host=%s port=%d dbname=%s sslmode=%s user=%s password=%s"
	DBConnectionStringWithSchema = "host=%s port=%d dbname=%s sslmode=%s user=%s password=%s search_path=%s"
)

const (
	DialectPostgres string = "postgres"
)

type contextKey int

const (
	// used to set the db instance in context in case of transactions
	ContextKeyDatabase contextKey = iota
)

var ErrorUndefinedDialect = errors.New("dialect for the db is not defined")

// IdbConfig interface has methods to read various DB configurations.
type IDbConnectionConfig interface {
	GetDialect() string
	GetDatabaseName() string
	GetConnectionPath() string
	GetMaxIdleConnections() int
	GetMaxOpenConnections() int
	GetConnMaxLifetime() time.Duration
	GetConnMaxIdleTime() time.Duration
	IsDebugMode() bool
}

// ConnectionConfig implements ConnectionReader.
type DbConnectionConfig struct {
	Dialect               string
	Protocol              string
	URL                   string
	Port                  int
	Username              string
	Password              string
	SslMode               string
	Name                  string
	Schema                string //can be used in Postgres (optional)
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionLifetime    time.Duration
	ConnectionMaxIdleTime time.Duration
	Debug                 bool
}

// GetDialect returns a dialect identifier
func (c DbConnectionConfig) GetDialect() string {
	return c.Dialect
}

// GetDatabaseName returns a database name
func (c DbConnectionConfig) GetDatabaseName() string {
	return c.Name
}

// GetConnectionPath returns connection string to be used by gorm basis dialect.
func (c DbConnectionConfig) GetConnectionPath() string {
	switch c.Dialect {
	case DialectPostgres:
		if c.Schema == "" {
			return fmt.Sprintf(DBConnectionString, c.URL, c.Port, c.Name, c.SslMode, c.Username, c.Password)
		}
		return fmt.Sprintf(DBConnectionStringWithSchema, c.URL, c.Port, c.Name, c.SslMode, c.Username, c.Password, c.Schema)
	default:
		return ""
	}
}

// GetMaxOpenConnections returns max open connections for the db.
func (c DbConnectionConfig) GetMaxOpenConnections() int {
	return c.MaxOpenConnections
}

// GetMaxIdleConnections returns max idle connections for the db.
func (c DbConnectionConfig) GetMaxIdleConnections() int {
	return c.MaxIdleConnections
}

// GetConnMaxLifetime returns configurable max lifetime of any connection of db.
func (c DbConnectionConfig) GetConnMaxLifetime() time.Duration {
	return c.ConnectionLifetime
}

// GetConnMaxIdleTime returns configurable max idle time of any connection of db.
func (c DbConnectionConfig) GetConnMaxIdleTime() time.Duration {
	return c.ConnectionMaxIdleTime
}

// IsDebugMode returns true if the debug logs for the DB are to be enabled
func (c DbConnectionConfig) IsDebugMode() bool {
	return c.Debug
}

// DB is the specific wrapper holding gorm db instance.
type DB struct {
	dbConfig   IDbConnectionConfig
	dialector  gorm.Dialector
	gormConfig *gorm.Config
	instance   *gorm.DB
}

func GormConfig(c *gorm.Config) func(*DB) error {
	return func(db *DB) error {
		db.gormConfig = c
		return nil
	}
}

func Dialector(gd gorm.Dialector) func(*DB) error {
	return func(db *DB) error {
		db.dialector = gd
		return nil
	}
}

func NewDb(dbConfig IDbConnectionConfig) (*DB, error) {
	if dbConfig == nil {
		dbConfig = &DbConnectionConfig{}
	}

	db := &DB{dbConfig: dbConfig}

	if db.dialector == nil {
		if err := db.initDialector(); err != nil {
			return nil, err
		}
	}

	if db.gormConfig == nil {
		db.gormConfig = &gorm.Config{
			AllowGlobalUpdate:      false,
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			// Set log level based on debug mode
			Logger: logger.Default.LogMode(getLogLevelByDebugMode(dbConfig.IsDebugMode())),
		}
	}

	if err := db.connect(); err != nil {
		return nil, err
	}

	return db, nil
}

// Instance returns underlying instance of gorm db.
// If the transaction/session in progress then it'll return
// the *gorm.DB from the context.
func (db *DB) Instance(ctx context.Context) *gorm.DB {
	if instance, ok := ctx.Value(ContextKeyDatabase).(*gorm.DB); ok {
		return instance
	}
	return db.instance
}

// GetInstance returns the instance from the gorm db
// This function skips the check in the context
func (db *DB) GetInstance(ctx context.Context) *gorm.DB {
	return db.instance
}

// Alive executes a select query and checks if connection exists and is alive.
func (db *DB) Alive() error {
	if dbi, err := db.instance.DB(); err != nil {
		return err
	} else {
		return dbi.Ping()
	}
}

func (db *DB) Dialector(ctx context.Context) gorm.Dialector {
	return db.dialector
}

// initDialector initializes a new dialector for the DB using the connReader
func (db *DB) initDialector() (err error) {
	var d gorm.Dialector
	if d, err = getDialector(db.dbConfig); err == nil {
		db.dialector = d
	}
	return
}

// connect opens a gorm connection and configures other connection details.
func (db *DB) connect() error {
	var err error

	if db.instance, err = gorm.Open(db.dialector, db.gormConfig); err != nil {
		return err
	}

	var dbConn *sql.DB
	if dbConn, err = db.instance.DB(); err != nil {
		return err
	}
	dbConn.SetMaxIdleConns(db.dbConfig.GetMaxIdleConnections())
	dbConn.SetMaxOpenConns(db.dbConfig.GetMaxOpenConnections())
	dbConn.SetConnMaxLifetime(db.dbConfig.GetConnMaxLifetime() * time.Second)
	dbConn.SetConnMaxIdleTime(db.dbConfig.GetConnMaxIdleTime() * time.Second)

	return nil
}

func getDialector(connReader IDbConnectionConfig) (gorm.Dialector, error) {
	switch connReader.GetDialect() {
	case DialectPostgres:
		return postgres.Open(connReader.GetConnectionPath()), nil
	default:
		return nil, errors.New("dialect for the db is not defined")
	}
}

func (db *DB) GetDatabaseName() string {
	return db.dbConfig.GetDatabaseName()
}

// getLogLevelByDebugMode return logger log level based on debug mode.
// If app db is in debug mode, make log level as info
// Default log level for gorm db is warning, overriding that by this method.
func getLogLevelByDebugMode(debug bool) logger.LogLevel {
	if !debug {
		return logger.Silent
	} else {
		return logger.Info
	}
}
