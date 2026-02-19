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
	fmt.Println("🚀 Starting database synchronization...")

	for charID, charData := range profile.Response.Characters.Data {
		err := saveCharacter(ctx, repo, charID, charData)
		if err != nil {
			logger.GetLogger().Error("Error saving Characters: %v", err)
		}
	}

	// 1. Iterate Characters: Equipped (205)
	for charID, equipment := range profile.Response.CharacterEquipment.Data {
		location := fmt.Sprintf("Equipped:%s", charID)
		for _, item := range equipment.Items {
			if err := SaveItem(ctx, repo, item, profile, manifest, location, charID); err != nil {
				return fmt.Errorf("error saving equipped item at %s: %w", location, err)
			}
		}
	}

	// 2. Iterate Characters: Inventory/Backpack (201)
	for charID, inventory := range profile.Response.CharacterInventories.Data {
		location := fmt.Sprintf("Inventory:%s", charID)
		for _, item := range inventory.Items {
			if err := SaveItem(ctx, repo, item, profile, manifest, location, charID); err != nil {
				return fmt.Errorf("error saving inventory item at %s: %w", location, err)
			}
		}
	}

	// 3. Iterate Vault (102)
	for _, item := range profile.Response.ProfileInventory.Data.Items {
		if err := SaveItem(ctx, repo, item, profile, manifest, "Vault", "Vault"); err != nil {
			return fmt.Errorf("error saving vault item: %w", err)
		}
	}

	fmt.Println("✅ Synchronization completed.")
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
		return nil // Not a weapon, not an error
	}

	extractor := NewWeaponExtractor()
	validator := NewPerkValidator()

	// 1. Save weapon base data
	metadata := extractor.ExtractMetadata(item, profile)
	if err := saveWeaponBase(ctx, repo, item, def, metadata, location, charID); err != nil {
		return fmt.Errorf("error saving weapon base data: %w", err)
	}

	// 2. Save stats
	if err := saveWeaponStats(ctx, repo, item, profile, extractor); err != nil {
		return fmt.Errorf("error saving weapon stats: %w", err)
	}

	// 3. Save perks
	if err := saveWeaponPerks(ctx, repo, item, profile, manifest, extractor, validator); err != nil {
		return fmt.Errorf("error saving weapon perks: %w", err)
	}

	return nil
}

// saveWeaponBase inserts or updates the weapon base
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
		Slot:        sql.NullString{String: GetSlotName(def.Inventory.BucketTypeHash), Valid: true},
		DamageType:  sql.NullString{String: GetDamageName(uint32(def.DefaultDamageTypeHash)), Valid: true},
		AmmoType:    sql.NullInt64{Int64: int64(def.EquippingBlock.AmmoType), Valid: true},
	})
}

// saveWeaponStats saves or updates the weapon stats
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

// saveWeaponPerks saves or updates the weapon perks
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

		// Case: Socket with multiple options (perks tree)
		if socketsData.HasReusable && ok && len(options) > 0 {
			if err := savePerkOptions(ctx, repo, item, options, socket, manifest, i, validator); err != nil {
				return err
			}
		} else {
			// Case: Simple socket (no options)
			if err := saveSinglePerk(ctx, repo, item, socket, manifest, i, validator); err != nil {
				return err
			}
		}
	}

	return nil
}

// savePerkOptions saves multiple perk options in a socket
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
		//
		// // ESTA LÍNEA ES LA CLAVE:
		//
		// if item.ItemInstanceId == "6917530146978745979" {
		// 	if item.ItemInstanceId == "TU_INSTANCE_ID_DE_MINT" || true {
		// 		fmt.Printf("🔍 Slot %d | Perk: %s | Tipo: '%s'\n",
		// 			socketIndex, perkDef.DisplayProperties.Name, perkDef.ItemTypeDisplayName)
		// 	}
		// }
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

// saveSinglePerk saves a single perk with no options
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
