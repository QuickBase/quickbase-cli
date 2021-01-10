package qbclient

import (
	"os"
	"strings"
)

// Filepath is a cross-platform solution for generating filepaths.
func Filepath(dirs ...string) string {
	return strings.Join(dirs, string(os.PathSeparator))
}

// FileExists returns true if filename exists and is a file.
func FileExists(filename string) bool {
	exists, isDir := exists(filename)
	return exists && !isDir
}

// DirExists returns true if dirname exists and is a directory.
func DirExists(dirname string) bool {
	exists, isDir := exists(dirname)
	return exists && isDir
}

func exists(name string) (exists bool, isDir bool) {
	info, err := os.Stat(name)
	exists = !os.IsNotExist(err)
	if err == nil {
		isDir = info.IsDir()
	}
	return
}
