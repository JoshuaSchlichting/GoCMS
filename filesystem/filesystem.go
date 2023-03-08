package filesystem

import (
	"log"
	"os"
	"path/filepath"
)

// Filesystem is a wrapper around the os package
type Filesystem struct {
	path string
}

// NewLocalFilesystem returns a new Filesystem
func NewLocalFilesystem(rootDir string) *Filesystem {
	_, err := os.Stat(rootDir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(rootDir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	} else if err != nil {
		log.Println(err)
	}

	return &Filesystem{
		path: rootDir,
	}
}

// GetFileContents returns the contents of a file
func (f *Filesystem) GetFileContents(path string) ([]byte, error) {
	// Read the file
	file, err := os.Open(filepath.Join(f.path, path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Read the file
	bytes := make([]byte, stat.Size())
	_, err = file.Read(bytes)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// WriteFileContents writes the contents to a file
func (f *Filesystem) WriteFileContents(path string, contents []byte) error {
	// Write the file
	file, err := os.Create(filepath.Join(f.path, path))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(contents)
	if err != nil {
		return err
	}

	return nil
}

func (f *Filesystem) DeleteFile(path string) error {
	return os.Remove(filepath.Join(f.path, path))
}

func (f *Filesystem) ListDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(filepath.Join(f.path, path))
}
