# Date Apps API

## Prerequesites

1. Docker
2. Golang-Migrate
3. Swag
4. Golang
5. Mockery

## Development Guide

### Getting Started

1. Run `make dev args="up"` to run local environment with hot reload, everything inside `args` is basically the `docker compose` arguments. You can exit by `ctrl+c`.
2. The db migration will automatically running.
3. If you want to tear down all the containers you can run `make dev args="down"` or add `-v` to the args if you want to remove the volume too, for example `make dev args="down -v"`
4. To rebuild the docker image if you change something, you can run `make dev args="build"`. This will left dangling images, you need to remove the dangling images manually.
5. dont forget to create `.env` file inside `deployments/development` folder and change the environment variables to your own.

### Custom Docker Compose configuration

if you need to add more configuration and it's different with the default one, you can create a `docker-compose.override.yml` inside `deployments/development` folder, at the same level with the `docker-compose.yaml` file.

for example if you need to expose port for database you can do something like this inside the `docker-compose.override.yml` file

```
version: "3"
services:
  database:
    ports:
      - 3306:3306
```

### Create Migration

You need to install `golang-migrate` manually in your device.

You can install it by `go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`.

Run `migrate create -ext sql -dir migration {migration_name}` to create your migration file

### Running DB Migration manually

Please make sure the `date-apps-be` and `database` containers are running.

everything inside `args` are the `golang-migrate` arguments

Run `make dev-migrate args="up"` to create or modify tables

Run `make dev-migrate args="drop -f"` to drop all tables

Run `make dev-migrate args="down [migration_version]"` to delete some tables based on migration version

Run `make dev-migrate args="force [migration_version]"` to fix the dirty migration

## Testing Guide

###

### Running Test

Run `make test`

### Create Mocks

Run `make mock` if using windows `make mock-win`

## Deployment Guide

### DB Migration

DB Migration on Staging or Prod need to do it manually by running `migrate -database "mysql://<DB_USER>:<DB_PASS>@tcp(<DB_HOST>:<DB_PORT>)/<DB_NAME>" -path migration [up|<any migrate argument>]`

### Deployment

The docker setup is inside `deployments/deploy` folder


## API Documentation
You need to install [swag](https://github.com/swaggo/swag) manually in your device.

You can install it by `go install github.com/swaggo/swag/cmd/swag`.

We use [swag](https://github.com/swaggo/swag) to generate necearry Swagger files for API documentation. Everytime we run `make gen-swagger`, the Swagger documentation will be updated.

To configure your API documentation, please refer to [swag's declarative comments format](https://github.com/swaggo/swag#declarative-comments-format) and [examples](https://github.com/swaggo/swag#examples).

To access the documentation, after you run the app please visit [API DOCUMENTATION](http://localhost:8080/v1/docs/index.html).


## Project Structure

- `cmd/webservice/main.go`: The entry point of the application where the server is initialized and started.
- `internal/api/http/router.go`: Configures all API routes and their corresponding handlers.
- `internal/container/container.go`: Manages dependency injection and application-wide component initialization.
- `internal/api/http/handler`: Contains HTTP handlers that process incoming requests and return responses.
- `internal/api/http/middleware`: Houses middleware functions for authentication, logging, error handling etc.
- `internal/api/http/request`: Defines request structs and validation for incoming API requests.
- `internal/api/http/response`: Contains response structs and helper functions for API responses.
- `internal/api/http/routes`: Groups related API routes and their configurations.
- `internal/model`: Contains domain models and entities that represent the core business objects and data structures.
- `internal/repository`: Data access layer that handles database operations and data persistence.
- `internal/service`: Implements core business logic and coordinates between different layers.
- `internal/usecase`: Orchestrates the flow of data and business rules between different services.
- `infrastructure/config`: Manages application configuration from environment variables and config files.
- `infrastructure/database`: Handles database connections, migrations and database-specific configurations.
- `infrastructure/database/migrations`: Contains database migration files for schema changes and data updates.
- `docs`: Stores auto-generated API documentation and other project documentation.
- `deployments`: Contains Docker files, deployment scripts and environment configurations.
- `pkg`: Reusable packages and utilities that can be shared across projects.
- `internal/test`: Contains test files, test utilities, and test configurations.

