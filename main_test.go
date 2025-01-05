package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var servers = map[string]string{
	"valid":   "https://webhook-test.com/a16b3f8aa3c1f1c1a64945eb66898a75",
	"invalid": "http://localhost:3251",
}
var sampleBody = `{"test": "ok"}`
var urls = map[string]string{
	"correct":   "aHR0cDovL2xvY2FsaG9zdDozMjUx",
	"incorrect": "invalid_b64",
	"non-url":   "c29tZXRleHQ=",
}

func TestHandler(t *testing.T) {
	h := &Handler{}

	t.Run("Missing servers parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:3250", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		if w.Body.String() != "wrong parameters" {
			t.Errorf("Expected response body 'wrong parameters', got '%s'", w.Body.String())
		}
	})

	t.Run("Valid servers parameter", func(t *testing.T) {
		serverUrl := servers["valid"]
		encodedServer := base64.StdEncoding.EncodeToString([]byte(serverUrl))
		reqBody := []byte(sampleBody)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:3250?servers="+encodedServer, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		if w.Body.String() != "ok" {
			t.Errorf("Expected response body 'ok', got '%s'", w.Body.String())
		}
	})

	t.Run("Invalid servers parameter", func(t *testing.T) {
		invalidServer := "invalid_base64"
		req := httptest.NewRequest(http.MethodPost, "http://localhost:3250?servers="+invalidServer, nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		if w.Body.String() != "ok" {
			t.Errorf("Expected response body 'ok', got '%s'", w.Body.String())
		}
	})

	t.Run("Unavailable server", func(t *testing.T) {
		encodedServer := base64.StdEncoding.EncodeToString([]byte(servers["invalid"]))
		reqBody := []byte(sampleBody)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:3250?servers="+encodedServer, bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
		}

		if w.Body.String() != "ok" {
			t.Errorf("Expected response body 'ok', got '%s'", w.Body.String())
		}
	})
}

func TestUrlParser(t *testing.T) {
	t.Run("Parse correct base64", func(t *testing.T) {
		res, err := ParseUrl(urls["correct"])
		if err != nil {
			t.Error(err)
			return
		}
		if res.String() != "http://localhost:3251" {
			t.Error("Wrong url returned")
			return
		}
	})

	t.Run("Parse incorrect base64", func(t *testing.T) {
		_, err := ParseUrl(urls["incorrect"])
		if err.Error() != "illegal base64 data at input byte 7" {
			t.Error(err)
			return
		}
	})

	t.Run("Parse non-url base64", func(t *testing.T) {
		_, err := ParseUrl(urls["non-url"])
		if err.Error() != "Invalid URL" {
			t.Error(err)
			return
		}
	})
}

func TestSendRequest(t *testing.T) {
	httpurl := &url.URL{
		Scheme: "https",
		Host:   "webhook-test.com",
		Path:   "/a16b3f8aa3c1f1c1a64945eb66898a75",
	}

	req := &http.Request{
		Method: "POST",
		URL:    httpurl,
		Header: map[string][]string{"Content-Type": {"application/json"}},
		Body:   io.NopCloser(bytes.NewBufferString(sampleBody)),
		Proto:  "HTTP/1.1",
	}

	t.Run("Sending request", func(t *testing.T) {
		SendRequest(req, httpurl, []byte(sampleBody))
	})
}
