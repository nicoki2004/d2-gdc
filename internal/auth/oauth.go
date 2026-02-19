package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/nicoki2004/g2-drc/internal/config"
	"github.com/nicoki2004/g2-drc/internal/decoding"
	"github.com/nicoki2004/g2-drc/internal/models"
)

// AuthURL builds the Bungie OAuth authorization URL.
func AuthURL(cfg config.Config) (string, error) {
	v := url.Values{}
	v.Set("client_id", cfg.ClientID)
	v.Set("response_type", "code")
	// v.Set("state", state)
	v.Set("redirect_uri", cfg.RedirectURL)

	return models.AUTH_URL_PREFIX + v.Encode(), nil
}

// Change a code for a Token
func ExchangeCode(cfg *config.Config, code string) (*models.Token, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	if code == "" {
		return nil, fmt.Errorf("authorization code cannot be empty")
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.Secret)
	if cfg.RedirectURL != "" {
		form.Set("redirect_uri", cfg.RedirectURL)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", models.AUTH_TOKEN_URL_PREFIX, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cfg.ApiKey != "" {
		req.Header.Set("X-API-Key", cfg.ApiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errorDetail any
		if err := json.NewDecoder(resp.Body).Decode(&errorDetail); err != nil {
			errorDetail = "unable to parse error details"
		}
		return nil, fmt.Errorf("token exchange failed: %s - Detail: %v", resp.Status, errorDetail)
	}

	t, err := decoding.DecodeToken(resp)
	if err != nil {
		return nil, err
	}
	t.ReceivedAt = time.Now()
	return t, nil
}

// getTokenFile returns the path to the token file from environment variable TOKEN_FILE,
// with a default of "token.json" if not set.
func getTokenFile() string {
	if path := os.Getenv("TOKEN_FILE"); path != "" {
		return path
	}
	return "token.json"
}

// Save Token to a JSON File
func SaveTokenJSON(token *models.Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getTokenFile(), data, 0o644)
}

// Reload a Token from a JSON file
func ReadTokenJSON() (*models.Token, error) {
	data, err := os.ReadFile(getTokenFile())
	if err != nil {
		return nil, err
	}
	var token models.Token
	err = json.Unmarshal(data, &token)
	return &token, err
}

// Refres token and return a new token
func RefreshToken(cfg *config.Config, oldToken *models.Token) (*models.Token, error) {
	if cfg == nil {
		return &models.Token{}, fmt.Errorf("config cannot be nil")
	}
	if oldToken == nil {
		return &models.Token{}, fmt.Errorf("oldToken cannot be nil")
	}
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", oldToken.RefreshToken)
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.Secret)

	resp, err := http.PostForm("https://www.bungie.net/platform/app/oauth/token/", data)
	if err != nil {
		return &models.Token{}, err
	}
	newToken, err := decoding.DecodeToken(resp)
	if err != nil {
		return &models.Token{}, err
	}
	SaveTokenJSON(newToken)

	return newToken, nil
}
