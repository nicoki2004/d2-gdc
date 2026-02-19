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
			fmt.Fprintf(w, "Esperando código de Bungie...")
			return
		}

		log.Info("¡Código recibido!: %s", code)
		fmt.Fprintf(w, "¡Autorización exitosa! Puedes cerrar esta pestaña y volver a la terminal.")

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
		log.Info("✅ Token guardado de forma segura.")

		done <- true
	})

	go func() {
		log.Info("Servidor escuchando en https://localhost:4200...")
		if err := http.ListenAndServeTLS(":4200", "localhost.pem", "localhost-key.pem", nil); err != nil {
			log.Fatal("No se pudo iniciar el servidor: %v", err)
		}
	}()

	<-done
	fmt.Println("Cerrando servidor de autenticación. ¡Listo para usar la API!")
}
