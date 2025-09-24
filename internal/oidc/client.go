package oidc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jeffrpowell/listaway/internal/constants"
	"golang.org/x/oauth2"
)

// OIDCClient represents an OIDC client configuration
type OIDCClient struct {
	Provider     *oidc.Provider
	OAuth2Config oauth2.Config
	Verifier     *oidc.IDTokenVerifier
	ProviderName string
}

// OIDCClaims represents the claims we extract from OIDC tokens
type OIDCClaims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

var globalOIDCClient *OIDCClient

// InitOIDCClient initializes the global OIDC client if OIDC is enabled
func InitOIDCClient() error {
	if constants.OIDC_ENABLED != "true" {
		return nil // OIDC is disabled
	}

	if constants.OIDC_PROVIDER_URL == "" || constants.OIDC_CLIENT_ID == "" || constants.OIDC_CLIENT_SECRET == "" {
		return fmt.Errorf("OIDC is enabled but required configuration is missing")
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, constants.OIDC_PROVIDER_URL)
	if err != nil {
		return fmt.Errorf("failed to create OIDC provider: %v", err)
	}

	// Parse scopes
	scopes := strings.Fields(constants.OIDC_SCOPES)
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"}
	}

	// Configure OAuth2
	oauth2Config := oauth2.Config{
		ClientID:     constants.OIDC_CLIENT_ID,
		ClientSecret: constants.OIDC_CLIENT_SECRET,
		RedirectURL:  constants.OIDC_REDIRECT_URL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	// Configure ID token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: constants.OIDC_CLIENT_ID,
	})

	// Determine provider name from URL
	providerName := getProviderName(constants.OIDC_PROVIDER_URL)

	globalOIDCClient = &OIDCClient{
		Provider:     provider,
		OAuth2Config: oauth2Config,
		Verifier:     verifier,
		ProviderName: providerName,
	}

	return nil
}

// GetOIDCClient returns the global OIDC client
func GetOIDCClient() *OIDCClient {
	return globalOIDCClient
}

// IsOIDCEnabled returns true if OIDC is enabled and configured
func IsOIDCEnabled() bool {
	return constants.OIDC_ENABLED == "true" && globalOIDCClient != nil
}

// GenerateAuthURL generates an OAuth2 authorization URL with state parameter
func (c *OIDCClient) GenerateAuthURL(state string) string {
	return c.OAuth2Config.AuthCodeURL(state)
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (c *OIDCClient) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.OAuth2Config.Exchange(ctx, code)
}

// VerifyIDToken verifies and extracts claims from an ID token
func (c *OIDCClient) VerifyIDToken(ctx context.Context, rawIDToken string) (*OIDCClaims, error) {
	idToken, err := c.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ID token: %v", err)
	}

	var claims OIDCClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %v", err)
	}

	return &claims, nil
}

// GenerateState generates a cryptographically secure random state parameter
func GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// getProviderName extracts a friendly provider name from the provider URL
func getProviderName(providerURL string) string {
	switch {
	case strings.Contains(providerURL, "accounts.google.com"):
		return "google"
	case strings.Contains(providerURL, "github.com"):
		return "github"
	case strings.Contains(providerURL, "login.microsoftonline.com"):
		return "microsoft"
	case strings.Contains(providerURL, "auth0.com"):
		return "auth0"
	default:
		return "oidc"
	}
}
