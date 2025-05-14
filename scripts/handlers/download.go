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

	encodedBlobID := r.URL.Query().Get("blobID")
	encodedKey := r.URL.Query().Get("key")
	if encodedBlobID == "" || encodedKey == "" {
		http.Error(w, "Missing blobID or key parameter", http.StatusBadRequest)
		return
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		http.Error(w, "Invalid key encoding", http.StatusBadRequest)
		return
	}

	blobIDBytes, err := base64.StdEncoding.DecodeString(encodedBlobID)
	if err != nil {
		http.Error(w, "Invalid blobID encoding", http.StatusBadRequest)
		return
	}
	blobID := string(blobIDBytes)

	// Direct client creation (aggregator endpoint)
	client := walrus.NewClient(
		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31415"}),
	)

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
