package destiny

import "slices"

// PerkValidator encapsula la lógica de validación de tipos de perk
type PerkValidator struct{}

// NewPerkValidator crea un nuevo validador
func NewPerkValidator() *PerkValidator {
	return &PerkValidator{}
}

// IsActualPerk determina si un tipo de perk debe ser mostrado
func (pv *PerkValidator) IsActualPerk(typeName string) bool {
	validTypes := []string{"Intrinsic", "Trait", "Barrel", "Magazine", "Blade", "Guard", "Arrow", "String", "Stock", "Grip"}
	return slices.Contains(validTypes, typeName)
}

// IsIntrinsic determina si es un perk intrínseco (marco del arma)
func (pv *PerkValidator) IsIntrinsic(typeName string) bool {
	return typeName == "Intrinsic"
}

// GetValidPerkTypes devuelve los tipos de perk válidos
func (pv *PerkValidator) GetValidPerkTypes() []string {
	return []string{"Intrinsic", "Trait", "Barrel", "Magazine", "Blade", "Guard", "Arrow", "String", "Stock", "Grip"}
}
