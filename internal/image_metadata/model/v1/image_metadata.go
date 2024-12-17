package model

import (
	"fmt"
	"strings"

	"github.com/danushk97/image-analyzer/pkg/errors"
	"github.com/danushk97/image-analyzer/pkg/storage/sql"
)

const (
	// ImageMetadaIDPrefix ...
	ImageMetadaIDPrefix = "image_"

	EntityImageMetadata = "images_metadata"
)

// Image represents the image metadata table
type ImageMetadata struct {
	sql.Model             // Unique Image ID
	UserID         string `gorm:"type:uuid;not null" json:"user_id"`          // User ID (Uploader)
	Filename       string `gorm:"type:varchar(255);not null" json:"filename"` // Original Filename
	FileType       string `gorm:"type:varchar(50);not null" json:"file_type"` // File Type (e.g., JPEG, PNG)
	FileSize       int64  `gorm:"not null" json:"file_size"`                  // File Size in Bytes
	Width          int    `gorm:"not null" json:"width"`                      // Image Width in Pixels
	Height         int    `gorm:"not null" json:"height"`                     // Image Height in Pixels
	Status         string `gorm:"type:varchar(50);not null" json:"status"`    // State (e.g., INITIATED, COMPLETED)
	AnalysisResult string `gorm:"type:text" json:"analysis_result"`           // JSON Blob for Analysis Results
}

func NewImageMetadata() *ImageMetadata {
	return &ImageMetadata{}
}

// GetID retrieves the Image ID
func (i *ImageMetadata) GetID() string {
	return i.ID
}

// GetUserID retrieves the User ID
func (i *ImageMetadata) GetUserID() string {
	return i.UserID
}

// GetFilename retrieves the original filename
func (i *ImageMetadata) GetFilename() string {
	return i.Filename
}

// GetFileType retrieves the file type
func (i *ImageMetadata) GetFileType() string {
	return i.FileType
}

// GetFileSize retrieves the file size in bytes
func (i *ImageMetadata) GetFileSize() int64 {
	return i.FileSize
}

// GetDimensions retrieves the image dimensions in "Width x Height" format
func (i *ImageMetadata) GetDimensions() string {
	return fmt.Sprintf("%dx%d", i.Width, i.Height)
}

// GetState retrieves the state of the image
func (i *ImageMetadata) GetStatus() string {
	return i.Status
}

// GetAnalysisResult retrieves the analysis result
func (i *ImageMetadata) GetAnalysisResult() string {
	return i.AnalysisResult
}

// GetAnalysisResult retrieves the analysis result
func (i *ImageMetadata) TableName() string {
	return EntityImageMetadata
}

// GetPublicID returns public id of the offer
func (o *ImageMetadata) GetPublicID() string {
	return GetOImageMetdataIdWithPrefix(o.ID)
}

// GetPublicID returns public id of the offer
func (o *ImageMetadata) EntityName() string {
	return EntityImageMetadata
}

// GetPublicID returns public id of the offer
func (o *ImageMetadata) SetDefaults() errors.IError {
	return nil
}

// GetOfferIDWithPrefix adds the offer_id prefix if does not exist
func GetOImageMetdataIdWithPrefix(ID string) string {
	if strings.HasPrefix(ID, ImageMetadaIDPrefix) {
		return ID
	}
	return fmt.Sprintf("%s%s", ImageMetadaIDPrefix, ID)
}
