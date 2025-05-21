package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"example.com/trial1/services"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             // Allow cross-origin requests
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS") // Allow PUT method
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow Content-Type header
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests (CORS OPTIONS method)
	if r.Method == http.MethodOptions {
		return
	}

	// Parse the form data (file + optional form fields)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to Parse form", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the form
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is Missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Temporary file path to save the uploaded file
	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, handler.Filename)

	// Create the temporary file
	out, err := os.Create(tempFilePath)
	if err != nil {
		http.Error(w, "Failed to create Temp file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Save the uploaded file to the temporary file
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Failed to save temp file", http.StatusInternalServerError)
		return
	}

	// Get the optional "epochs" field from the form
	epochs := 3
	if r.FormValue("epochs") != "" {
		if parsed, err := strconv.Atoi(r.FormValue("epochs")); err == nil {
			epochs = parsed
		}
	}

	// Generate key (32 bytes for AES-256)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		http.Error(w, "Failed to generate encryption key", http.StatusInternalServerError)
		return
	}

	// Upload the file to Walrus (assumes UploadToWalrus function is correct and can accept encryption key)
	blobID, err := services.UploadToWalrus(tempFilePath, epochs, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload Failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Clean up the temporary file after upload
	_ = os.Remove(tempFilePath)

	// Encode the encryption key to Base64 for the response
	encodedKey := base64.StdEncoding.EncodeToString(key)

	// Respond with JSON containing the blobID and encoded encryption key
	err = json.NewEncoder(w).Encode(map[string]string{
		"blobID": blobID,
		"key":    encodedKey,
	})

	if err != nil {
		log.Println("JSON encode error:", err)
	}
}
