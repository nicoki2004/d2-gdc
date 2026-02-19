package tui

import (
	db "github.com/nicoki2004/g2-drc/db/sqlc"
)

// convertToWeaponData converts a slice of db.Weapon to WeaponData
func convertToWeaponData(weapons []db.Weapon) []WeaponData {
	result := make([]WeaponData, 0, len(weapons))
	for _, w := range weapons {
		result = append(result, convertSingleWeapon(w))
	}
	return result
}

// convertSingleWeapon converts a single db.Weapon to WeaponData
func convertSingleWeapon(w db.Weapon) WeaponData {
	characterID := ""
	if w.CharacterID.Valid {
		characterID = w.CharacterID.String
	}

	return WeaponData{
		InstanceID:  w.InstanceID,
		Hash:        w.Hash,
		Name:        w.Name,
		Type:        w.Type,
		Power:       int(w.Power),
		Kills:       int(w.Kills),
		Level:       int(w.Level),
		Location:    w.Location,
		CharacterID: characterID,
		Tier:        w.Tier.String,
		IconUrl:     w.IconUrl.String,
		Slot:        w.Slot.String,
		DamageType:  w.DamageType.String,
		Perks:       []PerkData{},         // TODO: Load perks separately
		Stats:       make(map[string]int), // TODO: Load stats separately
	}
}

// convertToCharacterData converts a slice of db.Character to CharacterData
func convertToCharacterData(characters []db.Character) []CharacterData {
	result := make([]CharacterData, 0, len(characters))
	for _, c := range characters {
		result = append(result, convertSingleCharacter(c))
	}
	return result
}

// convertSingleCharacter converts a single db.Character to CharacterData
func convertSingleCharacter(c db.Character) CharacterData {
	return CharacterData{
		CharacterID: c.CharacterID,
		ClassType:   int(c.ClassType),
		LightLevel:  int(c.LightLevel),
		EmblemUrl:   c.EmblemUrl.String,
		LastPlayed:  c.LastPlayed.Time.String(),
	}
}

// convertWeaponsToViewData is an alias for convertToWeaponData for views package
func convertWeaponsToViewData(weapons []WeaponData) []WeaponData {
	return weapons
}
