package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/trial1/db"
)

// func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
// 	files, err := db.GetAllFiles()
// 	if err != nil {
// 		http.Error(w, "Failed to fetch files", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	json.NewEncoder(w).Encode(files)
// }

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	files, err := db.GetAllFiles()
	if err != nil {
		http.Error(w, "Failed to fetch files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(files); err != nil {
		http.Error(w, "Failed to encode files: "+err.Error(), http.StatusInternalServerError)
	}
}
