package main

import (
	"fmt"
	"net/http"

	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/logger"
)

func main() {
	log := logger.GetLogger()
	done := make(chan bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Fprintf(w, "Waiting for Bungie code...")
			return
		}

		log.Info("Code received!: %s", code)
		fmt.Fprintf(w, "Authorization successful! You can close this tab and return to the terminal.")

		cfg := config.Get()
		token, err := auth.ExchangeCode(cfg, code)
		if err != nil {
			log.Error("exchange error: %v", err)
			return
		}
		if err := auth.SaveTokenJSON(token); err != nil {
			log.Error("save keychain error: %v", err)
			return
		}
		log.Info("✅ Token saved securely.")

		done <- true
	})

	go func() {
		log.Info("Server listening on https://localhost:4200...")
		if err := http.ListenAndServeTLS(":4200", "localhost.pem", "localhost-key.pem", nil); err != nil {
			log.Fatal("Failed to start server: %v", err)
		}
	}()

	<-done
	fmt.Println("Closing authentication server. Ready to use the API!")
}
