package destiny

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nicoki2004/g2-drc/internal/models"
)

func GetMembershipForCurrentUser(client *Client) (*MembershipResponse, error) {
	resp, err := client.DoRequest("GET", models.URL_MEMBERSHIP_USER)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data MembershipResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decodificando membership: %w", err)
	}

	if data.ErrorCode != 1 {
		return nil, fmt.Errorf("error de Bungie: %s", data.Message)
	}

	return &data, nil
}

// Get profile and save to json. If not exists call bungie API
// GetProfile acepta una lista de componentes para traer exactamente lo necesario.
func GetProfile(client *Client, components ...string) (*ProfileResponse, error) {
	cacheFile := "profile_cache.json"

	// 1. Lógica de Caché (Mantenla para desarrollo)
	if data, err := os.ReadFile(cacheFile); err == nil {
		var cachedProfile ProfileResponse
		if err := json.Unmarshal(data, &cachedProfile); err == nil {
			fmt.Println("⚡ Cargando perfil desde caché local...")
			return &cachedProfile, nil
		}
	}

	if len(components) == 0 {
		components = []string{
			CharactersComponent,               // 200
			CharacterEquipmentComponent,       // 205
			ItemInstancesComponent,            // 300
			ItemStatsComponentNumber,          // 304
			ItemSocketsComponentNumber,        // 305
			ItemReusablePlugsComponentNumber,  // 310
			ItemObjectivesComponentNumber,     // 301
			CharacterInventoriesComponent,     // 201
			ItemPlugObjectivesComponentNumber, // 309
		}
	}

	queryString := strings.Join(components, ",")
	destUrl := fmt.Sprintf("%s?components=%s",
		models.GetProfileURL(PlayStation, client.memberShipId),
		queryString,
	)

	fmt.Printf("🌐 Llamando a Bungie: %s\n", destUrl)

	resp, err := client.DoRequest("GET", destUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decodificando perfil: %w", err)
	}

	if data.ErrorCode != 1 {
		return nil, fmt.Errorf("error de Bungie [%d]: %s", data.ErrorCode, data.Message)
	}

	saveCache(cacheFile, data)

	return &data, nil
}

func saveCache(filename string, data any) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Error serializando caché: %v", err)
		return
	}

	err = os.WriteFile(filename, jsonData, 0o644)
	if err != nil {
		log.Printf("Error escribiendo archivo de caché: %v", err)
		return
	}
	fmt.Println("💾 Perfil guardado en caché local (profile_cache.json)")
}
