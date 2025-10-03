package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"vercel_deploy_service/utils"
)

type SendURlRequestBody struct {
	UrlString string `json:"url"`
}

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
	w.Write([]byte(id))
}

func main() {
	fmt.Printf("Vercel Deploy Service Started")
	http.HandleFunc("/send", sendURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
