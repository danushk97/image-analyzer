package dtos

import "github.com/danushk97/image-analyzer/internal/image_metadata/model/v1"

// ImageMetadataResponse represents the response with additional fields
type ImageMetadataResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Filename       string `json:"filename"`
	FileType       string `json:"file_type"`
	FileSize       int64  `json:"file_size"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Status         string `json:"status"`
	AnalysisResult string `json:"analysis_result"`
	UploadURL      string `json:"upload_url"`
	DownloadURL    string `json:"download_url"`
}

// FromModel populates the ImageMetadataResponse from an ImageMetadata instance
func ImageMetadataResponseFromModel(image *model.ImageMetadata) *ImageMetadataResponse {
	r := &ImageMetadataResponse{}
	r.ID = image.GetPublicID()
	r.UserID = image.GetUserID()
	r.Filename = image.GetFilename()
	r.FileType = image.GetFileType()
	r.FileSize = image.GetFileSize()
	r.Width = image.Width
	r.Height = image.Height
	r.Status = image.GetStatus()
	r.AnalysisResult = image.GetAnalysisResult()

	return r
}
