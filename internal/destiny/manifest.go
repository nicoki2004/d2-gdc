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
	} `json:"displayProperties"`
	ItemTypeDisplayName string `json:"itemTypeDisplayName"`
	ItemType            int    `json:"itemType"`
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
