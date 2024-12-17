package dtos

import (
	"github.com/danushk97/image-analyzer/pkg/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CreateImageMetadata defines the structure of the create request for image metadata
type CreateImageMetadataRequest struct {
	FileName string `json:"file_name"` // Name of the image
}

func (c *CreateImageMetadataRequest) Validate() errors.IError {

	err := validation.ValidateStruct(
		c,
		validation.Field(
			&c.FileName,
			validation.Required,
			validation.Length(1, 255),
		),
	)

	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}

	return nil
}
