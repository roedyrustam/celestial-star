package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"specialty-coffee-brewer/brewer"
)

func main() {
	// API endpoint
	http.HandleFunc("/api/score", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var recipe brewer.Recipe
		err := json.NewDecoder(r.Body).Decode(&recipe)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		score := brewer.CalculateScore(recipe)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(score)
	})

	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	port := ":8080"
	fmt.Printf("Mulai menyeduh! Server berjalan di http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server gagal dijalankan: %v", err)
	}
}
