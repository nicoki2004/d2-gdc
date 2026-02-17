package decoding

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nicoki2004/g2-drc/internal/models"
)

func DecodeToken(resp *http.Response) (*models.Token, error) {
	defer resp.Body.Close()

	var token models.Token

	err := json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar el token de Bungie: %w", err)
	}

	return &token, nil
}
