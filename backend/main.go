package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/trial1/handlers"
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
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)
	http.HandleFunc("/health", handlers.HealthHandler)

	fmt.Println("Server Started on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
