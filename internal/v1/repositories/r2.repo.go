package repositories

import (
	"context"
	"fmt"
	"io"

	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Repo struct {
	c *s3.Client
}

func NewR2Repo(c *s3.Client) *R2Repo {
	return &R2Repo{c: c}
}

func (r *R2Repo) UploadFile(ctx context.Context, bucket string, file io.Reader, objectKey string) (string, *s3.PutObjectOutput, error) {
	output, err := r.c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Body:   file,
		Key:    &objectKey,
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed uploading to r2: %w", err)
	}

	return fmt.Sprintf("https://%s/%s", configs.Env.R2PublicDomain, objectKey), output, nil
}
