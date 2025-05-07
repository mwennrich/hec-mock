package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
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
    json.NewEncoder(w).Encode(resp)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/services/collector/event", hecHandler)
    mux.HandleFunc("/services/collector", hecHandler)

    log.Fatal(http.ListenAndServe(":8080", mux))
}
