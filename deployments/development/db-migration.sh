#!/bin/bash

migrate -database "mysql://$DBMASTERUSER:$DBMASTERPASS@tcp(database:3306)/$DBMASTERNAME?multiStatements=true" -path infrastucture/db/migrations $@