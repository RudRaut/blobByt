package services

import (
	"crypto/rand"
	"fmt"

	walrus "github.com/namihq/walrus-go"
	"github.com/namihq/walrus-go/encryption"
)

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

func UploadToWalrus(filePath string, epochs int, key []byte) (string, error) {
	client := walrus.NewClient(
		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31416"}),
	)

	resp, err := client.StoreFile(filePath, &walrus.StoreOptions{
		Epochs: epochs,
		Encryption: &walrus.EncryptionOptions{
			Key:   key,
			Suite: encryption.AES256GCM,
		},
	})
	if err != nil {
		return "", err
	}

	if resp.NewlyCreated != nil {
		return resp.NewlyCreated.BlobObject.BlobID, nil
	}
	if resp.AlreadyCertified != nil {
		return resp.AlreadyCertified.BlobID, nil
	}
	return "", fmt.Errorf("Upload response empty")
}

func DownloadFromWalrus(blobID string, key []byte) ([]byte, error) {
	client := walrus.NewClient(
		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31415"}),
	)

	return client.Read(blobID, &walrus.ReadOptions{
		Encryption: &walrus.EncryptionOptions{
			Key:   key,
			Suite: encryption.AES256GCM,
		},
	})
}
