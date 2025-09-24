# OIDC/OAuth2 Authentication Setup Guide

This guide explains how to configure and use OIDC (OpenID Connect) authentication in Listaway.

## Overview

Listaway now supports OIDC/OAuth2 authentication alongside traditional email/password authentication. Users can:

- Sign in with OIDC providers (Google, GitHub, Microsoft, Auth0, etc.)
- Link OIDC accounts to existing email/password accounts
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

## Database Schema Changes

The implementation automatically adds the following columns to the `listaway.user` table:

- `oidc_provider` (VARCHAR): The OIDC provider name (e.g., "google", "github")
- `oidc_subject` (VARCHAR): The unique subject identifier from the OIDC provider
- `oidc_email` (VARCHAR): The email address from the OIDC provider (may differ from primary email)

## API Endpoints

### OIDC Authentication Endpoints

- `GET /auth/oidc/login` - Initiates OIDC authentication flow
- `GET /auth/oidc/callback` - Handles OIDC callback and completes authentication
- `POST /auth/oidc/link` - Links OIDC account to existing user (requires authentication)
- `POST /auth/oidc/unlink` - Removes OIDC authentication from user account (requires authentication)
- `GET /api/oidc/status` - Returns OIDC configuration status for frontend

### Traditional Authentication Endpoints (Unchanged)

- `GET /auth` - Login page
- `POST /auth` - Email/password login
- `DELETE /auth` - Logout
- `POST /reset` - Request password reset
- `GET /reset/{token}` - Password reset form
- `POST /reset/{token}` - Complete password reset

## User Experience

### New User Flow

1. User visits `/auth` (login page)
2. If OIDC is enabled, they see both email/password form and OIDC button
3. Clicking OIDC button redirects to provider
4. After successful authentication, user is redirected to `/list`
5. New user account is automatically created

### Existing User Flow

1. Existing users can continue using email/password authentication
2. They can optionally link their OIDC account for future convenience
3. Once linked, they can use either authentication method

### Account Linking

- Users with existing email/password accounts can link OIDC accounts
- If OIDC email matches existing account email, accounts are automatically linked
- Users can unlink OIDC accounts if needed

## Security Features

- **State Parameter**: CSRF protection using cryptographically secure random state
- **Nonce Validation**: ID token replay protection
- **Token Verification**: Full OIDC ID token signature and claims validation
- **Session Management**: Consistent session handling with existing authentication
- **Account Linking**: Secure linking based on email verification

## Troubleshooting

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

### Logs

Enable debug logging to troubleshoot issues:

```bash
# Check application logs for OIDC-related messages
docker logs listaway-container

# Look for messages like:
# "OIDC authentication routes registered"
# "OIDC authentication successful for user ID X"
# "Failed to initialize OIDC client: ..."
```

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

## Production Considerations

1. **HTTPS Required**: Most OIDC providers require HTTPS in production
2. **Secure Secrets**: Store client secrets securely (environment variables, secrets management)
3. **Session Security**: Use secure session configuration (HTTPS, secure cookies)
4. **Rate Limiting**: Consider rate limiting on authentication endpoints
5. **Monitoring**: Monitor authentication success/failure rates
6. **Backup Authentication**: Always maintain email/password as backup authentication method

## Migration Guide

### Existing Installations

1. Update to the latest version with OIDC support
2. Run database migrations (automatic on startup)
3. Configure OIDC environment variables
4. Test authentication flows
5. Communicate changes to users

### Rollback Plan

If issues occur:

1. Set `OIDC_ENABLED=false`
2. Restart application
3. OIDC features will be disabled, traditional auth continues working
4. OIDC data in database remains intact for future re-enablement
