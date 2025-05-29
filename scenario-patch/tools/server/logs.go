package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

type LogHandler struct {
	handler http.Handler
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	sw := &statusWriter{ResponseWriter: w}
	var bodyLog string

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		bodyLog = "Error reading body"
	} else {
		bodyLog = string(bodyBytes)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	h.handler.ServeHTTP(sw, r)

	duration := time.Since(start)
	status := sw.status
	if status == 0 {
		status = http.StatusOK
	}

	log.Printf("%s - %s %s %d %s body: %s", r.RemoteAddr, r.Method, r.URL.Path, status, duration, bodyLog)
}
