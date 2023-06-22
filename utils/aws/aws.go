package aws

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	appConfig "github.com/playground-pro-project/playground-pro-api/app/config"
)

const (
	AWS_S3_REGION = "ap-southeast-2"
	AWS_S3_BUCKET = "aws-pgp-bucket"
)

type AWSService struct {
	S3Client *s3.Client
}

func (awsSvc AWSService) UploadFile(key string, fileType string, file multipart.File) error {
	_, err := awsSvc.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(AWS_S3_BUCKET),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileType),
	})
	if err != nil {
		log.Println("Error while uploading the file", err)
	}

	return err
}

func ConfigS3(cfg *appConfig.AppConfig) AWSService {
	creds := credentials.NewStaticCredentialsProvider(
		cfg.AWS_ACCESS_KEY_ID, cfg.AWS_SECRET_ACCESS_KEY, "",
	)

	config, err := config.LoadDefaultConfig(
		context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(AWS_S3_REGION),
	)
	if err != nil {
		log.Println("Error while loading the aws config", err)
	}

	awsService := AWSService{
		S3Client: s3.NewFromConfig(config),
	}

	return awsService
}

func InitS3() AWSService {
	cfg := appConfig.InitConfig()
	return ConfigS3(cfg)
}
