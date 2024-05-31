package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	log.Fatalf("Error on starting web server: %s", http.ListenAndServe(":3250", &handler{}))
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		data, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			log.Printf("error on decoding: %s", err)
			continue
		}
		parsedUrl, err := url.Parse(string(data))
		if err != nil {
			log.Printf("error on url parsing: %s", err)
			continue
		}
		go sendRequest(r, parsedUrl, bb)
	}
	w.Write([]byte("ok"))
}

func sendRequest(r *http.Request, url *url.URL, body []byte) {
	client := http.Client{}
	rd := bytes.NewReader(body)
	rc := io.NopCloser(rd)
	headers := map[string][]string{"Content-Type": {r.Header.Get("Content-Type")}}
	req := &http.Request{
		Method:     r.Method,
		URL:        url,
		Proto:      r.Proto,
		ProtoMajor: r.ProtoMajor,
		ProtoMinor: r.ProtoMinor,
		Header:     headers,
		Body:       rc,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("<! %s: %s", url, err)
	}
	b, _ := io.ReadAll(res.Body)
	log.Printf("<~ %s: %s", url, string(b))
	log.Printf("<~ %s: status: %d", url, res.StatusCode)
}
