package destiny

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
	"github.com/nicoki2004/g2-drc/internal/logger"
	"github.com/nicoki2004/g2-drc/internal/repository"
)

func SyncInventory(
	ctx context.Context,
	repo repository.WeaponRepository,
	profile *ProfileResponse,
	manifest map[string]ManifestItem,
) error {
	fmt.Println("🚀 Iniciando sincronización con la base de datos...")

	for charID, charData := range profile.Response.Characters.Data {
		err := saveCharacter(ctx, repo, charID, charData)
		if err != nil {
			logger.GetLogger().Error("Error saving Characters: %v", err)
		}
	}

	// 1. Recorrer Personajes: Equipado (205)
	for charID, equipment := range profile.Response.CharacterEquipment.Data {
		location := fmt.Sprintf("Equipped:%s", charID)
		for _, item := range equipment.Items {
			if err := SaveItem(ctx, repo, item, profile, manifest, location, charID); err != nil {
				return fmt.Errorf("error guardando item equipado en %s: %w", location, err)
			}
		}
	}

	// 2. Recorrer Personajes: Inventario/Mochila (201)
	for charID, inventory := range profile.Response.CharacterInventories.Data {
		location := fmt.Sprintf("Inventory:%s", charID)
		for _, item := range inventory.Items {
			if err := SaveItem(ctx, repo, item, profile, manifest, location, charID); err != nil {
				return fmt.Errorf("error guardando item inventario en %s: %w", location, err)
			}
		}
	}

	// 3. Recorrer el Vault (102)
	for _, item := range profile.Response.ProfileInventory.Data.Items {
		if err := SaveItem(ctx, repo, item, profile, manifest, "Vault", "Vault"); err != nil {
			return fmt.Errorf("error guardando item vault: %w", err)
		}
	}

	fmt.Println("✅ Sincronización finalizada.")
	return nil
}

func SaveItem(
	ctx context.Context,
	repo repository.WeaponRepository,
	item Item,
	profile *ProfileResponse,
	manifest map[string]ManifestItem,
	location string,
	charID string,
) error {
	// Validar que existe en manifest y es un arma
	hashStr := fmt.Sprintf("%d", item.ItemHash)
	def, ok := manifest[hashStr]
	if !ok || def.ItemType != 3 {
		return nil // No es un arma, no es error
	}

	extractor := NewWeaponExtractor()
	validator := NewPerkValidator()

	// 1. Guardar datos básicos del arma
	metadata := extractor.ExtractMetadata(item, profile)
	if err := saveWeaponBase(ctx, repo, item, def, metadata, location, charID); err != nil {
		return fmt.Errorf("error guardando datos base del arma: %w", err)
	}

	// 2. Guardar stats
	if err := saveWeaponStats(ctx, repo, item, profile, extractor); err != nil {
		return fmt.Errorf("error guardando stats del arma: %w", err)
	}

	// 3. Guardar perks
	if err := saveWeaponPerks(ctx, repo, item, profile, manifest, extractor, validator); err != nil {
		return fmt.Errorf("error guardando perks del arma: %w", err)
	}

	return nil
}

// saveWeaponBase inserta o actualiza el arma base
func saveWeaponBase(
	ctx context.Context,
	repo repository.WeaponRepository,
	item Item,
	def ManifestItem,
	metadata WeaponMetadata,
	location string,
	charID string,
) error {
	return repo.UpsertWeapon(ctx, db.UpsertWeaponParams{
		InstanceID: item.ItemInstanceId,
		Hash:       int64(item.ItemHash),
		Name:       def.DisplayProperties.Name,
		Type:       def.ItemTypeDisplayName,
		Power:      int64(metadata.Power),
		Kills:      int64(metadata.Kills),
		Level:      int64(metadata.Level),
		Location:   location,
		// NEW
		CharacterID: sql.NullString{String: charID, Valid: true},
		Tier:        sql.NullString{String: def.Inventory.TierTypeName, Valid: true},
		IconUrl:     sql.NullString{String: BungieCDN + def.DisplayProperties.Icon, Valid: true},
		Slot:        sql.NullString{String: GetSlotName(def.Inventory.BucketTypeHash)},
		DamageType:  sql.NullString{String: GetDamageName(uint32(def.DefaultDamageTypeHash)), Valid: true},
		AmmoType:    sql.NullInt64{Int64: int64(def.EquippingBlock.AmmoType), Valid: true},
	})
}

// saveWeaponStats guarda o actualiza los stats del arma
func saveWeaponStats(
	ctx context.Context,
	repo repository.WeaponRepository,
	item Item,
	profile *ProfileResponse,
	extractor *WeaponExtractor,
) error {
	if err := repo.ClearWeaponStats(ctx, item.ItemInstanceId); err != nil {
		return err
	}

	stats := extractor.ExtractStats(item, profile)
	for statName, value := range stats {
		if err := repo.InsertWeaponStat(ctx, db.InsertWeaponStatParams{
			InstanceID: item.ItemInstanceId,
			StatName:   statName,
			Value:      int64(value),
		}); err != nil {
			return err
		}
	}

	return nil
}

// saveWeaponPerks guarda o actualiza los perks del arma
func saveWeaponPerks(
	ctx context.Context,
	repo repository.WeaponRepository,
	item Item,
	profile *ProfileResponse,
	manifest map[string]ManifestItem,
	extractor *WeaponExtractor,
	validator *PerkValidator,
) error {
	if err := repo.ClearWeaponPerks(ctx, item.ItemInstanceId); err != nil {
		return err
	}

	socketsData := extractor.ExtractSockets(item, profile)
	if !socketsData.HasSockets {
		return nil
	}

	for i, socket := range socketsData.Sockets.Sockets {
		if socket.PlugHash == 0 {
			continue
		}

		socketIndexStr := fmt.Sprintf("%d", i)
		options, ok := socketsData.Plugs.Plugs[socketIndexStr]

		// Caso: Socket con múltiples opciones (árbol de perks)
		if socketsData.HasReusable && ok && len(options) > 0 {
			if err := savePerkOptions(ctx, repo, item, options, socket, manifest, i, validator); err != nil {
				return err
			}
		} else {
			// Caso: Socket simple (sin opciones)
			if err := saveSinglePerk(ctx, repo, item, socket, manifest, i, validator); err != nil {
				return err
			}
		}
	}

	return nil
}

// savePerkOptions guarda múltiples opciones de perks en un socket
func savePerkOptions(
	ctx context.Context,
	repo repository.WeaponRepository,
	item Item,
	options []PlugEntry,
	socket SocketEntry,
	manifest map[string]ManifestItem,
	socketIndex int,
	validator *PerkValidator,
) error {
	for _, opt := range options {
		perkDef, exists := manifest[fmt.Sprintf("%d", opt.PlugItemHash)]
		if !exists || !validator.IsActualPerk(perkDef.ItemTypeDisplayName) {
			continue
		}

		if err := repo.InsertWeaponPerk(ctx, db.InsertWeaponPerkParams{
			InstanceID:  item.ItemInstanceId,
			PerkHash:    int64(opt.PlugItemHash),
			PerkName:    perkDef.DisplayProperties.Name,
			IsEquipped:  opt.PlugItemHash == socket.PlugHash,
			SocketIndex: int64(socketIndex),
		}); err != nil {
			return err
		}
	}

	return nil
}

// saveSinglePerk guarda un único perk sin opciones
func saveSinglePerk(
	ctx context.Context,
	repo repository.WeaponRepository,
	item Item,
	socket SocketEntry,
	manifest map[string]ManifestItem,
	socketIndex int,
	validator *PerkValidator,
) error {
	perkDef, exists := manifest[fmt.Sprintf("%d", socket.PlugHash)]
	if !exists || !validator.IsActualPerk(perkDef.ItemTypeDisplayName) {
		return nil
	}

	return repo.InsertWeaponPerk(ctx, db.InsertWeaponPerkParams{
		InstanceID:  item.ItemInstanceId,
		PerkHash:    int64(socket.PlugHash),
		PerkName:    perkDef.DisplayProperties.Name,
		IsEquipped:  true,
		SocketIndex: int64(socketIndex),
	})
}

func saveCharacter(ctx context.Context, repo repository.WeaponRepository, charID string, charData CharacterData) error {
	return repo.UpsertCharacter(ctx, db.UpsertCharacterParams{
		CharacterID: charID,
		ClassType:   int64(charData.ClassType),
		LightLevel:  int64(charData.Light),
		EmblemUrl:   sql.NullString{String: BungieCDN + charData.EmblemPath, Valid: true},
		LastPlayed:  sql.NullTime{Time: time.Time(charData.DateLastPlayed), Valid: true},
	})
}

func refreshData(ctx context.Context, repo repository.WeaponRepository) error {
	return repo.RefreshData(ctx)
}
