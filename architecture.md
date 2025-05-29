# Listaway Application Architecture

## Folder Structure

```
listaway/
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
│   └── handlers/         # HTTP request handlers
│       ├── helper/       # Helper functions for handlers
│       └── middleware/   # HTTP middleware
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

3. **Database** (`internal/database/`)
   - Handles database initialization and schema setup
   - Provides CRUD operations for the core entities:
     - Users: User management and authentication
     - Lists: List creation and management
     - Items: List items management
     - Collections: Grouping and organizing related lists

4. **Handlers** (`internal/handlers/`)
   - Implements HTTP request handlers for all application endpoints
   - Contains middleware for authentication, authorization, and CORS
   - Organizes routes by functional area (admin, authentication, items, lists, collections, sharing)

### Frontend (Web)

1. **Web Integration** (`web/web.go`)
   - Connects the Go backend with the frontend assets
   - Handles serving of static files and web application routes

2. **Frontend Application** (`web/app/`)
   - Contains the JavaScript application code
   - Likely implements the user interface for the list management functionality

3. **Build Configuration**
   - Uses Webpack for bundling assets (development and production configurations)
   - Implements Tailwind CSS for styling
   - Uses PostCSS for CSS processing

## Summary

Listaway is a web application for managing lists, with the following key features:

1. **User Management**: Registration, authentication, and user profiles
2. **List Management**: Creating, editing, and organizing lists
3. **Item Management**: Adding, updating, and removing items within lists
4. **Collections**: Grouping related lists together and managing them as a unit
5. **Sharing Functionality**: Ability to share lists and collections with other users
6. **Admin Features**: Administrative capabilities for user and instance management

The application follows a clean separation between backend (Go) and frontend (JavaScript) with a well-structured organization of code by functional domains. The backend is organized using a layered architecture with database operations, business logic, and HTTP handlers clearly separated.