# Listaway Application Architecture

## Folder Structure

```
listaway/
├── .devcontainer/        # Development container configuration
│   ├── Dockerfile        # Container definition
│   ├── .env              # Environment variables
│   ├── amd64/            # AMD64 architecture specific files
│   └── arm64/            # ARM64 architecture specific files
├── .github/              # GitHub configuration
│   ├── workflows/        # GitHub Actions workflows
│   │   └── release.yml   # Release automation
│   └── dependabot.yml    # Dependency update configuration
├── .vscode/              # VS Code configuration
│   └── launch.json       # Debug configuration
├── .windsurf/            # Windsurf AI assistant configuration
│   ├── features/         # Feature definitions
│   ├── rules/            # Custom rules
│   └── workflows/        # Workflow templates
├── cmd/
│   └── listaway/         # Application entry point
│       └── main.go       # Main application file
├── internal/
│   ├── constants/        # Application constants and types
│   │   ├── constants.go
│   │   ├── random/
│   │   └── types.go
│   ├── database/         # Database operations
│   │   ├── init.go       # Database initialization
│   │   ├── init.sql      # SQL schema
│   │   └── oidc.go       # OIDC-specific database operations
│   └── handlers/         # HTTP request handlers
│       ├── helper/       # Helper functions for handlers
│       └── middleware/   # HTTP middleware
│       └── oidc/         # OIDC provider integration
└── web/                  # Frontend code
    ├── app/              # Frontend application code
    ├── dist/             # Compiled frontend assets
    ├── node_modules/     # Frontend dependencies
    ├── package.json      # Frontend package configuration
    ├── postcss.config.js # PostCSS configuration
    ├── tailwind.config.js # Tailwind CSS configuration
    ├── web.go            # Web server and frontend integration
    ├── webpack.development.config.js
    └── webpack.production.config.js # Inherits from webpack.development.config.js
```

## Major Areas

### Backend (Go)

1. **Application Entry Point** (`cmd/listaway/main.go`)
   - Initializes and starts the HTTP server
   - Sets up necessary imports for database and handler initialization

2. **Constants** (`internal/constants/`)
   - Defines application-wide constants, types, and router configuration
   - Contains utility functions and data structures used throughout the application
   - Includes OIDC configuration constants (provider URL, client ID, scopes, etc.)

3. **Database** (`internal/database/`)
   - Handles database initialization and schema setup
   - Provides CRUD operations for the core entities:
     - Users: User management and password authentication
     - OIDC: User management and OIDC authentication
     - Lists: List creation and management
     - Items: List items management
     - Collections: Grouping and sharing related sets of lists
     - Reset Tokens: Password reset tokens
     - Group Settings: Settings that apply to an entire group

4. **Handlers** (`internal/handlers/`)
   - Implements HTTP request handlers for all application endpoints
   - Contains middleware for authentication, authorization, and CORS
   - Organizes routes by functional area (admin, authentication, items, lists, collections, sharing)

5. **OIDC Client** (`internal/oidc/`)
   - Manages OIDC provider integration and OAuth2 flow
   - Provider initialization and configuration
   - Authorization URL generation with CSRF protection
   - Token exchange and validation

### Frontend (Web)

1. **Web Integration** (`web/web.go`)
   - Connects the Go backend with the frontend assets
   - Handles serving of static files and web application routes

2. **Frontend Application** (`web/app/`)
   - All user interface code for the application
   - Built with nested HTML snippets (powered by the Go html/template package via `web/web.go`), TailwindCSS, and vanilla Javascript.

3. **Build Configuration**
   - Uses Webpack for bundling assets (development and production configurations)
   - Implements Tailwind CSS for styling
   - Uses PostCSS for CSS processing

### Development Infrastructure

1. **Development Container** (`.devcontainer/`)
   - Provides a consistent development environment using containers
   - Includes configuration for both AMD64 and ARM64 architectures
   - Contains Dockerfile and environment settings

2. **GitHub Integration** (`.github/`)
   - Workflows for CI/CD automation via GitHub Actions
   - Release automation for application deployment
   - Dependabot configuration for automated dependency updates

3. **IDE Configuration** (`.vscode/`)
   - Visual Studio Code settings for consistent development experience
   - Debug configurations for Go application

4. **AI Assistant Configuration** (`.windsurf/`)
   - Configuration for Windsurf AI coding assistant
   - Custom workflows, rules, and feature definitions
   - Helps with maintaining project standards and automating common tasks

## Summary

Listaway is a web application for managing lists, with the following key features:

1. **User Management**: Registration, authentication, and user profiles
2. **Authentication**: Dual support for traditional email/password and OIDC/OAuth2
3. **List Management**: Creating, editing, and organizing lists
4. **Item Management**: Adding, updating, and removing items within lists
5. **Collections**: Grouping related lists together and managing them as a unit (including lists shared from the user group)
6. **Sharing Functionality**: Ability to share lists and collections to the internet or with other users
7. **Admin Features**: Administrative capabilities for user and instance management

The application follows a clean separation between backend (Go) and frontend (JavaScript) with a well-structured organization of code by functional domains. The backend is organized using a layered architecture with database operations, business logic, and HTTP handlers clearly separated.