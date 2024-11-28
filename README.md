# Date Apps API

## Prerequesites

1. Docker
2. Golang-Migrate

## Development Guide

### Getting Started

1. Run `make dev args="up"` to run local environment with hot reload, everything inside `args` is basically the `docker compose` arguments. You can exit by `ctrl+c`.
2. The db migration will automatically running.
3. If you want to tear down all the containers you can run `make dev args="down"` or add `-v` to the args if you want to remove the volume too, for example `make dev args="down -v"`
4. To rebuild the docker image if you change something, you can run `make dev args="build"`. This will left dangling images, you need to remove the dangling images manually.

### Custom Docker Compose configuration

if you need to add more configuration and it's different with the default one, you can create a `docker-compose.override.yml` inside `docker/development` folder, at the same level with the `docker-compose.yaml` file.

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

## Testing Guid

###

### Running Test

Run `make test`

### Create Mocks

Run `make mock`

## Deployment Guide

### DB Migration

DB Migration on Staging or Prod need to do it manually by running `migrate -database "mysql://<DB_USER>:<DB_PASS>@tcp(<DB_HOST>:<DB_PORT>)/<DB_NAME>" -path migration [up|<any migrate argument>]`

### Deployment

The docker setup is inside `deployments/deploy` folder

#### Staging

1. Please make sure you already setup the required environment variable inside the Server at `bashrc` file and `source` it
2. Application will be deployed to staging automatically when the PR got merged.

#### Production

Coming Soon
