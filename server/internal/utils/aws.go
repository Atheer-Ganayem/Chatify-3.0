package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var s3Client *s3.Client

func InitAWS() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func UploadFileToS3(file []byte, fileHeader *multipart.FileHeader) (string, error) {
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		panic("AWS_BUCKET_NAME is a required env variable.")
	}
	fileName := fmt.Sprintf("chatify-3/%s%s", uuid.New().String(), strings.ReplaceAll(fileHeader.Filename, " ", ""))

	_, err := s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(file),
		ContentType: aws.String("image/webp"),
	})

	return fileName, err
}

func DeleteFile(filePath string) {
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	if bucketName == "" {
		panic("AWS_BUCKET_NAME is a required env variable.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(filePath),
	}

	_, err := s3Client.DeleteObject(ctx, input)
	if err != nil {
		log.Println(err)
	}
}
