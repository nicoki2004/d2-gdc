package destiny

import (
	"encoding/json"
	"fmt"
	"os"
)

// Definimos solo lo que necesitamos para ahorrar memoria
type ManifestItem struct {
	DisplayProperties struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"displayProperties"`
	ItemTypeDisplayName string `json:"itemTypeDisplayName"`
	ItemType            int    `json:"itemType"`
	Inventory           struct {
		TierTypeName   string `json:"tierTypeName"`   // "Exotic"
		BucketTypeHash uint32 `json:"bucketTypeHash"` // Para el Slot
	} `json:"inventory"`
	EquippingBlock struct {
		AmmoType int32 `json:"ammoType"` // 1: Primary, 2: Special, 3: Heavy
	} `json:"equippingBlock"`
	DefaultDamageTypeHash uint32 `json:"defaultDamageTypeHash"`
}

func LoadManifestMap(filePath string) (map[string]ManifestItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create the map where we'll store the "translation"
	// We use string as key because the manifest JSON has hashes as strings
	manifestMap := make(map[string]ManifestItem)

	decoder := json.NewDecoder(file)

	// The manifest is a giant object { "hash": {data}, "hash2": {data} }
	// Decode() will process it efficiently
	err = decoder.Decode(&manifestMap)
	if err != nil {
		return nil, fmt.Errorf("error decoding manifest: %w", err)
	}

	return manifestMap, nil
}
