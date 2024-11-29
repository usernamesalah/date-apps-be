API_DOCS_PATH = docs/

.PHONY: dev
dev:
	@echo "> Run Date API Service for Development with default config ..."
	@docker compose --project-directory ./deployments/development -p date-apps-be $(args)

.PHONY: dev-migrate
dev-migrate:
	@echo "> Running database migration ..."
	@docker exec date-apps-be-development ./docker/development/db-migration.sh $(args)

mock:
	@./scripts/generate_mocks.sh

mock-win:
	@powershell scripts/generate_mocks.sh

test-report: 
	go test ./internal/... -v -coverprofile cover.out
	go tool cover -html=cover.out

test:
	go test ./internal/... -v

gen-swagger:
	@echo "Updating API documentation..."
	@swag init -o ${API_DOCS_PATH} -g cmd/webservice/main.go