package destiny

// WeaponMetadata contiene información sobre un arma
type WeaponMetadata struct {
	Level    int
	Kills    int
	Progress float64
	Power    int
}

// WeaponExtractor extrae información de armas desde el perfil
type WeaponExtractor struct{}

// NewWeaponExtractor crea un nuevo extractor
func NewWeaponExtractor() *WeaponExtractor {
	return &WeaponExtractor{}
}

// ExtractMetadata extrae toda la información de una arma
func (we *WeaponExtractor) ExtractMetadata(item Item, profile *ProfileResponse) WeaponMetadata {
	metadata := WeaponMetadata{}

	// Obtener level, kills y progress
	level, kills, progress := GetWeaponMetadata(item.ItemInstanceId, profile)
	metadata.Level = level
	metadata.Kills = kills
	metadata.Progress = progress

	// Obtener power desde itemComponents.instances
	if inst, ok := profile.Response.ItemComponents.Instances.Data[item.ItemInstanceId]; ok {
		metadata.Power = inst.PrimaryStat.Value
	}

	return metadata
}

// ExtractStats extrae los stats estadísticos del arma
func (we *WeaponExtractor) ExtractStats(item Item, profile *ProfileResponse) map[string]int {
	stats := make(map[string]int)

	if statsData, ok := profile.Response.ItemComponents.Stats.Data[item.ItemInstanceId]; ok {
		for hash, stat := range statsData.Stats {
			if statName, exists := StatsDictionary[hash]; exists {
				stats[statName] = stat.Value
			}
		}
	}

	return stats
}

// ExtractSockets extrae información de sockets del arma
func (we *WeaponExtractor) ExtractSockets(item Item, profile *ProfileResponse) SocketsData {
	data := SocketsData{}

	socketsData, hasSockets := profile.Response.ItemComponents.Sockets.Data[item.ItemInstanceId]
	reusableData, hasReusable := profile.Response.ItemComponents.ReusablePlugs.Data[item.ItemInstanceId]

	data.HasSockets = hasSockets
	data.HasReusable = hasReusable
	data.Sockets = socketsData
	data.Plugs = reusableData

	return data
}

// SocketsData contiene información de sockets
type SocketsData struct {
	HasSockets  bool
	HasReusable bool
	Sockets     ItemSocketsComponent
	Plugs       ItemReusablePlugsComponent
}
