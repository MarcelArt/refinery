package configs

import (
	"log"

	"github.com/MarcelArt/refinery/pkg/r2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func ConnectR2() *s3.Client {
	r2Client, err := r2.New(Env.R2AccessKeyID, Env.R2SecretKeyID, Env.R2AccountID)
	if err != nil {
		log.Fatalf("failed connecting to r2: %s", err.Error())
	}
	return r2Client
}
