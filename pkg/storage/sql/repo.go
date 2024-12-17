package sql

import (
	"context"

	"gorm.io/gorm"

	"github.com/danushk97/image-analyzer/pkg/errors"
)

const updatedAtField = "updated_at"

type Repo struct {
	Db *DB
}

// FindByID fetches the record which matches the ID provided from the entity defined by receiver
// and the result will be loaded into receiver
func (repo Repo) FindByID(ctx context.Context, receiver IModel, id string) errors.IError {
	q := repo.DBInstance(ctx).Where("id = ?", id).First(receiver)

	return GetDBError(q)
}

// Create inserts a new record in the entity defined by the receiver
// all data filled in the receiver will inserted
// Also, modifies current instance db's context to ctx using WithContext.
func (repo Repo) Create(ctx context.Context, receiver IModel) errors.IError {
	if err := receiver.SetDefaults(); err != nil {
		return err
	}

	if err := receiver.Validate(); err != nil {
		return err
	}

	q := repo.DBInstance(ctx).WithContext(ctx).Create(receiver)

	return GetDBError(q)
}

// Delete deletes the given model
// Soft or hard delete of model depends on the models implementation
// if the model composites SoftDeletableModel then it'll be soft deleted
func (repo Repo) Delete(ctx context.Context, receiver IModel) errors.IError {
	q := repo.DBInstance(ctx).Delete(receiver)

	return GetDBError(q)
}

// Transaction will manage the execution inside a transactions
// adds the txn db in the context for downstream use case
func (repo Repo) Transaction(ctx context.Context, fc func(ctx context.Context) errors.IError) errors.IError {
	var err = repo.DBInstance(ctx).Transaction(func(tx *gorm.DB) error {

		// This will ensure that when db.Instance(context) we return the txn on the context
		// & all repo queries are done on this txn. Refer usage in test.
		if err := fc(context.WithValue(ctx, ContextKeyDatabase, tx)); err != nil {
			return err
		}

		return GetDBError(tx)
	})

	if err == nil {
		return nil
	}

	// tx.Commit can throw an error which will not be an IError
	if iErr, ok := err.(errors.IError); ok {
		return iErr
	}

	// use the default code and wrap err in internal
	return errors.NewServerError(errDBError).Wrap(err)
}

// IsTransactionActive returns true if a transaction is active
func (repo Repo) IsTransactionActive(ctx context.Context) bool {
	_, ok := ctx.Value(ContextKeyDatabase).(*gorm.DB)
	return ok
}

// DBInstance returns gorm instance.
// If replicas are specified, for Query, Row callback, will use replicas, unless Write mode specified.
// For Raw callback, statements are considered read-only and will use replicas if the SQL starts with SELECT.
func (repo Repo) DBInstance(ctx context.Context) *gorm.DB {
	return repo.Db.Instance(ctx)
}
