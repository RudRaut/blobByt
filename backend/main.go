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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	http.Handle("/upload", enableCORS(http.HandlerFunc(handlers.UploadHandler)))
	http.Handle("/download", enableCORS(http.HandlerFunc(handlers.DownloadHandler)))
	http.Handle("/health", enableCORS(http.HandlerFunc(handlers.HealthHandler)))
	http.Handle("/files", enableCORS(http.HandlerFunc(handlers.ListFilesHandler)))

	fmt.Println("Server Started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
