package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nicoki2004/g2-drc/db/sqlc"
	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/destiny"
)

func main() {
	// 1. Abrir conexión al archivo .db
	dbConn, err := sql.Open("sqlite3", "./arsenal.db")
	if err != nil {
		log.Fatal("No se pudo abrir la DB:", err)
	}
	defer dbConn.Close()

	// 2. Instanciar las queries generadas por SQLC
	queries := db.New(dbConn)

	fmt.Println("Base de datos conectada y lista.")

	// Get config form .env
	cfg := initConfig()

	apiClient := fetchUserData(cfg)
	apiClient.SetDbQuery(queries)

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

	destiny.SyncInventory(context.Background(), apiClient.DbQuery, dataProfile, manifest)

	// destiny.PrintCharacters(dataProfile)
	// destiny.PrintCharactersItems(dataProfile, manifest)
	// destiny.PrintAllWeapons(dataProfile, manifest)
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
