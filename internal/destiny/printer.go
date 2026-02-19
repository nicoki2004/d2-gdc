package destiny

import (
	"fmt"
)

func PrintCharacters(profile *ProfileResponse) {
	fmt.Println("\n=== TUS GUARDIANES ===")
	for id, char := range profile.Response.Characters.Data {
		className := ""
		switch char.ClassType {
		case TitanClassType:
			className = "Titan"
		case HunterClassType:
			className = "Hunter"
		case WarlockClassType:
			className = "Warlock"
		}

		fmt.Printf("[%s] ID: %s | Luz: %d | Tiempo: %s min\n",
			className, id, char.Light, char.DateLastPlayed)
	}
}

func PrintAllWeapons(profile *ProfileResponse, manifest map[string]ManifestItem) {
	fmt.Println("\n=== RASTREO COMPLETO DE ARMAMENTO ===")
	// Every Character
	for charId, char := range profile.Response.Characters.Data {
		className := GetClassName(char.ClassType)
		fmt.Printf("\n--- %s (%s) ---\n", className, charId)
		// A. Lo equipado (205)
		if equip, ok := profile.Response.CharacterEquipment.Data[charId]; ok {
			for _, item := range equip.Items {
				renderItem(item, profile, manifest, "EQUIPADO")
			}
		}

		// B. Lo que lleva en la mochila (201)
		if inv, ok := profile.Response.CharacterInventories.Data[charId]; ok {
			for _, item := range inv.Items {
				renderItem(item, profile, manifest, "MOCHILA")
			}
		}
	}
	// 2. EL VAULT (102)
	fmt.Println("\n--- EL DEPÓSITO (VAULT) ---")
	for _, item := range profile.Response.ProfileInventory.Data.Items {
		renderItem(item, profile, manifest, "VAULT")
	}
}

func renderItem(item Item, profile *ProfileResponse, manifest map[string]ManifestItem, location string) {
	val, ok := manifest[fmt.Sprintf("%d", item.ItemHash)]
	if !ok || val.ItemType != Weapon {
		return
	}

	level, kills, progress := GetWeaponMetadata(item.ItemInstanceId, profile)
	power := 0
	if inst, ok := profile.Response.ItemComponents.Instances.Data[item.ItemInstanceId]; ok {
		power = inst.PrimaryStat.Value
	}
	// Si es del Vault, el InstanceId puede ser clave para buscar stats
	fmt.Printf("%s --> 🔫 %s [%s] @%d --> [Nivel %d - %.1f%%] [%d Bajas]\n", location,
		val.DisplayProperties.Name, val.ItemTypeDisplayName, power, level, progress, kills)
}

func PrintCharactersItems(profile *ProfileResponse, manifest map[string]ManifestItem) {
	fmt.Println("Cargando base de datos de items...")

	fmt.Println("\n=== ITEMS DE TUS GUARDIANES (MODO ÁRBOL) ===")
	for _, charEquipment := range profile.Response.CharacterEquipment.Data {
		for _, item := range charEquipment.Items {
			hashStr := fmt.Sprintf("%d", item.ItemHash)

			val, ok := manifest[hashStr]
			if !ok || val.ItemType != Weapon {
				continue
			}

			level, kills, progress := GetWeaponMetadata(item.ItemInstanceId, profile)
			power := 0
			if inst, ok := profile.Response.ItemComponents.Instances.Data[item.ItemInstanceId]; ok {
				power = inst.PrimaryStat.Value
			}
			// fmt.Printf("🔫 %s [%s] --> [Nivel %d - %.1f%%] [%d Bajas]\n", val.DisplayProperties.Name, val.ItemTypeDisplayName, level, progress, kills)

			fmt.Println(item.ItemInstanceId)
			fmt.Printf("🔫 %s [%s] @%d --> [Nivel %d - %.1f%%] [%d Bajas]\n",
				val.DisplayProperties.Name, val.ItemTypeDisplayName, power, level, progress, kills)
			PrintWeaponStats(item.ItemInstanceId, profile)

			if socketsData, exists := profile.Response.ItemComponents.Sockets.Data[item.ItemInstanceId]; exists {
				for i, socket := range socketsData.Sockets {
					if socket.PlugHash == 0 {
						continue
					}

					perkDef := manifest[fmt.Sprintf("%d", socket.PlugHash)]
					typeName := perkDef.ItemTypeDisplayName

					// 1. Special case: Intrinsic (The weapon frame)
					if typeName == "Intrinsic" {
						fmt.Printf("   ⭐ %s: %s\n", typeName, perkDef.DisplayProperties.Name)
						continue // Move to next socket
					}

					// 2. Case: Combat perks (Traits)
					validator := NewPerkValidator()
					if validator.IsActualPerk(typeName) {
						fmt.Printf("   └─ Socket %d (%s):\n", i, typeName)

						socketIndexStr := fmt.Sprintf("%d", i)
						reusableData, hasReusable := profile.Response.ItemComponents.ReusablePlugs.Data[item.ItemInstanceId]

						// Check if there are options in component 310 for this index
						options, ok := reusableData.Plugs[socketIndexStr]
						if hasReusable && ok && len(options) > 0 {
							for _, opt := range options {
								marker := "[ ]"
								if opt.PlugItemHash == socket.PlugHash {
									marker = "[X]"
								}
								optDef := manifest[fmt.Sprintf("%d", opt.PlugItemHash)]
								fmt.Printf("        %s %s\n", marker, optDef.DisplayProperties.Name)
							}
						} else {
							// Si no hay 310 o no hay opciones extras, mostramos el equipado
							fmt.Printf("        [X] %s\n", perkDef.DisplayProperties.Name)
						}
					}
				}
			}
		}
	}
}

func PrintWeaponStats(instanceId string, profile *ProfileResponse) {
	itemStats, ok := profile.Response.ItemComponents.Stats.Data[instanceId]
	if !ok {
		return
	}

	// Map to avoid duplicating stats
	printed := make(map[uint32]bool)

	// 1. DIM order for firearms (Commemoration/Drang)
	mainOrder := []uint32{4232813984, 4043523819, 1240592695, 155624089, 943549884, 4284893193, 3871231066}

	fmt.Print("   📊 Stats: ")
	for _, hash := range mainOrder {
		if stat, exists := itemStats.Stats[hash]; exists {
			fmt.Printf("%s %d | ", StatsDictionary[hash], stat.Value)
			printed[hash] = true
		}
	}

	// 2. Whatever's missing: This recovers Swing Speed, Charge Rate, etc.
	for hash, stat := range itemStats.Stats {
		if !printed[hash] {
			if name, found := StatsDictionary[hash]; found {
				fmt.Printf("%s %d | ", name, stat.Value)
			}
		}
	}
	fmt.Println()
}

func GetWeaponMetadata(instanceId string, profile *ProfileResponse) (level int, kills int, progress float64) {
	// 1. Try component 301 (Standard weapons)
	if objData, ok := profile.Response.ItemComponents.Objectives.Data[instanceId]; ok {
		for _, obj := range objData.Objectives {
			processObjective(obj, &kills, &progress)
		}
	}

	// 2. Try component 309 (Crafted/enhanced weapons like your Commemoration)
	if plugData, ok := profile.Response.ItemComponents.ItemPlugObjectives.Data[instanceId]; ok {
		for _, objectivesList := range plugData.ObjectivesPerPlug {
			for _, obj := range objectivesList {
				processObjective(obj, &kills, &progress)
			}
		}
	}
	return level, kills, progress
}

func processObjective(obj ObjectiveData, kills *int, progress *float64) {
	switch obj.ObjectiveHash {
	// Kills tracker (Generic PvE/PvP counter)
	case 73837075:
		if obj.Progress > 0 {
			*kills = obj.Progress
		}
	// Other progress trackers (can be level or killcounts)
	case 562334711, 867865505, 1970111194:
		if obj.CompletionValue > 0 && obj.CompletionValue > 1 {
			// If completionValue > 1, it's a tracker with multiple values
			*progress = (float64(obj.Progress) / float64(obj.CompletionValue)) * 100
		} else if obj.Progress > 100 {
			// If progress > 100, it's probably kills
			*kills = obj.Progress
		}
	}
}
