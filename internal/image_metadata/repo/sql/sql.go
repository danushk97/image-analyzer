package sql

import (
	"context"

	"gorm.io/gorm"

	internalErr "github.com/danushk97/image-analyzer/internal/errors"
	"github.com/danushk97/image-analyzer/internal/image_metadata/model/v1"
	"github.com/danushk97/image-analyzer/pkg/errors"
	pkgLogger "github.com/danushk97/image-analyzer/pkg/logger"
	"github.com/danushk97/image-analyzer/pkg/storage/sql"
)

// Repo is used to interact offers with the storage
// It can keep offer domain parts and use storage pkg as a layer to write
type Repo struct {
	dataStore *sql.Repo
}

// NewRepo creates a new repo for interacting with storage
func NewRepo(db *sql.Repo) *Repo {
	return &Repo{
		dataStore: db,
	}
}

// InstanceWithContext returns underlying instance of gorm db
// with the context attached to it
func (r Repo) InstanceWithContext(ctx context.Context) *gorm.DB {
	return r.dataStore.DBInstance(ctx).
		WithContext(context.WithoutCancel(ctx))
}

// Transaction performs given function inside the transaction block using
// spine transaction method
func (r Repo) Transaction(ctx context.Context,
	fn func(ctx context.Context) errors.IError) errors.IError {
	return r.dataStore.Transaction(ctx, func(ctx context.Context) errors.
		IError {
		return fn(ctx)
	})
}

// IsActive checks if transaction is active or not
func (r Repo) IsActive(ctx context.Context) bool {
	return r.dataStore.IsTransactionActive(ctx)
}

// CreateOffer creates a new offer in the Aurora DB. It creates different
// entities associated with offer in transactional manner.
func (r Repo) CreateImageMetadata(
	ctx context.Context,
	image *model.ImageMetadata,
) errors.IError {
	logger := pkgLogger.Ctx(ctx)
	err := r.dataStore.Create(ctx, image)

	if err != nil {
		logger.WithError(err).Error(
			"IMAGE_METADATA_CREATE_ERROR",
		)
		return errors.NewServerError(
			internalErr.ServerErrorDBCreateError).
			Wrap(err)
	}

	return nil
}
