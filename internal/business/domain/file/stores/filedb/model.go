package filedb

import (
	"time"

	"github.com/machilan1/cruise/internal/business/domain/file"
)

type dbFile struct {
	ID        int       `db:"file_id"`
	Path      string    `db:"path"`
	Original  string    `db:"original_filename"`
	MimeType  string    `db:"mime_type"`
	Size      int       `db:"size"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func toDBFile(fl file.File) dbFile {
	return dbFile{
		ID:        fl.ID,
		Path:      fl.Path.String(),
		Original:  fl.Original,
		MimeType:  fl.MimeType,
		Size:      fl.Size,
		CreatedAt: fl.CreatedAt,
		UpdatedAt: fl.UpdatedAt,
	}
}

func toCoreFile(fl dbFile) (file.File, error) {
	p, err := file.ParsePath(fl.Path)
	if err != nil {
		return file.File{}, err
	}

	return file.File{
		ID:        fl.ID,
		Path:      p,
		Original:  fl.Original,
		MimeType:  fl.MimeType,
		Size:      fl.Size,
		CreatedAt: fl.CreatedAt,
		UpdatedAt: fl.UpdatedAt,
	}, nil
}

func toCoreFiles(fls []dbFile) ([]file.File, error) {
	result := make([]file.File, len(fls))
	for i, fl := range fls {
		r, err := toCoreFile(fl)
		if err != nil {
			return nil, err
		}
		result[i] = r
	}
	return result, nil
}
