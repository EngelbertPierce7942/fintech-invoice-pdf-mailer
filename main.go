package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/invoice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var in Invoice
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := issueAndSend(c, in); err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "sent", "invoice": in.Number})
	})
	log.Println("invoice service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
