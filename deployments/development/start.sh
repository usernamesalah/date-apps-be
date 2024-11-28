#!/bin/bash

go mod tidy
go mod vendor

./docker/development/db-migration.sh up
exec air -c deployments/development/.air.multitenant.toml