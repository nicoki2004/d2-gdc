package destiny

// WeaponMetadata contains information about a weapon
type WeaponMetadata struct {
	Level    int
	Kills    int
	Progress float64
	Power    int
}

// WeaponExtractor extracts weapon information from the profile
type WeaponExtractor struct{}

// NewWeaponExtractor creates a new extractor
func NewWeaponExtractor() *WeaponExtractor {
	return &WeaponExtractor{}
}

// ExtractMetadata extracts all information about a weapon
func (we *WeaponExtractor) ExtractMetadata(item Item, profile *ProfileResponse) WeaponMetadata {
	metadata := WeaponMetadata{}

	// Get level, kills and progress
	level, kills, progress := GetWeaponMetadata(item.ItemInstanceId, profile)
	metadata.Level = level
	metadata.Kills = kills
	metadata.Progress = progress

	// Get power from itemComponents.instances
	if inst, ok := profile.Response.ItemComponents.Instances.Data[item.ItemInstanceId]; ok {
		metadata.Power = inst.PrimaryStat.Value
	}

	return metadata
}

// ExtractStats extracts weapon statistical stats
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

// ExtractSockets extracts weapon socket information
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

// SocketsData contains socket information
type SocketsData struct {
	HasSockets  bool
	HasReusable bool
	Sockets     ItemSocketsComponent
	Plugs       ItemReusablePlugsComponent
}
