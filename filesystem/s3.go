package filesystem

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Filesystem is a wrapper around the S3 service
type S3Filesystem struct {
	bucket string
	client *s3.Client
}

// NewS3Filesystem returns a new S3Filesystem
func NewS3Filesystem(bucket string, region string) (*S3Filesystem, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &S3Filesystem{
		bucket: bucket,
		client: s3.NewFromConfig(cfg),
	}, nil
}

// GetFileContents performs GetObject on the S3 bucket and returns the contents as a byte array.
func (f *S3Filesystem) GetFileContents(path string) ([]byte, error) {
	resp, err := f.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &f.bucket,
		Key:    &path,
	})
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// WriteFileContents writes the contents to a file
func (f *S3Filesystem) WriteFileContents(path string, contents []byte) error {
	uploader := manager.NewUploader(f.client)
	_, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: &f.bucket,
		Key:    &path,
		Body:   bytes.NewReader(contents),
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *S3Filesystem) DeleteFile(path string) error {
	_, err := f.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &f.bucket,
		Key:    &path,
	})
	if err != nil {
		return err
	}
	return nil
}

func (f *S3Filesystem) ListDir(path string) ([]string, error) {
	resp, err := f.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: &f.bucket,
		Prefix: &path,
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, item := range resp.Contents {
		keys = append(keys, *item.Key)
	}
	return keys, nil
}
