package destiny

import (
	"time"
)

// MembershipResponse contains Bungie subscription information for the user
// Returns all Destiny accounts connected to the Bungie account
type MembershipResponse struct {
	Response struct {
		DestinyMemberships  []DestinyMembership `json:"destinyMemberships"`
		PrimaryMembershipId string              `json:"primaryMembershipId"`
	} `json:"Response"`
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// DestinyMembership represents an individual gaming platform (PS5, Xbox, Steam, etc)
type DestinyMembership struct {
	MembershipType          int    `json:"membershipType"`
	MembershipId            string `json:"membershipId"`
	DisplayName             string `json:"displayName"`
	BungieGlobalDisplayName string `json:"bungieGlobalDisplayName"`
}

type MembershipType int

const (
	PlayStation = 2
	Steam       = 3
)

// ProfileResponse is the complete response from Bungie's Profile endpoint
// Contains characters, inventory, equipment, stats and all technical details
type ProfileResponse struct {
	Response struct {
		// 100 - PROFILE INFORMATION
		Characters struct {
			Data map[string]CharacterData `json:"data"`
		} `json:"characters"`

		// 102 - THE VAULT (STORAGE)
		ProfileInventory struct {
			Data struct {
				Items []Item `json:"items"`
			} `json:"data"`
		} `json:"profileInventory"`

		// 201 - THE BACKPACKS (INVENTORIES)
		CharacterInventories struct {
			Data map[string]CharacterEquipmentData `json:"data"`
		} `json:"characterInventories"`
		// 205 - EQUIPPED ITEMS
		CharacterEquipment struct {
			Data map[string]CharacterEquipmentData `json:"data"`
		} `json:"characterEquipment"`

		// 300-309 - TECHNICAL DETAILS
		ItemComponents struct {
			Instances struct {
				Data map[string]ItemInstanceData `json:"data"`
			} `json:"instances"`

			Objectives struct {
				Data map[string]ItemObjectivesComponent `json:"data"`
			} `json:"objectives"`

			Stats struct {
				Data map[string]ItemStatsComponent `json:"data"`
			} `json:"stats"`

			// COMPONENT 309: The one you need for the kills tracker
			ItemPlugObjectives struct {
				Data map[string]ItemPlugObjectivesComponent `json:"data"`
			} `json:"plugObjectives"`
			// Other optional but recommended
			Sockets struct {
				Data map[string]ItemSocketsComponent `json:"data"`
			} `json:"sockets"`

			ReusablePlugs struct {
				Data map[string]ItemReusablePlugsComponent `json:"data"`
			} `json:"reusablePlugs"`
		} `json:"itemComponents"`
	} `json:"Response"`
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// Auxiliary struct for 309
type ItemPlugObjectivesComponent struct {
	// CHANGE: Bungie calls it "objectivesPerPlug" internally for each socket
	ObjectivesPerPlug map[string][]ObjectiveData `json:"objectivesPerPlug"`
}

// This one stays clean, only with the sockets
type ItemSocketsComponent struct {
	Sockets []SocketEntry `json:"sockets"`
}

// These remain the same
type ItemReusablePlugsComponent struct {
	Plugs map[string][]PlugEntry `json:"plugs"`
}

type PlugEntry struct {
	PlugItemHash uint32 `json:"plugItemHash"`
	CanInsert    bool   `json:"canInsert"`
}

type SocketEntry struct {
	PlugHash  uint32 `json:"plugHash"` // The Hash of the Perk or Modifier
	IsEnabled bool   `json:"isEnabled"`
	IsVisible bool   `json:"isVisible"`
}

// CharacterData represents a user character (Component 200)
type CharacterData struct {
	CharacterId    string    `json:"characterId"`
	ClassType      int       `json:"classType"`
	Light          int       `json:"light"`
	DateLastPlayed time.Time `json:"dateLastPlayed"`
	EmblemPath     string    `json:"emblemBackgroundPath"`
}

type CharacterEquipmentData struct {
	Items []Item `json:"items"`
}

// Item representa un elemento en el inventario (arma, armadura, consumible, etc)
type Item struct {
	ItemHash       uint32 `json:"itemHash"`       // Definition hash for the item (reference to manifest)
	ItemInstanceId string `json:"itemInstanceId"` // Unique ID for this specific instance
}

// ItemStatsComponent represents the stats of an instance (Component 304)
type ItemStatsComponent struct {
	Stats map[uint32]StatData `json:"stats"`
}

// StatData contains the individual stat value
type StatData struct {
	StatHash uint32 `json:"statHash"`
	Value    int    `json:"value"`
}

// For Power (Component 300)
type ItemInstanceData struct {
	PrimaryStat struct {
		Value int `json:"value"` // Power lives here (460)
	} `json:"primaryStat"`
}

// For Kills and Level (Component 301)
type ItemObjectivesComponent struct {
	Objectives []ObjectiveData `json:"objectives"`
}

type ObjectiveData struct {
	ObjectiveHash   uint32 `json:"objectiveHash"`   // MUST be camelCase
	Progress        int    `json:"progress"`        // MUST be lowercase
	CompletionValue int    `json:"completionValue"` // MUST be camelCase
	Complete        bool   `json:"complete"`
	Visible         bool   `json:"visible"`
}

const (
	// Profile (100-199)
	ProfilesComponent           = "100"
	VendorReceiptsComponent     = "101"
	ProfileInventoriesComponent = "102"
	ProfileCurrenciesComponent  = "103"
	ProfileProgressionComponent = "104"
	PlatformSilverComponent     = "105"

	// Characters (200-299)
	CharactersComponent            = "200"
	CharacterInventoriesComponent  = "201"
	CharacterProgressionsComponent = "202"
	CharacterRenderDataComponent   = "203"
	CharacterActivitiesComponent   = "204"
	CharacterEquipmentComponent    = "205"
	CharacterLoadoutsComponent     = "206"

	// Items (300-399)
	ItemInstancesComponent            = "300"
	ItemObjectivesComponentNumber     = "301"
	ItemPerksComponent                = "302"
	ItemRenderDataComponent           = "303"
	ItemStatsComponentNumber          = "304"
	ItemSocketsComponentNumber        = "305"
	ItemTalentGridsComponent          = "306"
	ItemCommonDataComponent           = "307"
	ItemPlugStatesComponent           = "308"
	ItemPlugObjectivesComponentNumber = "309"
	ItemReusablePlugsComponentNumber  = "310"

	// Vendors & Social (400-599)
	VendorsComponent          = "400"
	VendorCategoriesComponent = "401"
	VendorSalesComponent      = "402"
	KiosksComponent           = "500"
	CurrencyLookupsComponent  = "600"

	// Collections & Triumphs (800-1100)
	PresentationNodesComponent = "700"
	CollectiblesComponent      = "800"
	RecordsComponent           = "900"
	TransitoryComponent        = "1000"
	MetricsComponent           = "1100"
	StringVariablesComponent   = "1200"
	CraftablesComponent        = "1300"
)

// GetCoreComponents returns the essential components to see
// characters, their equipped weapons, perks and stats.
func GetCoreComponents() []string {
	return []string{
		ProfilesComponent,                 // 100
		ProfileInventoriesComponent,       // 102
		CharactersComponent,               // 200
		CharacterInventoriesComponent,     // 201
		CharacterEquipmentComponent,       // 205
		ItemInstancesComponent,            // 300
		ItemObjectivesComponentNumber,     // 301
		ItemStatsComponentNumber,          // 304
		ItemSocketsComponentNumber,        // 305
		ItemPlugObjectivesComponentNumber, // 309
		ItemReusablePlugsComponentNumber,  // 310
	}
}

const (
	TitanClassType   = 0
	HunterClassType  = 1
	WarlockClassType = 2
	UnknownClassType = 3
)

var classNames = map[uint32]string{
	0: "Titan",
	1: "Hunter",
	2: "Warlock",
}

// GetClassName returns the guardian name based on its type (Titan, Hunter, Warlock)
func GetClassName(classType int) string {
	return classNames[uint32(classType)]
}

// ITEMTYPE from Manifest
const (
	None              = 0
	Currency          = 1
	Armor             = 2
	Weapon            = 3
	Message           = 7
	Engram            = 8
	Consumable        = 9
	ExchangeMaterial  = 10
	MissionReward     = 11
	QuestStep         = 12
	QuestStepComplete = 13
	Emblem            = 14
	Quest             = 15
	Subclass          = 16
	ClanBanner        = 17
	Aura              = 18
	Mod               = 19
)

// Hash constants para tipos de stats comunes en armas
const (
	// Main firearm statistics
	StatHashRPM         = uint32(4232813984) // Fire rate
	StatHashImpact      = uint32(4043523819) // Impact
	StatHashRange       = uint32(1240592695) // Range
	StatHashStability   = uint32(155624089)  // Stability
	StatHashHandling    = uint32(943549884)  // Handling
	StatHashReloadSpeed = uint32(4284893193) // Reload speed
	StatHashMagazine    = uint32(3871231066) // Magazine capacity

	// Technical statistics
	StatHashAimAssistance   = uint32(1345609583) // Aim assistance
	StatHashAirborne        = uint32(2714457168) // Airborne accuracy
	StatHashZoom            = uint32(3555269338) // Zoom
	StatHashRecoilDirection = uint32(2715839340) // Recoil direction

	// Sword statistics
	StatHashSwingSpeed      = uint32(2837207746) // Swing speed
	StatHashAmmoCapacity    = uint32(925767036)  // Ammo capacity
	StatHashGuardResistance = uint32(419712076)  // Guard resistance
	StatHashChargeRate      = uint32(3022301683) // Charge rate
	StatHashGuardEndurance  = uint32(3736848092) // Resistencia de guardia
)

var StatsDictionary = map[uint32]string{
	// --- Armas de Fuego ---
	4232813984: "RPM",          // Velocidad de fuego
	4043523819: "Impact",       // Impacto
	1240592695: "Range",        // Alcance
	155624089:  "Stability",    // Estabilidad
	943549884:  "Handling",     // Manejo
	4284893193: "Reload Speed", // Velocidad de recarga
	3871231066: "Magazine",     // Capacidad de cargador

	// --- Stats Técnicos / Extra ---
	1345609583: "Aim Assistance",   // Asistencia de puntería
	2714457168: "Airborne",         // Precisión en aire
	3555269338: "Zoom",             // Zoom
	2715839340: "Recoil Direction", // Dirección del retroceso

	// --- Espadas ---
	2837207746: "Swing Speed",      // Velocidad de oscilación
	925767036:  "Ammo Capacity",    // Capacidad de munición
	419712076:  "Guard Resistance", // Resistencia de guardia
	3022301683: "Charge Rate",      // Velocidad de carga
	3736848092: "Guard Endurance",  // Resistencia de guardia
}

// Hashes de objetivos especiales para rastreo de progreso
const (
	ObjectiveHashWeaponLevel     = uint32(2021564063) // Nivel del arma (ej: 598)
	ObjectiveHashEnemiesDefeated = uint32(38912240)   // Total de bajas/enemigos derrotados
	ObjectiveHashLevelProgress   = uint32(2146410756) // Progreso del nivel actual (%)
)

const BungieCDN = "https://www.bungie.net"

// Traduce el DamageTypeHash a texto legible
func GetDamageName(hash uint32) string {
	switch hash {
	case 1847026933, 3:
		return "Solar"
	case 2301139358, 2303181850, 2:
		return "Arc"
	case 3454344763, 3454344768, 4:
		return "Void"
	case 1513472331, 151347233, 6:
		return "Stasis"
	case 3946443463, 3949783978, 7:
		return "Strand"
	case 3373582059, 1, 0:
		return "Kinetic"
	default:
		return "Kinetic" // If it's 0 or unknown, the vast majority are kinetic
	}
}

func GetSlotName(hash uint32) string {
	switch hash {
	case 1498876634:
		return "Kinetic" // Slot 1
	case 2465295065:
		return "Energy" // Slot 2
	case 953998645:
		return "Power" // Slot 3
	default:
		// Si ves muchos "Other", imprime el hash aquí para debuguear
		return "Other"
	}
}

// Determines the ammo type (for icons or filters)
func GetAmmoTypeName(ammoType int32) string {
	switch ammoType {
	case 1:
		return "Primary"
	case 2:
		return "Special"
	case 3:
		return "Heavy"
	default:
		return "None"
	}
}
