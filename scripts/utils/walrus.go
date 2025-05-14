package utils

import (
	walrus "github.com/namihq/walrus-go"
	"github.com/namihq/walrus-go/encryption"
)

func UploadToWalrus(filePath string, epochs int, key []byte) (*walrus.StoreResponse, string, error) {

	client := walrus.NewClient(
		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31416"}),
	)

	// Upload file with encryption (GCM mode, IV handled by SDK)
	resp, err := client.StoreFile(filePath, &walrus.StoreOptions{
		Epochs: epochs,
		Encryption: &walrus.EncryptionOptions{
			Key:   key,
			Suite: encryption.AES256GCM, // Optional, defaults to this
		},
	})
	if err != nil {
		return nil, "", err
	}

	var blobID string
	if resp.NewlyCreated != nil {
		blobID = resp.NewlyCreated.BlobObject.BlobID
	} else if resp.AlreadyCertified != nil {
		blobID = resp.AlreadyCertified.BlobID
	}

	return resp, blobID, nil
}
