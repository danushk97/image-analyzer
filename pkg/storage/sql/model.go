package sql

import (
	"github.com/danushk97/image-analyzer/pkg/datatype"
	"github.com/danushk97/image-analyzer/pkg/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	AttributeID        = "id"
	AttributeCreatedAt = "created_at"
	AttributeUpdatedAt = "updated_at"
	AttributeDeletedAt = "deleted_at"
)

type Model struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type IModel interface {
	TableName() string
	EntityName() string
	GetID() string
	GetPrimaryKey() string
	Validate() errors.IError
	SetDefaults() errors.IError
}

// Validate validates base Model.
func (m *Model) Validate() errors.IError {
	return GetValidationError(
		validation.ValidateStruct(
			m,
			validation.Field(&m.CreatedAt, validation.By(datatype.IsTimestamp)),
			validation.Field(&m.UpdatedAt, validation.By(datatype.IsTimestamp)),
		),
	)
}

// GetID gets identifier of entity.
func (m *Model) GetID() string {
	return m.ID
}

// GetCreatedAt gets created time of entity.
func (m *Model) GetCreatedAt() int64 {
	return m.CreatedAt
}

// GetUpdatedAt gets last updated time of entity.
func (m *Model) GetUpdatedAt() int64 {
	return m.UpdatedAt
}

func (m *Model) GetPrimaryKey() string {
	return AttributeID
}

// BeforeCreate sets new id.
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
