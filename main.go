package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	file      *os.File
	fileMutex sync.RWMutex
	hecToken  string
)

func hecHandler(w http.ResponseWriter, r *http.Request) {
	if hecToken != "" && r.Header.Get("Authorization") != "Splunk "+hecToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, fmt.Sprintf("can't read body: %q", err), http.StatusBadRequest)
		return
	}
	if file != nil {
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, err := fmt.Fprintln(file, string(body))
		if err != nil {
			log.Println("Error writing to file:", err)
			http.Error(w, fmt.Sprintf("can't write to file: %q", err), http.StatusInternalServerError)
			return
		}
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
		http.Error(w, fmt.Sprintf("can't encode response: %q", err), http.StatusInternalServerError)
		return
	}
}

func main() {

	hecToken = os.Getenv("HEC_TOKEN")

	outputFile := os.Getenv("OUTPUT")
	if outputFile != "" {
		var err error
		file, err = os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Error opening file:", err)
			return
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
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
