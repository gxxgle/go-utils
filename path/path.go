package path

import (
	"os"
	"path/filepath"
)

// CurrentPath return the executable binary absolute path.
func CurrentPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// CurrentDir return the executable binary directory absolute path.
func CurrentDir() string {
	return filepath.Dir(CurrentPath())
}

// CurrentFilename return the executable binary filename.
func CurrentFilename() string {
	return filepath.Base(CurrentPath())
}

// RunDir return the running directory absolute path.
func RunDir() string {
	dir, _ := filepath.Abs("")
	return dir
}

// TopLevelDir return top-level directory of the directory absolute path.
func TopLevelDir(dir string) string {
	if dir == "/" {
		return dir
	}

	dir, _ = filepath.Abs(filepath.Join(dir, ".."))
	return dir
}

// Mkdir creates a directory named path.
func Mkdir(dir string) error {
	return os.MkdirAll(dir, 0777)
}
