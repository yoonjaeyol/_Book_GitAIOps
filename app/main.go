package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var counter int64
var version = "v0.2.2"

func main() {
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		podName = "unknown"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		id := atomic.AddInt64(&counter, 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":%d,"pod":"%s"}`, id, podName)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"version":"%s","pod":"%s"}`, version, podName)
	})

	log.Printf("Starting notiflex-api on :8080 (pod=%s)", podName)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}