package main

import (
	// "crypto/aes"
	// "crypto/cipher"
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

	walrus "github.com/namihq/walrus-go"
	"github.com/namihq/walrus-go/encryption"
)

/*
Run the Aggregator Node
walrus aggregator --bind-address "127.0.0.1:31415" --metrics-address "127.0.0.1:27182"
*/

/*Run the Publisher Node
PUBLISHER_WALLETS_DIR=~/.config/walrus/publisher-wallets
walrus publisher   --bind-address "127.0.0.1:31416"   --sub-wallets-dir "/home/rud/.config/walrus/publisher-wallets"   --n-clients 1 --metrics-address "127.0.0.1:27183"
*/

// func UploadToWalrus(filePath string, epochs int, key []byte, iv []byte) (*walrus.StoreResponse, string, error) {
// 	client := walrus.NewClient(
// 		walrus.WithPublisherURLs([]string{"http://127.0.0.1:31416"}),
// 	)

// 	// Create temp file to store encrypted data before upload
// 	encryptedFilePath := filePath + ".enc"

// 	// Encrypt File
// 	if err := encryptFile(filePath, encryptedFilePath, key, iv); err != nil {
// 		return nil, "", fmt.Errorf("failed to encrypt file: %v", err)
// 	}

// 	// Upload the encrypted file
// 	resp, err := client.StoreFile(encryptedFilePath, &walrus.StoreOptions{
// 		Epochs: epochs,
// 		Encryption: &walrus.EncryptionOptions{
// 			Suite: encryption.AES256GCM,
// 			Key:   key,
// 		},
// 	})

// 	if err != nil {
// 		return nil, "", err
// 	}

// 	var blobID string
// 	if resp.NewlyCreated != nil {
// 		blobID = resp.NewlyCreated.BlobObject.BlobID
// 	} else if resp.AlreadyCertified != nil {
// 		blobID = resp.AlreadyCertified.BlobID
// 	}

// 	return resp, blobID, nil
// }

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

func DownloadFromWalrus(blobID, outputPath string) error {
	client := walrus.NewClient(
		walrus.WithAggregatorURLs([]string{"http://127.0.0.1:31415"}),
	)
	return client.ReadToFile(blobID, outputPath, nil)
}

// func GenerateAES256Key() ([]byte, error) {
// 	key := make([]byte, 32)
// 	if _, err := rand.Read(key); err != nil{
// 		return nil, fmt.Errorf("Error generating key: %w", err)
// 	}
// 	return key, nil
// }

// func GenerateIV() ([]byte, error) {
// 	iv := make([]byte, 12)
// 	if _, err := rand.Read(iv); err != nil {
// 		return nil, fmt.Errorf("Error generating IV: %w", err)
// 	}
// 	return iv, nil
// }

// func encryptFile(inputPath, outputPath string, key, iv []byte) error {
// 	inputFile, err := os.Open(inputPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open input file: %v", err)
// 	}
// 	defer inputFile.Close()

// 	outputFile, err := os.Create(outputPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return fmt.Errorf("failed to create cipher block: %v", err)
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return fmt.Errorf("failed to create GCM: %v", err)
// 	}

// 	buffer := make([]byte, 4096)
// 	for {
// 		n, err := inputFile.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return fmt.Errorf("error reading file: %v", err)
// 		}

// 		// Encrypt data
// 		ciphertext := gcm.Seal(nil, iv, buffer[:n], nil)

// 		// Write encrypted data
// 		if _, err := outputFile.Write(ciphertext); err != nil {
// 			return fmt.Errorf("failed to write encrypted data: %v", err)
// 		}
// 	}

// 	return nil
// }

// func decryptData(ciphertext, key, iv []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gcm.Open(nil, iv, ciphertext, nil)
// }

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
	_, blobID, err := UploadToWalrus(tempFilePath, epochs, key)
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

func main() {
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/download", DownloadHandler)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("Server Started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to Start server: %v", err)
	}
}
