package file

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/machilan1/cruise/internal/business/sdk/mimetype"
)

type Path struct {
	dir      Dir
	filename string
	ext      string
}

func (p Path) Dir() string {
	return p.dir.String()
}

// Base returns the base of the path, which is the filename with the extension.
func (p Path) Base() string {
	return p.filename + "." + p.ext
}

// Filename returns the filename of the path, without the extension.
func (p Path) Filename() string {
	return p.filename
}

func (p Path) Ext() string {
	return p.ext
}

func (p Path) Mime() string {
	return mimetype.DetectFileExt(p.ext)
}

// String returns the path as a string.
func (p Path) String() string {
	suffix := ""
	if p.ext != "" {
		suffix = "." + p.ext
	}

	if p.dir.IsRoot() {
		return fmt.Sprintf("%s%s", p.filename, suffix)
	}

	return fmt.Sprintf("%s/%s%s", p.dir, p.filename, suffix)
}

func (p Path) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p Path) Equal(p2 Path) bool {
	return p.dir.Equal(p2.dir) && p.filename == p2.filename && p.ext == p2.ext
}

func (p Path) IsZero() bool {
	return p.dir.IsZero() && p.filename == "" && p.ext == ""
}

// NewPath creates a new Path.
func NewPath(dir Dir, filename, ext string) (Path, error) {
	if filename == "" {
		return Path{}, fmt.Errorf("invalid filename[%s]: must not be empty", filename)
	}

	// check if the base is a directory
	if filepath.Base(filename) == filepath.Dir(filename) {
		return Path{}, fmt.Errorf("invalid filename[%s]: must not end with a slash", filename)
	}

	// remove the dot from the extension, if any
	if ext != "" && ext[0] == '.' {
		ext = ext[1:]
	}

	return Path{
		dir:      dir,
		filename: filename,
		ext:      ext,
	}, nil
}

// NewRandomPath creates a new Path with a random filename.
// We're using random filenames to avoid conflicts when uploading files.
func NewRandomPath(dir Dir, ext string) (Path, error) {
	return NewPath(dir, uuid.New().String(), ext)
}

// ParsePath parses a path string into a Path.
func ParsePath(path string) (Path, error) {
	dir, err := ParseDirFromPath(path)
	if err != nil {
		return Path{}, err
	}

	base := filepath.Base(path)
	ext := filepath.Ext(base)
	filename := base[:len(base)-len(ext)]

	return NewPath(dir, filename, ext)
}

func MustParsePath(path string) Path {
	p, err := ParsePath(path)
	if err != nil {
		panic(err)
	}
	return p
}
