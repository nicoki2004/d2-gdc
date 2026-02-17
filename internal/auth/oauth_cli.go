package auth

import (
	"fmt"
	"log"

	"github.com/nicoki2004/g2-drc/internal/config"
)

func CmdAuthUrl(cfg config.Config) {
	url, err := AuthURL(cfg)
	if err != nil {
		log.Fatal("Error creating the URL")
	}
	fmt.Println(url)
}
