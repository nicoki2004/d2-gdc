package destiny

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/nicoki2004/g2-drc/db/sqlc"
	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/models"
)

type Client struct {
	httpClient   *http.Client
	cfg          *config.Config
	token        *models.Token
	memberShipId string
	DbQuery      *db.Queries
}

func NewClient(cfg *config.Config, token *models.Token) *Client {
	return &Client{
		httpClient: &http.Client{},
		cfg:        cfg,
		token:      token,
	}
}

func (c *Client) SetDbQuery(query *db.Queries) {
	c.DbQuery = query
}

func (c *Client) SetMembershipId(mId string) {
	c.memberShipId = mId
}

func (c *Client) DoRequest(method, url string) (*http.Response, error) {
	makeReq := func() (*http.Request, error) {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-API-Key", c.cfg.ApiKey)
		req.Header.Add("Authorization", "Bearer "+c.token.AccessToken)
		return req, nil
	}

	req, err := makeReq()
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Si el token expiró (401)
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()

		log.Println("🔄 Token expirado, intentando refresh...")

		newToken, err := auth.RefreshToken(c.cfg, c.token)
		if err != nil {
			return nil, fmt.Errorf("sesión expirada: %w", err)
		}

		c.token = newToken

		req, err = makeReq()
		if err != nil {
			return nil, err
		}

		log.Println("✅ Token refrescado, reintentando petición original...")
		return c.httpClient.Do(req)
	}

	return resp, nil
}
