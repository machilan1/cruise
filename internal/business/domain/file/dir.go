package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DirRoot represents the root directory.
var DirRoot = Dir{}

// Dir represents a directory path, which does NOT start with
// a slash, and does NOT end with a slash.
// The root directory is represented by an empty Dir.
//
// For example, the following are valid Dir values:
//   - ""
//   - "dir"
//   - "dir/subdir"
//   - "dir/subdir/subsubdir"
type Dir struct {
	value string
}

func (d Dir) IsRoot() bool {
	return d.value == ""
}

func (d Dir) String() string {
	return d.value
}

func MarshalText(d Dir) ([]byte, error) {
	return []byte(d.String()), nil
}

func (d Dir) Equal(d2 Dir) bool {
	return d.value == d2.value
}

func (d Dir) IsZero() bool {
	return d.value == ""
}

// NewDir creates a new Dir.
// The value must be a valid directory path, which should NOT contain
// file names.
func NewDir(value string) (Dir, error) {
	if !strings.HasSuffix(value, "/") {
		value += "/"
	}

	return newDir(value)
}

func newDir(value string) (Dir, error) {
	// get the directory
	dir := filepath.Dir(value)

	// if dir start as "..", return error
	if strings.HasPrefix(dir, "..") {
		return Dir{}, fmt.Errorf("dir[%s]: directory with '..' is not allowed", value)
	}

	// if dir is "." or start with "/", remove the "." or "/" from the dir
	if dir == "." || strings.HasPrefix(dir, "/") {
		dir = dir[1:]
	}

	return Dir{value: dir}, nil
}

// ParseDirFromPath parses a path string into a Dir.
// The path will be treated as a file path, and the element after the last slash
// will be ignored.
func ParseDirFromPath(path string) (Dir, error) {
	return newDir(path)
}

func MustParseDirFromPath(path string) Dir {
	dir, err := ParseDirFromPath(path)
	if err != nil {
		panic(err)
	}
	return dir
}
