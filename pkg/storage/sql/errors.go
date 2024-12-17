package sql

import (
	goerr "errors"

	"github.com/danushk97/image-analyzer/pkg/errors"
	"gorm.io/gorm"
)

const (
	errDBError           = "db_error"
	errNoRowAffected     = "no_row_affected"
	errRecordNotFound    = "record_not_found"
	errValidationFailure = "validation_failure"
)

// GetDBError accepts db instance and the details
// creates appropriate error based on the type of query result
// if there is no error then returns nil
func GetDBError(db *gorm.DB) errors.IError {
	if db.Error == nil {
		return nil
	}

	// Construct error based on type of db operation
	err := func() errors.IError {
		switch true {
		case goerr.Is(db.Error, gorm.ErrRecordNotFound):
			return errors.NewBadRequestError(errRecordNotFound)

		default:
			return errors.NewServerError(errDBError)
		}
	}()

	// add specific details of error
	return err.Wrap(db.Error)
}

// GetValidationError wraps the error and returns instance of ValidationError
// if the provided error is nil then it just returns nil
func GetValidationError(err error) errors.IError {
	if err != nil {
		return errors.NewBadRequestError(errValidationFailure).
			Wrap(err)
	}

	return nil
}
