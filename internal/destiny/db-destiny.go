package destiny

import (
	"context"
	"fmt"
	"strings"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
)

func SyncInventory(
	ctx context.Context,
	queries *db.Queries,
	profile *ProfileResponse,
	manifest map[string]ManifestItem,
) error {
	fmt.Println("🚀 Iniciando sincronización con la base de datos...")

	// 1. Recorrer Personajes: Equipado (205)
	for charID, equipment := range profile.Response.CharacterEquipment.Data {
		location := fmt.Sprintf("Equipped:%s", charID)
		for _, item := range equipment.Items {
			SaveItem(ctx, queries, item, profile, manifest, location)
		}
	}

	// 2. Recorrer Personajes: Inventario/Mochila (201)
	for charID, inventory := range profile.Response.CharacterInventories.Data {
		location := fmt.Sprintf("Inventory:%s", charID)
		for _, item := range inventory.Items {
			SaveItem(ctx, queries, item, profile, manifest, location)
		}
	}

	// 3. Recorrer el Vault (102)
	for _, item := range profile.Response.ProfileInventory.Data.Items {
		SaveItem(ctx, queries, item, profile, manifest, "Vault")
	}

	fmt.Println("✅ Sincronización finalizada.")
	return nil
}

func SaveItem(
	ctx context.Context,
	q *db.Queries,
	item Item,
	profile *ProfileResponse,
	manifest map[string]ManifestItem,
	location string,
) {
	hashStr := fmt.Sprintf("%d", item.ItemHash)
	def, ok := manifest[hashStr]

	if !ok || def.ItemType != 3 {
		return
	}

	// 1. Metadatos básicos
	level, kills, _ := GetWeaponMetadata(item.ItemInstanceId, profile)
	power := 0
	if inst, ok := profile.Response.ItemComponents.Instances.Data[item.ItemInstanceId]; ok {
		power = inst.PrimaryStat.Value
	}

	// 2. Upsert del arma
	q.UpsertWeapon(ctx, db.UpsertWeaponParams{
		InstanceID: item.ItemInstanceId,
		Hash:       int64(item.ItemHash),
		Name:       def.DisplayProperties.Name,
		Type:       def.ItemTypeDisplayName,
		Power:      int64(power),
		Kills:      int64(kills),
		Level:      int64(level),
		Location:   location,
	})

	// 3. Stats (Limpiar e Insertar)
	q.ClearWeaponStats(ctx, item.ItemInstanceId)
	if statsData, ok := profile.Response.ItemComponents.Stats.Data[item.ItemInstanceId]; ok {
		for hash, stat := range statsData.Stats {
			if statName, exists := StatsDictionary[hash]; exists {
				q.InsertWeaponStat(ctx, db.InsertWeaponStatParams{
					InstanceID: item.ItemInstanceId,
					StatName:   statName,
					Value:      int64(stat.Value),
				})
			}
		}
	}

	// 4. Perks (Aquí está el cambio clave)
	q.ClearWeaponPerks(ctx, item.ItemInstanceId)

	socketsData, hasSockets := profile.Response.ItemComponents.Sockets.Data[item.ItemInstanceId]
	reusableData, hasReusable := profile.Response.ItemComponents.ReusablePlugs.Data[item.ItemInstanceId]

	if hasSockets {
		for i, socket := range socketsData.Sockets {
			if socket.PlugHash == 0 {
				continue
			}

			// Verificamos si este socket tiene opciones múltiples (árbol de perks)
			socketIndexStr := fmt.Sprintf("%d", i)
			options, ok := reusableData.Plugs[socketIndexStr]

			if hasReusable && ok && len(options) > 0 {
				// Caso: El arma tiene múltiples opciones en este socket
				for _, opt := range options {
					perkDef, exists := manifest[fmt.Sprintf("%d", opt.PlugItemHash)]
					if exists && isActualPerk(perkDef.ItemTypeDisplayName) {
						q.InsertWeaponPerk(ctx, db.InsertWeaponPerkParams{
							InstanceID:  item.ItemInstanceId,
							PerkHash:    int64(opt.PlugItemHash),
							PerkName:    perkDef.DisplayProperties.Name,
							IsEquipped:  opt.PlugItemHash == socket.PlugHash, // Cambiado aquí
							SocketIndex: int64(i),
						})
					}
				}
			} else {
				// Caso: Socket simple (sin opciones extras, ej: el marco intrínseco)
				perkDef, exists := manifest[fmt.Sprintf("%d", socket.PlugHash)]
				if exists && isActualPerk(perkDef.ItemTypeDisplayName) {
					q.InsertWeaponPerk(ctx, db.InsertWeaponPerkParams{
						InstanceID:  item.ItemInstanceId,
						PerkHash:    int64(socket.PlugHash),
						PerkName:    perkDef.DisplayProperties.Name,
						IsEquipped:  true, // Cambiado aquí
						SocketIndex: int64(i),
					})
				}
			}
		}
	}
}

// Función auxiliar para filtrar basura (ornamentos, shaders, trackers)
func isActualPerk(typeName string) bool {
	valid := []string{"Intrinsic", "Trait", "Barrel", "Magazine", "Blade", "Guard", "Arrow", "String"}
	for _, v := range valid {
		if strings.Contains(typeName, v) {
			return true
		}
	}
	return false
}
