package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

func ParseUrl(encodedUrl string) (*url.URL, error) {
	bd, err := base64.StdEncoding.DecodeString(encodedUrl)
	if err != nil {
		return nil, err
	}
	parsedUrl, err := url.Parse(string(bd))
	if err != nil {
		return nil, err
	}
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return nil, errors.New("Invalid URL")
	}
	return parsedUrl, nil
}

func SendRequest(r *http.Request, url *url.URL, body []byte) {
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
		return
	}
	log.Printf("<~ %s: status: %d", url, res.StatusCode)
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("<! %s: %s", url, err)
		return
	}
	log.Printf("<~ %s: %s", url, string(b))
}
