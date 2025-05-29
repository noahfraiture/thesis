package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// New BasicAuth middleware
func BasicAuth(db *DB, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		// Check if it starts with "Basic "
		if !strings.HasPrefix(auth, "Basic ") {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		// Extract and decode the credentials
		encoded := strings.TrimPrefix(auth, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		// Split into username and password
		creds := string(decoded)
		parts := strings.SplitN(creds, ":", 2)
		if len(parts) != 2 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}

		user := parts[0]
		password := parts[1]
		pass := base64.StdEncoding.EncodeToString([]byte(password))

		if !isAdmin(db, user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Admin Area"`)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Unauthorized")
			return
		}
		handler.ServeHTTP(w, r)
	})
}
