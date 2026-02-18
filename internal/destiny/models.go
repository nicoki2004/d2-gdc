package destiny

type MembershipResponse struct {
	Response struct {
		DestinyMemberships  []DestinyMembership `json:"destinyMemberships"`
		PrimaryMembershipId string              `json:"primaryMembershipId"`
	} `json:"Response"`
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

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

type CharacterData struct {
	CharacterId  string `json:"characterId"`
	ClassType    int    `json:"classType"`
	Light        int    `json:"light"`
	MinutesTotal string `json:"minutesPlayedTotal"`
	EmblemPath   string `json:"emblemBackgroundPath"`
}

type CharacterEquipmentData struct {
	Items []Item `json:"items"`
}

type Item struct {
	ItemHash       uint32 `json:"itemHash"`       // ¡El famoso Hash!
	ItemInstanceId string `json:"itemInstanceId"` // ID único de tu arma específica
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

func getClassName(classType int) string {
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

var StatsDictionary = map[uint32]string{
	// --- Armas de Fuego (Basado en tus capturas de Commemoration y Drang) ---

	4232813984: "RPM",          // Valor 450 en Commemoration
	4043523819: "Impact",       // Valor 41 en Commemoration
	1240592695: "Range",        // Valor 70 en Commemoration
	155624089:  "Stability",    // Valor 57 en Commemoration
	943549884:  "Handling",     // Valor 77 en Commemoration
	4284893193: "Reload Speed", // Hash estándar de recarga
	3871231066: "Magazine",     // Valor 75 en Commemoration

	// --- Stats Técnicos / Extra ---
	1345609583: "Aim Assistance",   // Valor 64 en Commemoration
	2714457168: "Airborne",         // Valor 18 en Commemoration
	3555269338: "Zoom",             // Valor 16 en Commemoration
	2715839340: "Recoil Direction", // Valor 100 en Commemoration

	// --- Espadas (Tu Synanceia) ---
	2837207746: "Swing Speed",      // Valor 40 en Synanceia
	925767036:  "Ammo Capacity",    // Valor 71 en Synanceia
	419712076:  "Guard Resistance", // Valor 41 en Synanceia
	3022301683: "Charge Rate",      // Valor 50 en Synanceia
	3736848092: "Guard Endurance",  // Valor 43 en Synanceia
}

const (
	StatWeaponLevel     = 2021564063 // Nivel de arma (598)
	StatEnemiesDefeated = 38912240   // Bajas (23,186)
	StatLevelProgress   = 2146410756 // Progreso del nivel (24%)
)
