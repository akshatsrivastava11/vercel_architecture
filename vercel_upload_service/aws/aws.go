package aws

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadFile(fileName string, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Region:      "eu-north-1",
	}
	client := s3.NewFromConfig(cfg)
	input := &s3.PutObjectInput{
		Bucket: aws.String("vercel-arch"),
		Key:    aws.String(fileName),
		Body:   file,
	}
	_, err = client.PutObject(context.Background(), input)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

	fmt.Println("File uploaded successfully:", fileName)

}
