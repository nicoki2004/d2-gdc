package tui

// WeaponData represents a weapon with all its information
type WeaponData struct {
	InstanceID  string
	Hash        int64
	Name        string
	Type        string
	Power       int
	Kills       int
	Level       int
	Location    string
	CharacterID string
	Tier        string
	IconUrl     string
	Slot        string
	DamageType  string
	Perks       []PerkData
	Stats       map[string]int
}

// PerkData represents a perk/trait on a weapon
type PerkData struct {
	Hash       uint32
	Name       string
	IsEquipped bool
	SocketIdx  int
	Category   string
}

// CharacterData represents a Destiny character
type CharacterData struct {
	CharacterID string
	ClassType   int
	LightLevel  int
	EmblemUrl   string
	LastPlayed  string
}

// SearchFilters contains filtering options for weapon search
type SearchFilters struct {
	WeaponType string
	Rarity     string
	Slot       string
	DamageType string
	MinPower   int
	MaxPower   int
	Location   string
	SearchText string
}

// GodRollFilter contains criteria for finding god rolls
type GodRollFilter struct {
	WeaponName    string
	DesiredPerks  []string
	MinStatValues map[string]int
	Rarity        string
}

// GodRollResult contains a weapon and its score
type GodRollResult struct {
	Weapon WeaponData
	Score  float64
}

// ComparisonWeapon wraps a weapon for comparison view
type ComparisonWeapon struct {
	Weapon          WeaponData
	IsSelected      bool
	MatchPercentage float64
}

// DuplicateGroup groups weapons by hash
type DuplicateGroup struct {
	Hash     int64
	Name     string
	Weapons  []WeaponData
	MaxPower int
}

// ViewType defines available views
type ViewType string

const (
	ViewHome       ViewType = "home"
	ViewList       ViewType = "list"
	ViewDetail     ViewType = "detail"
	ViewDuplicates ViewType = "duplicates"
	ViewComparison ViewType = "comparison"
	ViewGodRoll    ViewType = "godroll"
	ViewHelp       ViewType = "help"
)

// SortBy defines sorting options
type SortBy string

const (
	SortByName     SortBy = "name"
	SortByPower    SortBy = "power"
	SortByKills    SortBy = "kills"
	SortByType     SortBy = "type"
	SortByRariry   SortBy = "rarity"
	SortByLocation SortBy = "location"
)
