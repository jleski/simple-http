package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	version string = "0.1.0" // semver (major.minor.patch)
	port    string = "8080"  // port to listen
)

func main() {
	log.Printf("Simple HTTP Service v%s\n", version)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/api/headers", ApiGet)
	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Simple HTTP Service!")
}

func ApiGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Served request to GET /api/headers from %s.\n", r.RemoteAddr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(r.Header)
	if err != nil {
		log.Println("ERROR: Failed to Encode HTTP request headers to JSON!")
	}
}
