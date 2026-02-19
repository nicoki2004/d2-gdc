package repository

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
)

// WeaponRepository abstracts database operations for weapons
type WeaponRepository interface {
	UpsertWeapon(ctx context.Context, params db.UpsertWeaponParams) error
	ClearWeaponStats(ctx context.Context, instanceID string) error
	InsertWeaponStat(ctx context.Context, params db.InsertWeaponStatParams) error
	ClearWeaponPerks(ctx context.Context, instanceID string) error
	InsertWeaponPerk(ctx context.Context, params db.InsertWeaponPerkParams) error
	UpsertCharacter(ctx context.Context, params db.UpsertCharacterParams) error
	ClearAllWeaponStats(ctx context.Context) error
	ClearAllWeaponsPerks(ctx context.Context) error
	ClearAllWeaponsData(ctx context.Context) error
	RefreshData(ctx context.Context) error

	// Query methods for TUI
	GetAllWeapons(ctx context.Context) ([]db.Weapon, error)
	GetAllCharacters(ctx context.Context) ([]db.Character, error)
	GetWeaponsByType(ctx context.Context, weaponType string) ([]db.Weapon, error)
	GetWeaponsByName(ctx context.Context, pattern string) ([]db.Weapon, error)
	GetDuplicateWeapons(ctx context.Context) ([]db.Weapon, error)
	GetWeaponsByHash(ctx context.Context, hash int64) ([]db.Weapon, error)
	GetWeaponComparison(ctx context.Context, ids []string) ([]db.Weapon, error)
	GetWeaponStats(ctx context.Context, instanceID string) (map[string]int, error)
	GetWeaponPerks(ctx context.Context, instanceID string) ([]WeaponPerkRecord, error)
}

// WeaponPerkRecord represents a perk row loaded from DB.
type WeaponPerkRecord struct {
	Hash       int64
	Name       string
	IsEquipped bool
	SocketIdx  int
}

func (r *SQLCWeaponRepository) GetAllWeapons(ctx context.Context) ([]db.Weapon, error) {
	return r.queries.GetAllWeapons(ctx)
}

func (r *SQLCWeaponRepository) ClearAllWeaponStats(ctx context.Context) error {
	return r.queries.ClearAllWeaponStats(ctx)
}

func (r *SQLCWeaponRepository) ClearAllWeaponsData(ctx context.Context) error {
	return r.queries.ClearAllWeaponsData(ctx)
}

func (r *SQLCWeaponRepository) ClearAllWeaponsPerks(ctx context.Context) error {
	return r.queries.ClearAllWeaponsPerks(ctx)
}

func (r *SQLCWeaponRepository) UpsertCharacter(ctx context.Context, params db.UpsertCharacterParams) error {
	return r.queries.UpsertCharacter(ctx, params)
}

// SQLCWeaponRepository implementa WeaponRepository usando SQLC
type SQLCWeaponRepository struct {
	queries *db.Queries
	db      *sql.DB
}

// NewSQLCWeaponRepository crea una nueva instancia del repositorio
func NewSQLCWeaponRepository(queries *db.Queries, db *sql.DB) *SQLCWeaponRepository {
	return &SQLCWeaponRepository{
		queries: queries,
		db:      db,
	}
}

func (r *SQLCWeaponRepository) UpsertWeapon(ctx context.Context, params db.UpsertWeaponParams) error {
	return r.queries.UpsertWeapon(ctx, params)
}

func (r *SQLCWeaponRepository) ClearWeaponStats(ctx context.Context, instanceID string) error {
	return r.queries.ClearWeaponStats(ctx, instanceID)
}

func (r *SQLCWeaponRepository) InsertWeaponStat(ctx context.Context, params db.InsertWeaponStatParams) error {
	return r.queries.InsertWeaponStat(ctx, params)
}

func (r *SQLCWeaponRepository) ClearWeaponPerks(ctx context.Context, instanceID string) error {
	return r.queries.ClearWeaponPerks(ctx, instanceID)
}

func (r *SQLCWeaponRepository) InsertWeaponPerk(ctx context.Context, params db.InsertWeaponPerkParams) error {
	return r.queries.InsertWeaponPerk(ctx, params)
}

func (r *SQLCWeaponRepository) RefreshData(ctx context.Context) error {
	// 1. Iniciamos la transacción desde la conexión DB
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting tx: %w", err)
	}

	// 2. Automatic rollback if something fails (before Commit)
	defer tx.Rollback()

	// 3. Get a version of queries that uses this transaction
	qtx := r.queries.WithTx(tx)

	// 4. Execute queries in referential integrity order
	// First delete what has foreign keys
	if err := qtx.ClearAllWeaponsPerks(ctx); err != nil {
		return fmt.Errorf("perks clear error: %w", err)
	}

	if err := qtx.ClearAllWeaponStats(ctx); err != nil {
		return fmt.Errorf("stats clear error: %w", err)
	}

	if err := qtx.ClearAllWeaponsData(ctx); err != nil {
		return fmt.Errorf("weapons clear error: %w", err)
	}

	// 5. Consolidamos los cambios
	return tx.Commit()
}

// GetAllCharacters retrieves all characters from the database
func (r *SQLCWeaponRepository) GetAllCharacters(ctx context.Context) ([]db.Character, error) {
	return r.queries.GetAllCharacters(ctx)
}

// GetWeaponsByType retrieves weapons filtered by type
func (r *SQLCWeaponRepository) GetWeaponsByType(ctx context.Context, weaponType string) ([]db.Weapon, error) {
	return r.queries.GetWeaponsByType(ctx, weaponType)
}

// GetWeaponsByName retrieves weapons filtered by name pattern
func (r *SQLCWeaponRepository) GetWeaponsByName(ctx context.Context, pattern string) ([]db.Weapon, error) {
	return r.queries.GetWeaponsByName(ctx, pattern)
}

// GetDuplicateWeapons retrieves all duplicate weapons
func (r *SQLCWeaponRepository) GetDuplicateWeapons(ctx context.Context) ([]db.Weapon, error) {
	weapons, err := r.queries.GetAllWeapons(ctx)
	if err != nil {
		return nil, err
	}
	// Filter to only return weapons that have duplicates (same hash)
	hashMap := make(map[int64]int)
	for _, w := range weapons {
		hashMap[w.Hash]++
	}
	var duplicates []db.Weapon
	for _, w := range weapons {
		if hashMap[w.Hash] > 1 {
			duplicates = append(duplicates, w)
		}
	}
	return duplicates, nil
}

// GetWeaponsByHash retrieves all weapons with the same hash
func (r *SQLCWeaponRepository) GetWeaponsByHash(ctx context.Context, hash int64) ([]db.Weapon, error) {
	return r.queries.GetDuplicatesByHash(ctx, hash)
}

// GetWeaponComparison retrieves specific weapons for comparison
func (r *SQLCWeaponRepository) GetWeaponComparison(ctx context.Context, ids []string) ([]db.Weapon, error) {
	return r.queries.GetWeaponComparison(ctx, ids)
}

// GetWeaponStats retrieves all saved stats for a weapon instance.
func (r *SQLCWeaponRepository) GetWeaponStats(ctx context.Context, instanceID string) (map[string]int, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT stat_name, value
		FROM weapon_stats
		WHERE instance_id = ?
		ORDER BY stat_name ASC
	`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var name string
		var value int
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		stats[name] = value
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stats, nil
}

// GetWeaponPerks retrieves all saved perks (equipped and available) for a weapon instance.
func (r *SQLCWeaponRepository) GetWeaponPerks(ctx context.Context, instanceID string) ([]WeaponPerkRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT perk_hash, perk_name, is_equipped, socket_index
		FROM weapon_perks
		WHERE instance_id = ?
		ORDER BY socket_index ASC, is_equipped DESC, perk_name ASC
	`, instanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perks []WeaponPerkRecord
	for rows.Next() {
		var row WeaponPerkRecord
		if err := rows.Scan(&row.Hash, &row.Name, &row.IsEquipped, &row.SocketIdx); err != nil {
			return nil, err
		}
		perks = append(perks, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return perks, nil
}
