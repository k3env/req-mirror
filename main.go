package main

import (
	"log"
	"net/http"
)

func main() {
	log.Fatalf("Error on starting web server: %s", http.ListenAndServe(":3250", &Handler{}))
}
