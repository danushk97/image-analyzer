package storage

import (
	"context"
	"fmt"

	sql "github.com/danushk97/image-analyzer/pkg/storage/sql"
)

// Choice represents the type of the storage to use
type Choice string

// Store is the interface supporting all storage operations
type Store struct {
	db *sql.DB
}

var (
	// SQLChoice is the dialect for all sql storages
	SQLChoice Choice = "sql"
)

// Config defines the database config
type Config struct {
	// Choice defines the storage choice: sql, dynamodb
	Choice Choice
	// SQL specifies the configuration
	// if database choice is mysql
	SQL sql.DbConnectionConfig
}

// New gives back a storage instance
func New(ctx context.Context, config Config) (*Store, error) {
	switch config.Choice {
	case SQLChoice:
		db, err := sql.NewDb(&config.SQL)
		if err != nil {
			return nil, err
		}
		return &Store{db: db}, nil
	}

	return nil, fmt.Errorf("unknown database choice: %v", config.Choice)
}
