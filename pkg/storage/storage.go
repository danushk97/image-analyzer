package storage

import (
	"context"
	"fmt"

	"github.com/danushk97/image-analyzer/pkg/errors"
	sql "github.com/danushk97/image-analyzer/pkg/storage/sql"
)

// Choice represents the type of the storage to use
type Choice string

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

// Store is the interface supporting all storage operations
type Store interface {
	Create(
		ctx context.Context, receiver sql.IModel) errors.IError
	FindByID(
		ctx context.Context, receiver sql.IModel, id string) errors.IError
	Delete(
		ctx context.Context, receiver sql.IModel) errors.IError
}

// New gives back a storage instance
func New(ctx context.Context, config Config) (Store, error) {
	switch config.Choice {
	case SQLChoice:
		db, err := sql.NewDb(&config.SQL)
		if err != nil {
			return nil, err
		}
		return &sql.Repo{Db: db}, nil
	}

	return nil, fmt.Errorf("unknown database choice: %v", config.Choice)
}
