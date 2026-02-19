package repository

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
)

// WeaponRepository abstrae las operaciones de base de datos para armas
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
}

func (r *SQLCWeaponRepository) ClearAllWeaponStats(ctx context.Context) error {
	return r.ClearAllWeaponStats(ctx)
}

func (r *SQLCWeaponRepository) ClearAllWeaponsData(ctx context.Context) error {
	return r.ClearAllWeaponsData(ctx)
}

func (r *SQLCWeaponRepository) ClearAllWeaponsPerks(ctx context.Context) error {
	return r.ClearAllWeaponsPerks(ctx)
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
		return fmt.Errorf("error al iniciar tx: %w", err)
	}

	// 2. Rollback automático si algo falla (antes del Commit)
	defer tx.Rollback()

	// 3. Obtenemos una versión de los queries que use esta transacción
	qtx := r.queries.WithTx(tx)

	// 4. Ejecutamos las queries en orden de integridad referencial
	// Borramos primero lo que tiene llaves foráneas
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
