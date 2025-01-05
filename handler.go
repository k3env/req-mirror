package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	servers, ok := r.URL.Query()["servers"]
	if !ok {
		w.Write([]byte("wrong parameters"))
		return
	}
	bb, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on getting req body, using empty")
		bb = []byte{}
	}
	log.Printf("~>: %s", strings.Join(strings.Fields(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(string(bb), "\n", " "), "\r", ""))), " "))
	for _, s := range servers {
		parsedUrl, err := ParseUrl(s)
		if err != nil {
			log.Printf("Error on parsing url: %s", err)
			continue
		}
		SendRequest(r, parsedUrl, bb)
	}
	w.Write([]byte("ok"))
}
