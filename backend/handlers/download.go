package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	walrus "github.com/namihq/walrus-go"
	"github.com/namihq/walrus-go/encryption"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	blobID := r.URL.Query().Get("blobID")
	keyBase64 := r.URL.Query().Get("key")
	if blobID == "" || keyBase64 == "" {
		http.Error(w, "Missing blobID or key parameter", http.StatusBadRequest)
		return
	}

	// Add the Logic later as per requirement
	// Fetch the key and nonce from frontend/Sui move module
	// Add logic accordingly

	// Decode the Base64 key
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		http.Error(w, "Invalid key encoding", http.StatusBadRequest)
		return
	}

	// aggregatorURL := "http://127.0.0.1:31415"

	// downloadURL := fmt.Sprintf("%s/v1/blobs/%s", aggregatorURL, blobID)

	// Create a Walrus client to read the file
	client := walrus.NewClient(
		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31415"}),
	)

	// Fetch the encrypted file form walrus
	retrievedData, err := client.Read(blobID, &walrus.ReadOptions{
		Encryption: &walrus.EncryptionOptions{
			Key:   key,
			Suite: encryption.AES256GCM,
		},
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve data from Walrus: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", blobID))
	w.Header().Set("Content-Type", "application/octet-stream")

	w.Write(retrievedData)
}
