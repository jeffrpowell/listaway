# OIDC/OAuth2 Authentication Setup Guide

This guide explains how to configure and use OIDC (OpenID Connect) authentication in Listaway.

## Overview

Listaway supports OIDC/OAuth2 authentication alongside traditional email/password authentication. Users can:

- Sign in with OIDC providers (Google, GitHub, Microsoft, Auth0, etc.)
- Link OIDC accounts to existing email/password accounts (NOT vice versa)
- Use both authentication methods interchangeably

## Configuration

### Environment Variables

Add the following environment variables to enable OIDC authentication:

```bash
# Enable OIDC authentication
OIDC_ENABLED=true

# OIDC Provider Configuration
OIDC_PROVIDER_URL=https://accounts.google.com  # Example for Google
OIDC_CLIENT_ID=your-client-id
OIDC_CLIENT_SECRET=your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/oidc/callback
OIDC_SCOPES="openid profile email"
```

### Provider-Specific Examples

#### Google OAuth2

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable Google+ API
4. Create OAuth2 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/oidc/callback`

```bash
OIDC_ENABLED=true
OIDC_PROVIDER_URL=https://accounts.google.com
OIDC_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com
OIDC_CLIENT_SECRET=GOCSPX-your-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/oidc/callback
OIDC_SCOPES="openid profile email"
```

#### GitHub OAuth2

1. Go to GitHub Settings > Developer settings > OAuth Apps
2. Create a new OAuth App
3. Set Authorization callback URL: `http://localhost:8080/auth/oidc/callback`

```bash
OIDC_ENABLED=true
OIDC_PROVIDER_URL=https://github.com
OIDC_CLIENT_ID=your-github-client-id
OIDC_CLIENT_SECRET=your-github-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/oidc/callback
OIDC_SCOPES="openid profile email"
```

#### Microsoft Azure AD

1. Go to Azure Portal > Azure Active Directory > App registrations
2. Create a new registration
3. Add redirect URI: `http://localhost:8080/auth/oidc/callback`

```bash
OIDC_ENABLED=true
OIDC_PROVIDER_URL=https://login.microsoftonline.com/{tenant-id}/v2.0
OIDC_CLIENT_ID=your-azure-client-id
OIDC_CLIENT_SECRET=your-azure-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/oidc/callback
OIDC_SCOPES="openid profile email"
```

#### Auth0

1. Go to Auth0 Dashboard > Applications
2. Create a new Single Page Application
3. Add callback URL: `http://localhost:8080/auth/oidc/callback`

```bash
OIDC_ENABLED=true
OIDC_PROVIDER_URL=https://your-domain.auth0.com
OIDC_CLIENT_ID=your-auth0-client-id
OIDC_CLIENT_SECRET=your-auth0-client-secret
OIDC_REDIRECT_URL=http://localhost:8080/auth/oidc/callback
OIDC_SCOPES="openid profile email"
```

## Docker Compose Example

```yaml
version: '3.8'
services:
  listaway:
    image: listaway:latest
    ports:
      - "8080:8080"
    environment:
      # Database configuration
      POSTGRES_HOST: postgres
      POSTGRES_USER: listaway
      POSTGRES_PASSWORD: listaway
      POSTGRES_DB: listaway
      
      # OIDC configuration
      OIDC_ENABLED: "true"
      OIDC_PROVIDER_URL: "https://accounts.google.com"
      OIDC_CLIENT_ID: "your-client-id"
      OIDC_CLIENT_SECRET: "your-client-secret"
      OIDC_REDIRECT_URL: "http://localhost:8080/auth/oidc/callback"
      OIDC_SCOPES: "openid profile email"
      
      # Application configuration
      APP_URL: "http://localhost:8080"
      LISTAWAY_AUTH_KEY: "your-32-character-secret-key-here"
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: listaway
      POSTGRES_PASSWORD: listaway
      POSTGRES_DB: listaway
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### Common Issues

1. **"OIDC not configured" error**
   - Ensure `OIDC_ENABLED=true`
   - Verify all required environment variables are set
   - Check provider URL is accessible

2. **"Invalid state parameter" error**
   - Usually indicates CSRF attack or session issues
   - Clear browser cookies and try again
   - Check session store configuration

3. **"Authentication failed" error**
   - Verify client ID and secret are correct
   - Check redirect URL matches exactly
   - Ensure provider is properly configured

4. **Database errors**
   - Run database migrations to add OIDC columns
   - Check database connectivity
   - Verify user table permissions

## Development

### Adding New Providers

To add support for additional OIDC providers:

1. Update `getProviderName()` function in `internal/oidc/client.go`
2. Add provider-specific configuration examples to this documentation
3. Test the provider configuration thoroughly

### Testing

1. Set up a test OIDC provider (e.g., Google OAuth2 playground)
2. Configure environment variables
3. Test authentication flow end-to-end
4. Verify account linking functionality
5. Test error scenarios (invalid tokens, network issues, etc.)