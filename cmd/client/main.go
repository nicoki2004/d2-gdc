package main

import (
	"fmt"
	"log"

	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/destiny"
)

func main() {
	// Get config form .env
	cfg := config.Get()
	fmt.Println("--------- Initialization Complete --------- ")
	// auth.CmdAuthUrl(*cfg)

	token, err := auth.ReadTokenJSON()
	if err != nil {
		log.Fatal("Error, no token available")
	}

	api_client := destiny.NewClient(cfg, token)
	res, err := destiny.GetMembershipForCurrentUser(api_client)
	if err != nil {
		log.Fatal("Error from destiny")
	}

	fmt.Println("--------- Datos del Guardián ---------")
	for _, m := range res.Response.DestinyMemberships {
		platform := ""
		switch m.MembershipType {
		case 2:
			platform = "PlayStation"
		case 3:
			platform = "Steam"
		}
		fmt.Printf("Plataforma: %s | Usuario: %s | ID: %s\n", platform, m.DisplayName, m.MembershipId)
	}
}
