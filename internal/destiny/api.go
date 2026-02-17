package destiny

import (
	"encoding/json"
	"fmt"

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

	// Siempre verifica el ErrorCode de Bungie
	if data.ErrorCode != 1 {
		return nil, fmt.Errorf("error de Bungie: %s", data.Message)
	}

	return &data, nil
}
