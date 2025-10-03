package aws

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
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

func GetAllFilesPath(folderPath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil

}

func CopyFinalDist(id string) {
	var folder_path = "./output/" + id + "/dist"
	allFiles, err := GetAllFilesPath(folder_path)

	if err != nil {
		log.Fatal(err)
	}
	for _, file := range allFiles {
		relPath, err := filepath.Rel(folder_path, file)
		if err != nil {
			log.Fatal(err)
		}
		key := filepath.ToSlash(filepath.Join("dist", id, relPath))

		UploadFile(key, file)
		fmt.Printf("File uploaded successfully: %s\n", key)
	}
	DeleteFolder("vercel-arch", id)
	listObjectsInABucket()
}

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
func listObjectsInABucket() {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		Region: "eu-north-1",
	}

	client := s3.NewFromConfig(cfg)
	data, err := client.ListObjects(context.Background(), &s3.ListObjectsInput{
		Bucket: aws.String("vercel-arch"),
	})
	if err != nil {
		log.Fatal(err)
	}
	// Iterate and print useful info
	for _, obj := range data.Contents {
		fmt.Printf("Name: %s, Size: %d, LastModified: %v\n",
			aws.ToString(obj.Key), obj.Size, obj.LastModified)
	}

}

func DeleteFolder(bucket, prefix string) {
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		Region: "eu-north-1",
	}

	client := s3.NewFromConfig(cfg)

	// List all objects with the prefix
	listResp, err := client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix + "/"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if len(listResp.Contents) == 0 {
		fmt.Println("No objects to delete")
		return
	}

	// Delete all objects
	for _, obj := range listResp.Contents {
		_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		})
		if err != nil {
			log.Printf("Failed to delete object %s: %v", *obj.Key, err)
		} else {
			fmt.Println("Deleted:", *obj.Key)
		}
	}

	fmt.Println("Folder deleted successfully:", prefix)
}
