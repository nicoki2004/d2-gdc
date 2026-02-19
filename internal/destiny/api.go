package destiny

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nicoki2004/g2-drc/internal/cache"
	"github.com/nicoki2004/g2-drc/internal/logger"
	"github.com/nicoki2004/g2-drc/internal/models"
)

// checkBungieResponse validates ErrorCode in Bungie API responses
func checkBungieResponse(errorCode int, message string) error {
	if errorCode != 1 {
		return fmt.Errorf("error de Bungie [ErrorCode: %d]: %s", errorCode, message)
	}
	return nil
}

// getCacheFile returns the path to the profile cache file from environment variable CACHE_FILE,
// with a default of "profile_cache.json" if not set.
func getCacheFile() string {
	if path := os.Getenv("CACHE_FILE"); path != "" {
		return path
	}
	return "profile_cache.json"
}

func GetMembershipForCurrentUser(client *Client) (*MembershipResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}
	resp, err := client.DoRequest("GET", models.URL_MEMBERSHIP_USER)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data MembershipResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decodificando membership: %w", err)
	}

	if err := checkBungieResponse(data.ErrorCode, data.Message); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetProfile acepta una lista de componentes para traer exactamente lo necesario.
// Intenta cargar desde caché primero, y si no existe, consulta la API de Bungie.
func GetProfile(client *Client, components ...string) (*ProfileResponse, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}
	cacheManager := cache.NewFileCache()
	cacheFile := getCacheFile()

	// 1. Intentar cargar desde caché
	var cachedProfile ProfileResponse
	if err := cacheManager.Load(cacheFile, &cachedProfile); err == nil {
		logger.GetLogger().Info("Cargando perfil desde caché local")
		return &cachedProfile, nil
	}

	// 2. Si no hay componentes especificados, usar los por defecto
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

	logger.GetLogger().Debug("Llamando a Bungie: %s", destUrl)

	resp, err := client.DoRequest("GET", destUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decodificando perfil: %w", err)
	}

	if err := checkBungieResponse(data.ErrorCode, data.Message); err != nil {
		return nil, err
	}

	// 3. Guardar en caché
	if err := cacheManager.Save(cacheFile, data); err != nil {
		logger.GetLogger().Warn("No se pudo guardar caché: %v", err)
	}

	return &data, nil
}
