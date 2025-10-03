package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func getAll(w http.ResponseWriter, r *http.Request) {
	hostname := r.URL.Hostname()
	id := strings.Split(hostname, ".")[0]
	filePath := r.URL.Path
	cfg := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Region:      "eu-north-1",
	}
	client := s3.NewFromConfig(cfg)
	contents, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("vercel-arch"),
		Key:    aws.String("dist/" + id + filePath),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Contents  are", contents)

}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/*", getAll)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)
	fmt.Println("VERCEL REQUEST HANDLER STARTED")
	log.Fatal(http.ListenAndServe(":3001", handler))

}
