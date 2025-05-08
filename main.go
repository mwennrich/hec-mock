package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func hecHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "can't read body %w", http.StatusBadRequest)
		return
	}

	fmt.Println(string(body))

	resp := map[string]interface{}{
		"text":  "Success",
		"code":  0,
		"ackId": 1234,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "can't encode response %w", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/services/collector/event", hecHandler)
	mux.HandleFunc("/services/collector", hecHandler)

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
