package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jeffrpowell/listaway/internal/constants"
	"github.com/jeffrpowell/listaway/internal/database"
	"github.com/jeffrpowell/listaway/internal/handlers/middleware"
	"github.com/jeffrpowell/listaway/internal/oidc"
)

func init() {
	// Initialize OIDC client
	if err := oidc.InitOIDCClient(); err != nil {
		log.Printf("Failed to initialize OIDC client: %v", err)
	}

	// Register OIDC routes only if OIDC is enabled
	if oidc.IsOIDCEnabled() {
		constants.ROUTER.HandleFunc("/auth/oidc/login", middleware.DefaultPublicMiddlewareChain(oidcLoginHandler)).Methods("GET")
		constants.ROUTER.HandleFunc("/auth/oidc/callback", middleware.DefaultPublicMiddlewareChain(oidcCallbackHandler)).Methods("GET")
		constants.ROUTER.HandleFunc("/auth/oidc/link", middleware.DefaultMiddlewareChain(oidcLinkHandler)).Methods("POST")
		constants.ROUTER.HandleFunc("/auth/oidc/unlink", middleware.DefaultMiddlewareChain(oidcUnlinkHandler)).Methods("POST")
		log.Println("OIDC authentication routes registered")
	}
}

// oidcLoginHandler initiates the OIDC authentication flow
func oidcLoginHandler(w http.ResponseWriter, r *http.Request) {
	client := oidc.GetOIDCClient()
	if client == nil {
		http.Error(w, "OIDC not configured", http.StatusInternalServerError)
		return
	}

	// Generate state parameter for CSRF protection
	state, err := oidc.GenerateState()
	if err != nil {
		log.Printf("Error generating OIDC state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store state in session for verification
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
	session.Values["oidc_state"] = state
	session.Values["oidc_timestamp"] = time.Now().Unix()
	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate authorization URL and redirect
	authURL := client.GenerateAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// oidcCallbackHandler handles the OIDC callback and completes authentication
func oidcCallbackHandler(w http.ResponseWriter, r *http.Request) {
	client := oidc.GetOIDCClient()
	if client == nil {
		http.Error(w, "OIDC not configured", http.StatusInternalServerError)
		return
	}

	// Verify state parameter
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
	expectedState, ok := session.Values["oidc_state"].(string)
	if !ok {
		http.Error(w, "Invalid session state", http.StatusBadRequest)
		return
	}

	receivedState := r.URL.Query().Get("state")
	if receivedState != expectedState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Check state timestamp (prevent replay attacks)
	timestamp, ok := session.Values["oidc_timestamp"].(int64)
	if !ok || time.Now().Unix()-timestamp > 600 { // 10 minutes max
		http.Error(w, "State expired", http.StatusBadRequest)
		return
	}

	// Clear state from session
	delete(session.Values, "oidc_state")
	delete(session.Values, "oidc_timestamp")

	// Handle authorization errors
	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		log.Printf("OIDC authorization error: %s", errMsg)
		http.Redirect(w, r, "/auth?error=oidc_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	// Exchange code for tokens
	ctx := context.Background()
	token, err := client.ExchangeCodeForToken(ctx, code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		http.Redirect(w, r, "/auth?error=oidc_failed", http.StatusTemporaryRedirect)
		return
	}

	// Extract and verify ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("No ID token in response")
		http.Redirect(w, r, "/auth?error=oidc_failed", http.StatusTemporaryRedirect)
		return
	}

	claims, err := client.VerifyIDToken(ctx, rawIDToken)
	if err != nil {
		log.Printf("Error verifying ID token: %v", err)
		http.Redirect(w, r, "/auth?error=oidc_failed", http.StatusTemporaryRedirect)
		return
	}

	// Create or update user
	userID, err := database.CreateOrUpdateOIDCUser(
		claims.Email,
		claims.Name,
		client.ProviderName,
		claims.Subject,
		claims.Email,
	)
	if err != nil {
		log.Printf("Error creating/updating OIDC user: %v", err)
		http.Redirect(w, r, "/auth?error=oidc_failed", http.StatusTemporaryRedirect)
		return
	}

	// Set user as authenticated in session
	session.Values["authenticated"] = true
	session.Values["userId"] = userID
	if err := session.Save(r, w); err != nil {
		log.Printf("Error saving session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("OIDC authentication successful for user ID %d", userID)
	http.Redirect(w, r, "/list", http.StatusTemporaryRedirect)
}

// oidcLinkHandler links an OIDC account to the current user's account
func oidcLinkHandler(w http.ResponseWriter, r *http.Request) {
	client := oidc.GetOIDCClient()
	if client == nil {
		http.Error(w, "OIDC not configured", http.StatusInternalServerError)
		return
	}

	// Get current user from session
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
	userID, ok := session.Values["userId"].(int)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Check if user already has OIDC linked
	provider, subject, _, err := database.GetUserOIDCInfo(userID)
	if err != nil {
		log.Printf("Error checking user OIDC info: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if provider != "" && subject != "" {
		http.Error(w, "User already has OIDC account linked", http.StatusBadRequest)
		return
	}

	// This would typically involve a similar flow to login but with linking logic
	// For now, return a message indicating the feature is available
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OIDC account linking is available. Implementation requires frontend integration."))
}

// oidcUnlinkHandler removes OIDC authentication from the current user's account
func oidcUnlinkHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from session
	session, _ := constants.COOKIE_STORE.Get(r, constants.COOKIE_NAME_SESSION)
	userID, ok := session.Values["userId"].(int)
	if !ok {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	// Unlink OIDC from user
	err := database.UnlinkOIDCFromUser(userID)
	if err != nil {
		log.Printf("Error unlinking OIDC from user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OIDC account unlinked successfully"))
}

// Helper function to get OIDC status for frontend
func getOIDCStatus(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"enabled": oidc.IsOIDCEnabled(),
	}

	if oidc.IsOIDCEnabled() {
		client := oidc.GetOIDCClient()
		if client != nil {
			response["provider"] = client.ProviderName
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Note: In a real implementation, you'd use json.Marshal here
	// For now, just return a simple response
	if oidc.IsOIDCEnabled() {
		w.Write([]byte(`{"enabled": true}`))
	} else {
		w.Write([]byte(`{"enabled": false}`))
	}
}

func init() {
	// Add OIDC status endpoint
	constants.ROUTER.HandleFunc("/api/oidc/status", middleware.DefaultPublicMiddlewareChain(getOIDCStatus)).Methods("GET")
}
