package service

import (
	imageSql "github.com/danushk97/image-analyzer/internal/image_metadata/repo/sql"
	"github.com/danushk97/image-analyzer/pkg/storage"
	sql "github.com/danushk97/image-analyzer/pkg/storage/sql"
)

// Option is an option to offer Service to set
// the dependencies and configurations
type Option func(*Service)

// WithStorage adds the storage being used for offer storage
func WithStorage(
	store storage.Store,
) Option {
	return func(opts *Service) {
		switch s := store.(type) {
		case *sql.Repo:
			opts.Repo = imageSql.NewRepo(s)
		}
	}
}

// NewOptions will create a new builder Service object and
// apply all the options to that object and returns pointer
// to the builder Service
func NewOptions(opts ...Option) *Service {
	s := &Service{}
	// Loop through each option
	for _, op := range opts {
		op(s)
	}

	return s
}
