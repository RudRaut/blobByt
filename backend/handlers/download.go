// package handlers

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"net/http"

// 	"example.com/trial1/db" // adjust import path to your project
// 	walrus "github.com/namihq/walrus-go"
// 	"github.com/namihq/walrus-go/encryption"
// )

// func DownloadHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")

// 	blobID := r.URL.Query().Get("blobID")
// 	if blobID == "" {
// 		http.Error(w, "Missing blobID parameter", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch file metadata from MongoDB using blobID
// 	metadata, err := db.FindByBlobID(blobID)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to fetch metadata: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	if metadata == nil {
// 		http.Error(w, "File metadata not found", http.StatusNotFound)
// 		return
// 	}

// 	// Decode the base64 encryption key stored in DB
// 	key, err := base64.StdEncoding.DecodeString(metadata.EncryptionKey)
// 	if err != nil {
// 		http.Error(w, "Failed to decode encryption key", http.StatusInternalServerError)
// 		return
// 	}

// 	// Create Walrus client
// 	client := walrus.NewClient(
// 		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31415"}),
// 	)

// 	// Fetch and decrypt the file from Walrus
// 	retrievedData, err := client.Read(blobID, &walrus.ReadOptions{
// 		Encryption: &walrus.EncryptionOptions{
// 			Key:   key,
// 			Suite: encryption.AES256GCM,
// 		},
// 	})
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Failed to retrieve data from Walrus: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	// Set headers with original filename if available
// 	filename := metadata.Name
// 	if filename == "" {
// 		filename = blobID
// 	}
// 	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
// 	w.Header().Set("Content-Type", "application/octet-stream")

// 	// Write the decrypted file bytes to response
// 	w.Write(retrievedData)
// }

package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"example.com/trial1/db" // Adjust to your actual module path
	"example.com/trial1/services"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
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
