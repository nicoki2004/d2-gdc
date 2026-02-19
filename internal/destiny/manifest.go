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

	// Creamos el mapa donde guardaremos la "traducción"
	// Usamos string como llave porque el JSON del manifest tiene los hashes como strings
	manifestMap := make(map[string]ManifestItem)

	decoder := json.NewDecoder(file)

	// El manifest es un objeto gigante { "hash": {datos}, "hash2": {datos} }
	// Decode() lo procesará de forma eficiente
	err = decoder.Decode(&manifestMap)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar manifest: %w", err)
	}

	return manifestMap, nil
}
