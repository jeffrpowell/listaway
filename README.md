<img src="./web/app/images/ListawayWordmarkDark.png" alt="Listaway logo" title="Listaway" height="300">

This self-hostable application allows authenticated users to publish one or more lists of items publicly. These lists can either be for tracking purposes (e.g. a list of books to read, a list of components in a custom computer build, a list of favorite local places, etc.) or for wishlist purposes (e.g. a gift wishlist, a list of tasks you need help with, etc.). The items can be freeform text or a URL to details about the item. Shared lists incorporate a random string in the URL to give a little protection against guessing (thus allowing you to share the link and access it without requiring authentication).

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
version: "3.9"
services:
  listaway:
    image: "ghcr.io/jeffrpowell/listaway:v1.11.0"
    ports:
      - "8080:8080"
    env_file:
      - .env
```
3. Make a `.env` file
```
# Required configuration
LISTAWAY_AUTH_KEY=[random alphanumeric 128-character string]
PORT=8080
POSTGRES_USER=listaway
POSTGRES_PASSWORD=password
POSTGRES_HOST=[pghost]
POSTGRES_DATABASE=listaway

# Optional SMTP configuration for password reset emails (defaults will cause email bodies to be logged instead of sent outbound)
# SMTP_HOST=smtp.example.com        # SMTP server hostname (default "")
# SMTP_PORT=587                     # SMTP server port (typically 25, 465, or 587, default 587)
# SMTP_USER=username               # SMTP authentication username (default "")
# SMTP_PASSWORD=password           # SMTP authentication password (default "")
# SMTP_FROM=noreply@example.com    # Sender email address (default "noreply@listaway.dev)
# SMTP_SECURE=true                 # Use TLS/SSL (true/false, default true)
# APP_URL=https://listaway.your-domain.com     # Base URL of your application (for reset links, default "http://localhost:8080")
```
4. `docker compose up`
5. [https://localhost:8080/](https://localhost:8080/) (All paths will 303 to [https://localhost:8080/admin/register](https://localhost:8080/admin/register))


## Password Reset Feature

Listaway includes a password reset feature that allows users to reset their password via email. When a user requests a password reset:

1. The system generates a secure token and stores it in the database.
2. An email with a password reset link is sent to the user's registered email.
3. The link is valid for 1 hour and can only be used once.

To enable email delivery for password resets, configure the SMTP settings in your `.env` file as shown above. If SMTP is not configured, the application will log the reset emails to the console instead of sending them.

## Build from source
This repository is provided with a configured devcontainer that is available to assist in quickly bootstrapping a local development environment suitable to build and run this application locally. 

1. Install Docker, VS Code, and the Dev Containers extension in VS Code
    * Make sure your Docker engine is running
2. Click the `><` in the bottom-left corner -> `Reopen in Container`

It is possible to do local development without using the devcontainer. Generally this involves following the steps encoded in the files under the `.devcontainer` folder.

Once you've installed the necessary tooling for this project (notably Go, Postgres, NodeJs, and Yarn), these steps will get the application server running:

1. `yarn --cwd web run build`
    * For the devcontainer users, the devcontainer will run this automatically for you each time you open it (you may need to manually "Refresh Explorer" to see the generated files for the first time)
2. `go run cmd/listaway/main.go`
    * Or leverage the `launch.json` file in VS Code
3. The URL of the server is printed to stdout