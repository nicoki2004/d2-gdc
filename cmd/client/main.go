package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/destiny"
	"github.com/nicoki2004/g2-drc/internal/logger"
	"github.com/nicoki2004/g2-drc/internal/repository"
)

// getDatabasePath returns the path to the SQLite database from environment variable DATABASE_FILE,
// with a default of "./arsenal.db" if not set.
func getDatabasePath() string {
	if path := os.Getenv("DATABASE_FILE"); path != "" {
		return path
	}
	return "./arsenal.db"
}

func main() {
	log := logger.GetLogger()
	refresh := flag.Bool("refresh", false, "Borra y recarga los datos")
	flag.Parse()

	// 1. Abrir conexión al archivo .db
	dbConn, err := sql.Open("sqlite3", getDatabasePath())
	if err != nil {
		log.Fatal("No se pudo abrir la DB: %v", err)
	}
	defer dbConn.Close()
	// 2. Instanciar las queries generadas por SQLC
	queries := db.New(dbConn)

	// 3. Crear el repositorio
	repo := repository.NewSQLCWeaponRepository(queries, dbConn)

	log.Info("Base de datos conectada y lista.")

	// Get config form .env
	cfg := initConfig()

	if *refresh {
		log.Info("🧹 Limpiando arsenal previo...")
		if err := repo.RefreshData(context.Background()); err != nil {
			log.Fatal("Fallo el refresh: %v", err)
		}
	}

	apiClient := fetchUserData(cfg)

	displayProfile(apiClient, repo)
}

// Shows the character and the items.
func displayProfile(apiClient *destiny.Client, repo repository.WeaponRepository) {
	log := logger.GetLogger()

	components := destiny.GetCoreComponents()
	dataProfile, err := destiny.GetProfile(apiClient, components...)
	if err != nil {
		log.Fatal("Error getting profile: %v", err)
	}
	manifest, err := destiny.LoadManifestMap("items_manifest.json")
	if err != nil {
		log.Error("Error loading manifest: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := destiny.SyncInventory(ctx, repo, dataProfile, manifest); err != nil {
		log.Fatal("Error sincronizando inventario: %v", err)
	}

	// destiny.PrintCharacters(dataProfile)
	// destiny.PrintCharactersItems(dataProfile, manifest)
	// destiny.PrintAllWeapons(dataProfile, manifest)
}

// Inicializa el ApiCLient con el Token y el MembershipId
func fetchUserData(cfg *config.Config) *destiny.Client {
	log := logger.GetLogger()

	token, err := auth.ReadTokenJSON()
	if err != nil {
		log.Fatal("Error, no token available: %v", err)
	}

	apiClient := destiny.NewClient(cfg, token)
	resp, err := destiny.GetMembershipForCurrentUser(apiClient)
	if err != nil {
		log.Fatal("Error from destiny: %v", err)
	}

	mId := resp.Response.PrimaryMembershipId
	apiClient.SetMembershipId(mId)

	return apiClient
}

// Inicializa el Config
func initConfig() *config.Config {
	cfg := config.Get()
	logger.GetLogger().Info("--------- Initialization Complete --------- ")
	return cfg
}
