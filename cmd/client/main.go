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
	cfg := initConfig()

	apiClient := fetchUserData(cfg)
	displayProfile(apiClient)
}

// Shows the character and the items.
func displayProfile(apiClient *destiny.Client) {
	components := destiny.GetCoreComponents()
	dataProfile, err := destiny.GetProfile(apiClient, components...)
	if err != nil {
		log.Fatal("Error getting profile")
	}
	manifest, err := destiny.LoadManifestMap("items_manifest.json")
	if err != nil {
		return
	}

	destiny.PrintCharacters(dataProfile)
	destiny.PrintCharactersItems(dataProfile, manifest)
}

// Inicializa el ApiCLient con el Token y el MembershipId
func fetchUserData(cfg *config.Config) *destiny.Client {
	token, err := auth.ReadTokenJSON()
	if err != nil {
		log.Fatal("Error, no token available")
	}

	apiClient := destiny.NewClient(cfg, token)
	resp, err := destiny.GetMembershipForCurrentUser(apiClient)
	if err != nil {
		log.Fatal("Error from destiny")
	}

	mId := resp.Response.PrimaryMembershipId
	apiClient.SetMembershipId(mId)

	return apiClient
}

// Inicializa el Config
func initConfig() *config.Config {
	cfg := config.Get()
	fmt.Println("--------- Initialization Complete --------- ")
	return cfg
}
