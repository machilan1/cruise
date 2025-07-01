package fileapi

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/file"
	"github.com/machilan1/cruise/internal/framework/validate"
)

type AppSignedFile struct {
	Path      string `json:"path"`
	SignedURL string `json:"signedUrl"`
}

type AppFileInput struct {
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

func (app AppFileInput) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}

	return nil
}

type AppFile struct {
	ID        int       `json:"id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toAppFile(fl file.File) AppFile {
	return AppFile{
		ID:        fl.ID,
		Path:      fl.Path.String(),
		CreatedAt: fl.CreatedAt,
		UpdatedAt: fl.UpdatedAt,
	}
}

func toAppFiles(fls []file.File) []AppFile {
	result := make([]AppFile, len(fls))
	for i, fl := range fls {
		result[i] = toAppFile(fl)
	}
	return result
}
