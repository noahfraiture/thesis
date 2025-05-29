package main

import (
	"crypto/aes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

//go:embed black.jpg
var blackImage []byte

//go:embed room.jpg
var room []byte

//go:embed unlocked.jpg
var unlocked []byte

//go:embed locked.jpg
var locked []byte

var hasLight = false
var inRoom = true
var isLock = false

// Constants and variables
var (
	// XOR key as a byte slice
	xorKey = []byte{0x73, 0x8B, 0x55, 0x44}
	// AES key as a byte slice (16 bytes)
	aesKey = []byte{0x27, 0x99, 0x77, 0xf6, 0x2f, 0x6c, 0xfd, 0x2d, 0x91, 0xcd, 0x75, 0xb8, 0x89, 0xce, 0x0c, 0x9a}
)

// xore performs XOR operation on data with a cycling key
func xore(data []byte, key []byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[i] ^ key[i%len(key)]
	}
	return result
}

// padTo16 pads data with null bytes to a multiple of 16
func padTo16(data []byte) []byte {
	padLen := 16 - len(data)%16
	padding := make([]byte, padLen)
	return append(data, padding...)
}

// encryptData encrypts the plaintext using XOR and AES ECB mode, prefixing with 16 zero bytes
func encryptData(plaintext []byte) []byte {
	// Apply XOR with the key
	xorData := xore(plaintext, xorKey)
	// Pad to a multiple of 16 bytes
	paddedData := padTo16(xorData)
	// Create AES cipher
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err) // In production, handle errors gracefully
	}
	// Encrypt in ECB mode (block by block)
	ciphertext := make([]byte, len(paddedData))
	for i := 0; i < len(paddedData); i += 16 {
		block.Encrypt(ciphertext[i:i+16], paddedData[i:i+16])
	}
	// Prefix with 16 zero bytes
	header := make([]byte, 16)
	return append(header, ciphertext...)
}

// isValidAuth checks if the base64-encoded auth parameter starts with "admin:"
func isValidAuth(auth string) bool {
	if auth == "" {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return false
	}
	return strings.HasPrefix(string(decoded), "admin:")
}

// generateConfig generates a static configuration file as a byte slice
func generateConfig() []byte {
	return []byte(`Camera Configuration File

DeviceInfo:
Model: DS-2CD2143G0-I
SerialNumber: DS-2CD2143G0-I20190909AAWR123456789
FirmwareVersion: V5.4.0 Build 160414

Network:
IPAddress: 192.168.1.100
SubnetMask: 255.255.255.0
Gateway: 192.168.1.1
DNS1: 8.8.8.8
DNS2: 8.8.4.4

Users:
admin
password123
Role: Administrator

operator
operator123
Role: Operator

Video:
Resolution: 1920x1080
FrameRate: 30
Bitrate: 4096

Storage:
SDCard: Enabled
Capacity: 128GB
`)
}

// configurationFileHandler serves the encrypted configuration file
func configurationFileHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	if !isValidAuth(auth) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	plaintext := generateConfig()
	encrypted := encryptData(plaintext)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(encrypted)
}

// snapshotHandler serves the embedded black.png image
func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	if !isValidAuth(auth) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	if inRoom {
		if !hasLight {
			w.Write(blackImage)
		} else {
			w.Write(room)
		}
	} else {
		if isLock {
			w.Write(locked)
		} else {
			w.Write(unlocked)
		}
	}
}

func turnLight(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	if !isValidAuth(auth) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.URL.Query().Get("light") {
	case "on":
		hasLight = true
		w.Write([]byte{1})
	case "off":
		hasLight = false
		w.Write([]byte{0})
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func switchCamera(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	if !isValidAuth(auth) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	switch r.URL.Query().Get("location") {
	case "corridor":
		inRoom = false
		w.Write([]byte{0})
	case "room1":
		inRoom = true
		w.Write([]byte{1})
	}
}

func blockDoor(w http.ResponseWriter, r *http.Request) {
	auth := r.URL.Query().Get("auth")
	if !isValidAuth(auth) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	isLock = true
	w.Write([]byte{1})
}

// defaultPageHandler serves a Hikvision-like login page
func defaultPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("X-Powered-By", "Hikvision") // Additional Hikvision-specific header
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Hikvision Login</title>
		</head>
		<body>
			<div style="text-align: center; margin-top: 50px;">
				<h1>Hikvision Camera Login</h1>
				<p>Device Model: DS-2CD2143G0-I</p>
				<p>Firmware: V5.4.0 Build 160414</p>
				<form action="/login" method="post">
					<input type="text" name="username" placeholder="Username"><br>
					<input type="password" name="password" placeholder="Password"><br>
					<input type="submit" value="Login">
				</form>
			</div>
		</body>
		</html>
	`)
}

// customHandler sets a custom Server header for all responses
type customHandler struct {
	handler http.Handler
}

func (h *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Hikvision-Webs/1.0")
	h.handler.ServeHTTP(w, r)
}

func main() {
	// Register handlers
	http.HandleFunc("/System/configurationFile", configurationFileHandler)
	http.HandleFunc("/onvif-http/snapshot", snapshotHandler)
	http.HandleFunc("/light", turnLight)
	http.HandleFunc("/location", switchCamera)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	http.HandleFunc("/", defaultPageHandler)

	// Wrap the default handler with the custom handler to set the Server header
	handler := &customHandler{handler: http.DefaultServeMux}

	// Start the server
	fmt.Println("Server starting on port 8888...")
	err := http.ListenAndServe(":8888", handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
