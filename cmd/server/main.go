package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
)

func main() {
	done := make(chan bool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			fmt.Fprintf(w, "Esperando código de Bungie...")
			return
		}

		fmt.Printf("¡Código recibido!: %s\n", code)
		fmt.Fprintf(w, "¡Autorización exitosa! Puedes cerrar esta pestaña y volver a la terminal.")

		cfg := config.Get()
		token, err := auth.ExchangeCode(cfg, code)
		if err != nil {
			log.Printf("exchange error: %v", err)
			return
		}
		if err := auth.SaveTokenJSON(token); err != nil {
			log.Printf("save keychain error: %v", err)
			return
		}
		fmt.Println("✅ Token guardado de forma segura.")

		done <- true
	})

	go func() {
		log.Println("Servidor escuchando en https://localhost:4200...")
		if err := http.ListenAndServeTLS(":4200", "localhost.pem", "localhost-key.pem", nil); err != nil {
			log.Fatalf("No se pudo iniciar el servidor: %v", err)
		}
	}()

	<-done
	fmt.Println("Cerrando servidor de autenticación. ¡Listo para usar la API!")
}
