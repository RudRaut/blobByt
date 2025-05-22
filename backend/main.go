package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/trial1/db"
	"example.com/trial1/handlers"
	"github.com/joho/godotenv"
)

/*
Run the Aggregator Node
walrus aggregator --bind-address "127.0.0.1:31415" --metrics-address "127.0.0.1:27182"
*/

/*Run the Publisher Node
PUBLISHER_WALLETS_DIR=~/.config/walrus/publisher-wallets
walrus publisher   --bind-address "127.0.0.1:31416"   --sub-wallets-dir "/home/rud/.config/walrus/publisher-wallets"   --n-clients 1 --metrics-address "127.0.0.1:27183"
*/

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("URI")

	err = db.Connect(uri)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := db.Client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB client: %v", err)
		}
	}()

	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/files", handlers.ListFilesHandler)

	fmt.Println("Server Started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
