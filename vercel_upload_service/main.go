package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"vercel_upload_service/aws"
	"vercel_upload_service/extractFilesPath"
	"vercel_upload_service/utils"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
)

type SendURlRequestBody struct {
	UrlString string `json:"url"`
}

var publisher = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", // Redis server address
	Password: "",               // no password set
	DB:       0,                // use default DB
})

func sendURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL SEND")
	if r.Method != http.MethodPost {
		http.Error(w, "ONLY POST Allowed", http.StatusMethodNotAllowed)
	}
	var body SendURlRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, "Invalid quest body", http.StatusBadRequest)
	}
	id := utils.Generate_random()
	destPath := fmt.Sprintf("output/%s", id)
	cmd := exec.Command("git", "clone", body.UrlString, destPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Git clone failed: %v\n%s", err, string(output))
	}
	folderPath := fmt.Sprintf("output/%s", id)

	extractedPath, err := extractFilesPath.GetAllFilesPath(folderPath)
	if err != nil {
		log.Fatalf("PATH extraction failed %v", err)
	}

	for _, file := range extractedPath {
		bucketKey, err := filepath.Rel("./output", file)
		if err != nil {
			log.Fatal("Error getting bucket key", err)
		}
		aws.UploadFile(bucketKey, file)
	}
	fmt.Println("FILE UPLOAD DONE")
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"id": id}
	ctx := context.Background()
	publisher.LPush(ctx, "build-queue", id)
	publisher.HSet(ctx, "status", id, "uploaded")

	json.NewEncoder(w).Encode(resp)
}

func status(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	response, err := publisher.HGet(ctx, "status", id).Result()
	if err == redis.Nil {
		response = "" // key does not exist
	} else if err != nil {
		http.Error(w, "Redis error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": response,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Printf("Vercel Deploy Service Started")
	mux := http.NewServeMux()
	mux.HandleFunc("/send", sendURL)
	mux.HandleFunc("/status", status)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
