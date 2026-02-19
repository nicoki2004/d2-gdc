package destiny

import "time"

// MembershipResponse contiene información de suscripción del usuario en Bungie
// Devuelve todas las cuentas de Destiny conectadas a la cuenta de Bungie
type MembershipResponse struct {
	Response struct {
		DestinyMemberships  []DestinyMembership `json:"destinyMemberships"`
		PrimaryMembershipId string              `json:"primaryMembershipId"`
	} `json:"Response"`
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// DestinyMembership representa una plataforma de juego individual (PS5, Xbox, Steam, etc)
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

// ProfileResponse es la respuesta completa del endpoint de Perfil de Bungie
// Contiene personajes, inventario, equipamiento, stats y todos los detalles técnicos
type ProfileResponse struct {
	Response struct {
		// 100 - INFORMACION DEL PERFIL
		Characters struct {
			Data map[string]CharacterData `json:"data"`
		} `json:"characters"`

		// 102 - EL VAULT (DEPÓSITO)
		ProfileInventory struct {
			Data struct {
				Items []Item `json:"items"`
			} `json:"data"`
		} `json:"profileInventory"`

		// 201 - LAS MOCHILAS (INVENTARIOS)
		CharacterInventories struct {
			Data map[string]CharacterEquipmentData `json:"data"`
		} `json:"characterInventories"`
		// 205 - LO EQUIPADO
		CharacterEquipment struct {
			Data map[string]CharacterEquipmentData `json:"data"`
		} `json:"characterEquipment"`

		// 300-309 - DETALLES TÉCNICOS
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

			// COMPONENTE 309: El que necesitas para las 23,186 bajas
			ItemPlugObjectives struct {
				Data map[string]ItemPlugObjectivesComponent `json:"data"`
			} `json:"plugObjectives"`
			// Otros opcionales pero recomendados
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

// Struct auxiliar para el 309
type ItemPlugObjectivesComponent struct {
	// CAMBIO: Bungie lo llama "objectivesPerPlug" internamente por cada socket
	ObjectivesPerPlug map[string][]ObjectiveData `json:"objectivesPerPlug"`
}

// Esta queda limpia, solo con los sockets
type ItemSocketsComponent struct {
	Sockets []SocketEntry `json:"sockets"`
}

// Estas se mantienen igual
type ItemReusablePlugsComponent struct {
	Plugs map[string][]PlugEntry `json:"plugs"`
}

type PlugEntry struct {
	PlugItemHash uint32 `json:"plugItemHash"`
	CanInsert    bool   `json:"canInsert"`
}

type SocketEntry struct {
	PlugHash  uint32 `json:"plugHash"` // El Hash del Perk o Modificador
	IsEnabled bool   `json:"isEnabled"`
	IsVisible bool   `json:"isVisible"`
}

// CharacterData representa un personaje del usuario (Componente 200)
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
	ItemHash       uint32 `json:"itemHash"`       // Hash de definición del item (referencia al manifest)
	ItemInstanceId string `json:"itemInstanceId"` // ID único de esta instancia específica
}

// ItemStatsComponent representa los stats de una instancia (Componente 304)
type ItemStatsComponent struct {
	Stats map[uint32]StatData `json:"stats"`
}

// StatData contiene el valor individual
type StatData struct {
	StatHash uint32 `json:"statHash"`
	Value    int    `json:"value"`
}

// Para el Poder (Componente 300)
type ItemInstanceData struct {
	PrimaryStat struct {
		Value int `json:"value"` // Aquí vive el Poder (460)
	} `json:"primaryStat"`
}

// Para Bajas y Nivel (Componente 301)
type ItemObjectivesComponent struct {
	Objectives []ObjectiveData `json:"objectives"`
}

type ObjectiveData struct {
	ObjectiveHash   uint32 `json:"objectiveHash"`   // DEBE ser camelCase
	Progress        int    `json:"progress"`        // DEBE ser minúscula
	CompletionValue int    `json:"completionValue"` // DEBE ser camelCase
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

// GetCoreComponents devuelve los componentes esenciales para ver
// personajes, sus armas equipadas, perks y stats.
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

// GetClassName devuelve el nombre del guardián a partir de su tipo (Titan, Hunter, Warlock)
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
	// Estadísticas principales de armas de fuego
	StatHashRPM         = uint32(4232813984) // Velocidad de fuego
	StatHashImpact      = uint32(4043523819) // Impacto
	StatHashRange       = uint32(1240592695) // Alcance
	StatHashStability   = uint32(155624089)  // Estabilidad
	StatHashHandling    = uint32(943549884)  // Manejo
	StatHashReloadSpeed = uint32(4284893193) // Velocidad de recarga
	StatHashMagazine    = uint32(3871231066) // Capacidad de cargador

	// Estadísticas técnicas
	StatHashAimAssistance   = uint32(1345609583) // Asistencia de puntería
	StatHashAirborne        = uint32(2714457168) // Precisión en aire
	StatHashZoom            = uint32(3555269338) // Zoom
	StatHashRecoilDirection = uint32(2715839340) // Dirección del retroceso

	// Estadísticas de espadas
	StatHashSwingSpeed      = uint32(2837207746) // Velocidad de oscilación
	StatHashAmmoCapacity    = uint32(925767036)  // Capacidad de munición
	StatHashGuardResistance = uint32(419712076)  // Resistencia de guardia
	StatHashChargeRate      = uint32(3022301683) // Velocidad de carga
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
		return "Kinetic" // Si es 0 o desconocido, la gran mayoría son cinéticas
	}
}

// Traduce el BucketHash (Slot) a texto legible
func GetSlotName(hash uint32) string {
	switch hash {
	case 1498876634:
		return "Kinetic" // Slot superior
	case 2465295065:
		return "Energy" // Slot medio
	case 95395402:
		return "Power" // Slot inferior
	default:
		return "Other" // Para el caso del segundo objeto que vimos con slot 24222...
	}
}

// Determina el tipo de munición (para iconos o filtros)
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
