package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/atotto/clipboard"
	qrterminal "github.com/mdp/qrterminal/v3"
)

type Payload struct {
	Text string `json:"text"`
}

func main() {
	port := 8000
	ip := getLocalIP()
	if ip == "" {
		log.Fatal("❌ Could not determine local IP address")
	}

	address := fmt.Sprintf("http://%s:%d", ip, port)
	fmt.Println("📡 Server available at:", address)
	fmt.Println("📱 Scan this QR code to connect:")

	// Show QR code in terminal
	qrterminal.GenerateHalfBlock(address, qrterminal.L)

	// Start HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var payload Payload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
				return
			}
			if err := clipboard.WriteAll(payload.Text); err != nil {
				http.Error(w, "❌ Failed to write to clipboard", http.StatusInternalServerError)
				return
			}
			log.Println("📋 Clipboard updated:", payload.Text)
			fmt.Fprintln(w, "✅ Clipboard updated")
		case http.MethodGet:
			fmt.Fprintln(w, "🖥️ Clipboard server is running.\nSend a POST with JSON: { \"text\": \"...\" }")
		default:
			http.Error(w, "❌ Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("🚀 Serving on %s\n", address)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// getLocalIP attempts to find a non-loopback local IP address.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok &&
			!ipNet.IP.IsLoopback() &&
			ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return ""
}