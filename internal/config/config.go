package config

import (
	"log"
	"os"
	"sync"

	"github.com/subosito/gotenv"
)

type Config struct {
	ApiKey      string
	ClientID    string
	Secret      string
	RedirectURL string
}

var (
	instance *Config
	once     sync.Once
)

// Get retorna la instancia única de la configuración
func Get() *Config {
	once.Do(func() {
		err := gotenv.Load()
		if err != nil {
			log.Println("Aviso: No se encontró archivo .env, usando variables de sistema")
		}

		instance = &Config{
			ApiKey:      os.Getenv("BUNGIE_API_KEY"),
			ClientID:    os.Getenv("BUNGIE_OAUTH_CLIENT_ID"),
			Secret:      os.Getenv("BUNGIE_OAUTH_CLIENT_SECRET"),
			RedirectURL: os.Getenv("BUNGIE_OAUTH_REDIRECT_URI"),
		}

		// Validación rápida
		if instance.ApiKey == "" {
			log.Fatal("Error: BUNGIE_API_KEY es obligatoria")
		}
	})

	return instance
}
