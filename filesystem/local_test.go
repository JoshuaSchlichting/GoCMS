package filesystem

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestLocalFilesystem_GetFileContents(t *testing.T) {
	tempDir := path.Join(os.TempDir(), "go-filesystem-test")

	defer os.RemoveAll(tempDir)

	fs := NewLocalFilesystem(tempDir)

	// Write some data to a file
	err := fs.WriteFileContents("myfile.txt", []byte("Hello, world!"))
	if err != nil {
		t.Fatalf("Failed to write file contents: %v", err)
	}

	// Read the contents of the file
	contents, err := fs.GetFileContents("myfile.txt")
	if err != nil {
		t.Fatalf("Failed to read file contents: %v", err)
	}
	if string(contents) != "Hello, world!" {
		t.Fatalf("Expected file contents \"Hello, world!\", but got \"%s\"", contents)
	}

	// Delete the file
	err = fs.DeleteFile("myfile.txt")
	if err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}
}

func TestLocalFilesystem_ListDir(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some subdirectories and files
	testFiles := []string{"file1.txt", "file2.txt"}
	testDirs := []string{"dir1", "dir2"}
	for _, dir := range testDirs {
		if err := os.Mkdir(filepath.Join(tempDir, dir), os.ModePerm); err != nil {
			t.Fatalf("failed to create temp dir %s: %v", dir, err)
		}
		for _, file := range testFiles {
			if _, err := os.Create(filepath.Join(tempDir, dir, file)); err != nil {
				t.Fatalf("failed to create temp file %s/%s: %v", dir, file, err)
			}
		}
	}

	// Create the filesystem
	fs := NewLocalFilesystem(tempDir)

	// Test listing directories
	files, err := fs.ListDir("")
	if err != nil {
		t.Fatalf("ListDir returned an error: %v", err)
	}
	if len(files) != len(testDirs) {
		t.Fatalf("expected %d directories, got %d", len(testDirs), len(files))
	}
	for _, dir := range testDirs {
		found := false
		for _, file := range files {
			if file.Name() == dir {
				if !file.IsDir() {
					t.Errorf("%s is not a directory", dir)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("directory %s not found", dir)
		}
	}

	// Test listing files in a subdirectory
	files, err = fs.ListDir("dir1")
	if err != nil {
		t.Fatalf("ListDir returned an error: %v", err)
	}
	if len(files) != len(testFiles) {
		t.Fatalf("expected %d files, got %d", len(testFiles), len(files))
	}
	for _, file := range testFiles {
		found := false
		for _, f := range files {
			if f.Name() == file {
				if f.IsDir() {
					t.Errorf("%s is a directory", file)
				}
				found = true
				break
			}
		}
		if !found {
			t.Errorf("file %s not found", file)
		}
	}
}
