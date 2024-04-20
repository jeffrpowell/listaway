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
```
4. `docker compose up`
4. [https://localhost:8080/](https://localhost:8080/)