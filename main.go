package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/atotto/clipboard"
	"github.com/grandcat/zeroconf"
)

type Payload struct {
	Text string `json:"text"`
}

func main() {
	port := 8000

	// Start HTTP clipboard server
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				var payload Payload
				err := json.NewDecoder(r.Body).Decode(&payload)
				if err != nil {
					http.Error(w, "❌ Invalid JSON", http.StatusBadRequest)
					return
				}
				err = clipboard.WriteAll(payload.Text)
				if err != nil {
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

		log.Printf("🚀 Serving on http://localhost:%d\n", port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}()

	// Register mDNS service
	host := "clipserver"
	service := "_http._tcp"
	domain := "local."

	server, err := zeroconf.Register(
		host,     // service instance name
		service,  // service type
		domain,   // domain
		port,     // service port
		[]string{"txtv=1", "lo=1", "clipboard=true"}, // optional metadata
		nil,      // use system's network interfaces
	)
	if err != nil {
		log.Fatal("❌ Failed to register mDNS:", err)
	}
	defer server.Shutdown()

	log.Println("🌐 Registered mDNS: http://" + host + ".local:" + fmt.Sprint(port))
	select {} // Keep running
}
