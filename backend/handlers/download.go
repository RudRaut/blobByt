package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"example.com/trial1/db" // Adjust to your actual module path
	"example.com/trial1/services"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // Or use "http://localhost:5173" for stricter control
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	blobID := r.URL.Query().Get("blobID")
	if blobID == "" {
		http.Error(w, "Missing blobID", http.StatusBadRequest)
		return
	}

	metadata, err := db.GetMetadataByBlobID(blobID)
	if err != nil {
		http.Error(w, "Metadata not found", http.StatusNotFound)
		return
	}

	// Decode encryption key from base64
	key, err := base64.StdEncoding.DecodeString(metadata.EncryptionKey)
	if err != nil {
		http.Error(w, "Invalid encryption key", http.StatusInternalServerError)
		return
	}

	// Download and decrypt the file from Walrus
	decryptedData, err := services.DownloadFromWalrus(blobID, key)
	if err != nil {
		http.Error(w, "Download failed", http.StatusInternalServerError)
		return
	}

	// Set headers and write response
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", metadata.Name))
	w.Header().Set("Content-Type", metadata.FileType)
	w.Write(decryptedData)
}
