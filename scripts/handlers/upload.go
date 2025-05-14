package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"example.com/trial1/utils"
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
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		http.Error(w, "Failed to generate encryption key", http.StatusInternalServerError)
		return
	}

	_, blobID, err := utils.UploadToWalrus(tempFilePath, epochs, key)
	_ = os.Remove(tempFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Upload failed: %v", err), http.StatusInternalServerError)
		return
	}

	encodedKey := base64.StdEncoding.EncodeToString(key)
	encodedBlobID := base64.StdEncoding.EncodeToString([]byte(blobID))

	json.NewEncoder(w).Encode(map[string]string{
		"blobID":        blobID,        // hex (for backend use)
		"blobID_base64": encodedBlobID, // base64 (for Move)
		"key":           encodedKey,    // base64
	})
}
