package destiny

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nicoki2004/g2-drc/internal/auth"
	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/models"
)

type Client struct {
	httpClient *http.Client
	cfg        *config.Config
	token      *models.Token
}

func NewClient(cfg *config.Config, token *models.Token) *Client {
	return &Client{
		httpClient: &http.Client{},
		cfg:        cfg,
		token:      token,
	}
}

func (c *Client) DoRequest(method, url string) (*http.Response, error) {
	// Usamos una función interna para no repetir la lógica de creación de headers
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
		// ¡Importante! Cerrar el cuerpo de la respuesta fallida antes de reintentar
		resp.Body.Close()

		log.Println("🔄 Token expirado, intentando refresh...")

		newToken, err := auth.RefreshToken(c.cfg, c.token)
		if err != nil {
			return nil, fmt.Errorf("sesión expirada: %w", err)
		}

		c.token = newToken

		// Crear nueva petición con el nuevo token
		req, err = makeReq()
		if err != nil {
			return nil, err
		}

		log.Println("✅ Token refrescado, reintentando petición original...")
		return c.httpClient.Do(req)
	}

	return resp, nil
}
