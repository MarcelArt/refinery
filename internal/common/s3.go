package common

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type IS3Repo interface {
	UploadFile(ctx context.Context, bucket string, file io.Reader, objectKey string) (string, *s3.PutObjectOutput, error)
}
