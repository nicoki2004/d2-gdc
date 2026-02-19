package destiny

import "slices"

// PerkValidator encapsulates perk type validation logic
type PerkValidator struct{}

// NewPerkValidator creates a new validator
func NewPerkValidator() *PerkValidator {
	return &PerkValidator{}
}

// IsActualPerk determines if a perk type should be displayed
func (pv *PerkValidator) IsActualPerk(typeName string) bool {
	validTypes := []string{
		"Intrinsic", "Enhanced Intrinsic", "Trait", "Enhanced Trait",
		"Origin Trait", "Barrel", "Launcher Barrel", "Magazine",
		"Battery", "Blade", "Guard", "Haft", "Arrow", "Bowstring",
		"String", "Stock", "Grip", "Enhanced Launcher Barrel", "Enhanced Magazine",
	}
	return slices.Contains(validTypes, typeName)
}

// IsIntrinsic determines if it's an intrinsic perk (weapon frame)
func (pv *PerkValidator) IsIntrinsic(typeName string) bool {
	return typeName == "Intrinsic"
}

// GetValidPerkTypes returns valid perk types
func (pv *PerkValidator) GetValidPerkTypes() []string {
	return []string{
		"Intrinsic", "Enhanced Intrinsic", "Trait", "Enhanced Trait",
		"Origin Trait", "Barrel", "Launcher Barrel", "Magazine",
		"Battery", "Blade", "Guard", "Haft", "Arrow", "Bowstring",
		"String", "Stock", "Grip", "Enhanced Launcher Barrel", "Enhanced Magazine",
	}
}
