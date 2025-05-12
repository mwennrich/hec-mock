package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var file *os.File

func hecHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "can't read body %w", http.StatusBadRequest)
		return
	}
	if file != nil {
		fmt.Fprintln(file, string(body))
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

	outputFile := os.Getenv("OUTPUT")
	if outputFile == "" {
		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Error opening file:", err)
			return
		}
		defer file.Close()
	}

	http.HandleFunc("/services/collector/event", hecHandler)
	http.HandleFunc("/services/collector", hecHandler)

	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 1 * time.Minute,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
