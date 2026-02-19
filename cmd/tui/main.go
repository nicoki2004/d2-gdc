package main

import (
	"context"
	"database/sql"
	"flag"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/destiny"
	"github.com/nicoki2004/g2-drc/internal/logger"
	"github.com/nicoki2004/g2-drc/internal/repository"
	"github.com/nicoki2004/g2-drc/internal/tui"
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
	refresh := flag.Bool("refresh", false, "Delete and reload data")
	flag.Parse()

	// 1. Open database file connection
	dbConn, err := sql.Open("sqlite3", getDatabasePath())
	if err != nil {
		log.Fatal("Failed to open database: %v", err)
	}
	defer dbConn.Close()
	// 2. Instantiate SQLC-generated queries
	queries := db.New(dbConn)

	// 3. Create repository
	repo := repository.NewSQLCWeaponRepository(queries, dbConn)

	log.Info("Database connected and ready.")

	if *refresh {
		log.Info("🧹 Cleaning up previous arsenal...")
		if err := repo.RefreshData(context.Background()); err != nil {
			log.Fatal("Refresh failed: %v", err)
		}
	}

	launchTUI(repo)
}

func launchTUI(repo repository.WeaponRepository) error {
	model := tui.NewModel(repo)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
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

// Initialize Config
func initConfig() *config.Config {
	cfg := config.Get()
	logger.GetLogger().Info("--------- Initialization Complete --------- ")
	return cfg
}
