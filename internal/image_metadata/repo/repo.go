package repo

import (
	"context"

	"github.com/danushk97/image-analyzer/internal/image_metadata/model/v1"

	"github.com/danushk97/image-analyzer/pkg/errors"
)

// Transactional is the interface for all task related to executing tasks
// within transactional block
type Transactional interface {
	// Transaction is used to execute given function inside transaction block
	Transaction(ctx context.Context, fn func(ctx context.Context) errors.
		IError) errors.IError
	// IsActive checks if transaction is active or not
	IsActive(ctx context.Context) bool
}

// Repo is the interface that is used to talk to storage layer
// This was added as an interface so that SQL and Dynamodb both
// can use implement this interface to write to their respective storage
// Those implementations are store under
// offer/repo/sql or
type Repo interface {
	Transactional

	CreateImageMetadata(context.Context, *model.ImageMetadata) errors.IError
}
