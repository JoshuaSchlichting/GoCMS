package filesystem_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/joshuaschlichting/gocms/filesystem"
)

const bucketName = ""
const region = "us-east-1"

func init() {
	if bucketName == "" || region == "" {
		panic("This is not an automated test. Please set the bucketName and region variables in 's3_test.go'.")
	}
}

func TestS3Filesystem_GetFileContents(t *testing.T) {
	// create a new S3Filesystem
	fs, err := filesystem.NewS3Filesystem(bucketName, region)
	assert.NoError(t, err)

	// write some data to a file
	err = fs.WriteFileContents("myfile.txt", []byte("Hello, world!"))
	assert.NoError(t, err)

	// read the contents of the file
	contents, err := fs.GetFileContents("myfile.txt")
	assert.NoError(t, err)
	assert.Equal(t, []byte("Hello, world!"), contents)

	// delete the file
	err = fs.DeleteFile("myfile.txt")
	assert.NoError(t, err)
}

func TestS3Filesystem_ListDir(t *testing.T) {
	// create a new S3Filesystem
	fs, err := filesystem.NewS3Filesystem(bucketName, region)
	assert.NoError(t, err)

	// create some files in the directory
	err = fs.WriteFileContents("dir1/file1.txt", []byte("File 1"))
	assert.NoError(t, err)
	err = fs.WriteFileContents("dir1/file2.txt", []byte("File 2"))
	assert.NoError(t, err)
	err = fs.WriteFileContents("dir2/file3.txt", []byte("File 3"))
	assert.NoError(t, err)

	// list the files in the directory
	files, err := fs.ListDir("dir1/")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"dir1/file1.txt", "dir1/file2.txt"}, files)

	files, err = fs.ListDir("dir2/")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"dir2/file3.txt"}, files)

	// delete the files
	err = fs.DeleteFile("dir1/file1.txt")
	assert.NoError(t, err)
	err = fs.DeleteFile("dir1/file2.txt")
	assert.NoError(t, err)
	err = fs.DeleteFile("dir2/file3.txt")
	assert.NoError(t, err)
}
