package file

import "time"

type File struct {
	ID        int
	Path      Path
	Original  string
	MimeType  string
	Size      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SignedFile struct {
	File
	SignedURL string
}

type NewFile struct {
	// Path is the path where the file is stored.
	Path Path
	// Original is the original basename of the file.
	// Including both the filename and extension.
	Original string
	// Size is the size of the file in bytes.
	Size int
}
