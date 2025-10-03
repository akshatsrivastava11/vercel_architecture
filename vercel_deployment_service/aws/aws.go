package aws

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DownloadS3Object(ctx context.Context, bucket, prefix, outputDir string) error {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Region:      "eu-north-1",
	}
	client := s3.NewFromConfig(cfg)

	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String("vercel-arch"),
		Prefix: aws.String(prefix),
	}
	paginator := s3.NewListObjectsV2Paginator(client, listInput)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Page is pageing", page)
		for _, object := range page.Contents {
			if object.Key == nil {
				continue
			}
			key := *object.Key
			fmt.Println("KEy is ", key)
			finalOutputPath := filepath.Join(outputDir, key)
			dirName := filepath.Dir(finalOutputPath)
			if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create dictionary %s : %w", dirName, err)
			}
			getInput := &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			}
			result, err := client.GetObject(ctx, getInput)
			if err != nil {
				fmt.Println(err)
			}
			defer result.Body.Close()
			outFile, err := os.Create(finalOutputPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", finalOutputPath, err)
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, result.Body); err != nil {
				return fmt.Errorf("failed to write file %s: %w", finalOutputPath, err)
			}

			fmt.Printf("Downloaded: %s\n", finalOutputPath)

		}
	}
	return nil

}
