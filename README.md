# Listaway

This self-hostable application allows authenticated users to publish one or more lists of items publicly. These lists can either be for tracking purposes (e.g. a list of books to read, a list of components in a custom computer build, a list of favorite local places, etc.) or for wishlist purposes (e.g. a gift wishlist, a list of tasks you need help with, etc.). The items can be freeform text or a URL to details about the item. The public link (which you control) to these lists do not require authentication.

## Quick start

1. Configure your PostgresDB instance
```sql
--connect to your postgres server with an admin role
CREATE ROLE listaway LOGIN PASSWORD 'password';
CREATE DATABASE listaway;
GRANT CONNECT ON DATABASE listaway TO listaway;
--connect to your new listaway database with an admin role
CREATE SCHEMA listaway;
GRANT CREATE, USAGE ON SCHEMA listaway to listaway;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA listaway TO listaway;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA listaway TO listaway;

```
2. Make a `docker-compose.yml` file
```yaml
version: '3.9'

services:
  listaway:
    image: "https://github.com/jeffrpowell/listaway:1.0"
    ports:
      - "8080:8080"
    env_file:
      - .env
```
3. Make a `.env` file
```
POSTGRES_USER=listaway
POSTGRES_PASSWORD=password
POSTGRES_HOST=[pghost]
POSTGRES_DATABASE=listaway
LISTAWAY_AUTH_KEY=[random alphanumeric 128-character string]
```
4. `docker compose up`
4. [https://localhost:8080/](https://localhost:8080/)

## Build from source
This repository is provided with a configured devcontainer available to assist in quickly bootstrapping a local developmment environment suitable to build and run this application locally. 

1. Install Docker, VS Code, and the Dev Containers extension in VS Code
    * Make sure your Docker engine is running
2. Click the `><` in the bottom-left corner -> `Reopen in Container`

It is possible to do local development without using the devcontainer. Generally this involves following the steps encoded in the files under the `.devcontainer` folder.

Once you've installed the necessary tooling for this project (notably Go, Postgres, and the TailindCSS Standalone CLI), these steps will get the application server running:

1. `tailwindcss -i web/root.css -o internal/handlers/assets/root.css --watch`
    * For the devcontainer users, the devcontainer will run this automatically for you
2. `go run cmd/listaway/main.go`
    * Or leverage the `launch.json` file in VS Code
3. The URL of the server is printed to stdout