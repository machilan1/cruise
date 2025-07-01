ENV_FILE := .env
-include $(ENV_FILE)

sinclude ./scripts/foundation/build.mk
sinclude ./scripts/foundation/rebase.mk
sinclude ./scripts/foundation/undertesting.mk
sinclude ./scripts/foundation/undertesting3.mk
sinclude ./scripts/foundation/undertesting4.mk
sinclude ./scripts/project/project.mk

# ==============================================================================
# Define environment variables

# get current directory without the full path
CURRENT_DIR := $(notdir $(patsubst %/,%,$(CURDIR)))
export COMPOSE_PROJECT_NAME := $(CURRENT_DIR)

# ==============================================================================
# Install dependencies

dev-brew:
	brew update
	brew list lefthook || brew install lefthook
	lefthook install

dev-gotooling: dev-gotooling-dev dev-gotooling-ci

dev-gotooling-dev:
	go install github.com/divan/expvarmon@latest
	go install github.com/air-verse/air@latest
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-gotooling-ci:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================

# lint uses the same linter as CI and tries to report the same results running
# locally. There is a chance that CI detects linter errors that are not found
# locally, but it should be rare.
lint:
	golangci-lint run --config .golangci.yaml && go run ./cmd/enhancelint ./internal/...

vuln-check:
	govulncheck ./...

fmt:
	go fmt ./...
	goimports -l -w cmd internal
	gofumpt -l -w cmd internal

# diff-check runs git-diff and fails if there are any changes.
.PHONY: diff-check
diff-check:
	@FINDINGS="$$(git status -s -uall)" ; \
		if [ -n "$${FINDINGS}" ]; then \
			echo "Changed files:\n\n" ; \
			echo "$${FINDINGS}\n\n" ; \
			echo "Diffs:\n\n" ; \
			git diff ; \
			git diff --cached ; \
			exit 1 ; \
		fi

.PHONY: generate
generate:
	@go generate ./...

.PHONY: generate-check
generate-check: generate diff-check

# ==============================================================================
# Testing

test:
	go test -count=1 -shuffle=on -timeout=5m ./internal/...

test-acc:
	go test -count=1 -shuffle=on -timeout=10m -race ./internal/... -coverprofile=coverage.out

test-coverage:
	go tool cover -func=./coverage.out

# ==============================================================================
# Running from within docker

dev-up:
	docker compose -f ./build/docker-compose.yaml up -d --build --remove-orphans

dev-up-no-build:
	docker compose -f ./build/docker-compose.yaml up -d --remove-orphans

dev-down:
	docker compose -f ./build/docker-compose.yaml down -v

# ==============================================================================
# Build containers

dev-build: dev-build-api dev-build-migrate dev-build-seed

dev-build-api:
	docker compose -f ./build/docker-compose.yaml build api

dev-build-migrate:
	docker compose -f ./build/docker-compose.yaml build migrate

dev-build-seed:
	docker compose -f ./build/docker-compose.yaml build seed

# ==============================================================================
# Build & restart containers

dev-update: dev-update-api

dev-update-api:
	docker compose -f ./build/docker-compose.yaml up -d --build --no-deps api

# ==============================================================================
# Logs

dev-logs:
	docker compose -f ./build/docker-compose.yaml logs api -f --no-log-prefix --tail=100 | go run cmd/logfmt/main.go

# ==============================================================================
# Code generation

codegen:
	@read -p "Please enter the domain name you want to add: " name; \
	if [ -d "./internal/app/domain/$$name" ] || [ -d "./internal/business/domain/$$name" ] || [ -d "./internal/business/domain/$$name/stores/$$name""db" ]; then \
		echo "The domain '$$name' already exists. Entering one to many mode..."; \
		read -p "Please enter the abbreviation name you want to add: " abbr; \
		read -p "Please enter the plural you want to add: " plur; \
		read -p "Please choose whether to record createdAt / updatedAt: " cutime; \
		read -p "Please choose whether to use pagination (Y/n): " pag; \
		read -p "Please choose whether to use optimistic concurrency control (Y/n): " occ; \
		read -p "Please choose whether to enable soft delete (Y/n): " del; \
		read -p "Please specify one-to-many name: " otmn; \
		read -p "Please specify plural one-to-many name: " otmnp; \
		go run ./cmd/codegen $$name $$abbr $$plur $$cutime $${pag:-Y} $${occ:-Y} $${del:-Y} "Y" $$otmn $$otmnp; \
	else \
		read -p "Please enter the abbreviation name you want to add: " abbr; \
		read -p "Please enter the plural you want to add: " plur; \
		read -p "Please choose whether to record createdAt / updatedAt: " cutime; \
		read -p "Please choose whether to use pagination (Y/n): " pag; \
		read -p "Please choose whether to use optimistic concurrency control (Y/n): " occ; \
		read -p "Please choose whether to enable soft delete (Y/n): " del; \
		go run ./cmd/codegen $$name $$abbr $$plur $$cutime $${pag:-Y} $${occ:-Y} $${del:-Y} "n" "" ""; \
	fi

codegen2:
	@read -p "Please enter the domain name you want to add: " name; \
	if [ -d "./internal/app/domain/$$name" ] || [ -d "./internal/business/domain/$$name" ] || [ -d "./internal/business/domain/$$name/stores/$$name""db" ]; then \
		echo "Error: The domain '$$name' already exists. Aborting."; \
		exit 1; \
	else \
		read -p "Please enter the abbreviation name you want to add: " abbr; \
		read -p "Please enter the plural you want to add: " plur; \
		go run ./cmd/codegen $$name $$abbr $$plur q n; \
	fi

# ======================================================================================================================

dev-migrate-up:
	docker compose -f ./build/docker-compose.yaml run --build --no-deps migrate | go run cmd/logfmt/main.go

dev-migrate-down:
	@read -p "Please enter the number of migrations you want to rollback: " num; \
	docker compose -f ./build/docker-compose.yaml run --build --no-deps -e MIGRATION_DOWN=true -e MIGRATION_VERSION=$$num migrate | go run cmd/logfmt/main.go

dev-migrate-down-all:
	docker compose -f ./build/docker-compose.yaml run --build --no-deps -e MIGRATION_DOWN=true -e MIGRATION_VERSION=0 migrate | go run cmd/logfmt/main.go

dev-run-seed:
	@read -p "Please enter the version of the seed you want to run (leave empty to run all): " version; \
	if [ -z "$$version" ]; then \
		docker compose -f ./build/docker-compose.yaml run --build --no-deps seed | go run cmd/logfmt/main.go; \
	else \
		docker compose -f ./build/docker-compose.yaml run --build --no-deps -e SEED_VERSION=$$version seed | go run cmd/logfmt/main.go; \
	fi

dev-run-seed-all:
	docker compose -f ./build/docker-compose.yaml run --build --no-deps seed | go run cmd/logfmt/main.go;

dev-resetdb:
	-@$(MAKE) dev-migrate-down-all
	-@$(MAKE) dev-migrate-up
	-@$(MAKE) dev-run-seed-all
# ======================================================================================================================
# Find hints


find-table:
	@read -p "Please enter the table name you want to find: " table; \
	if [ -z "$$table" ]; then \
		echo "Error: Table name cannot be empty. Aborting."; \
		exit 1; \
	else \
		echo "find in migration files:"; \
		echo "------------------------------------------------------------------------------------------------------" ; \
		echo "" ; \
		rg "$$table" ./migrations; \
		echo "" ; \
		echo "find in seed files:"; \
		echo "------------------------------------------------------------------------------------------------------" ; \
		echo "" ; \
		rg "$$table" ./cmd/seed; \
		echo "" ; \
		echo "find in code files:"; \
		echo "------------------------------------------------------------------------------------------------------" ; \
		echo "" ; \
		rg "$$table" ./internal; \
		echo "" ; \
	fi

# ======================================================================================================================
# Cloud initialization

cloud-init:
	export PROJECT_ID=${PROJECT_ID} && chmod +x ./deploy/cloud-init.sh && ./deploy/cloud-init.sh

cloud-init-db:
	export PROJECT_ID=${PROJECT_ID} && export INFRA_PROJECT_ID=${INFRA_PROJECT_ID} && chmod +x ./deploy/db-init.sh && ./deploy/db-init.sh

# DO THIS for other users
# 1. cloudsqlsuperuser
# 2. ALTER ROLE "mayainfo.co.ltd@gmail.com" SET ROLE cloudsqlsuperuser;

sql-permission:
	PGPASSWORD=$$(gcloud auth print-access-token) gcloud beta sql connect primary --user=$$(gcloud config get-value account) --quiet --database=${PROJECT_ID}-db

cloud-db-permission:
	@export PROJECT_ID=${PROJECT_ID} && export INFRA_PROJECT_ID=${INFRA_PROJECT_ID} && chmod +x ./deploy/db-perm.sh && ./deploy/db-perm.sh

# ======================================================================================================================
# Cloud database management

cloudsql-proxy:
	@docker run -p $$CLOUD_SQL_DB_PORT:5432 --rm gcr.io/cloud-sql-connectors/cloud-sql-proxy:latest "$$CLOUD_SQL_CONNECTION_NAME" \
			--address 0.0.0.0 \
			--auto-iam-authn \
			--token "$(shell gcloud auth application-default print-access-token)" \
			--login-token "$(shell gcloud sql generate-login-token)"

# Must run `cloudsql-proxy` before running this command
cloudsql-psql:
	@PGSSLMODE=disable psql -h 127.0.0.1 -p $$CLOUD_SQL_DB_PORT -d $$CLOUD_SQL_DB_NAME -U $(shell gcloud config get-value account)


#cloudsql-migrate-up:
	#go run cmd/migrate/main.go --db-name=$$CLOUD_SQL_DB_NAME --db-cloud-sql-connection-name=$$CLOUD_SQL_CONNECTION_NAME --db-user=$(shell gcloud config get-value account) --db-host=127.0.0.1:$$CLOUD_SQL_DB_PORT | go run cmd/logfmt/main.go

#cloudsql-migrate-down:
#	go run cmd/migrate/main.go --db-name=$$CLOUD_SQL_DB_NAME --db-cloud-sql-connection-name=$$CLOUD_SQL_CONNECTION_NAME --db-user=$(shell gcloud config get-value account) --db-host=127.0.0.1:$$CLOUD_SQL_DB_PORT --migration-down=true --migration-version=0 | go run cmd/logfmt/main.go

#cloudsql-seed:
#	go run cmd/seed/main.go --db-name=$$CLOUD_SQL_DB_NAME --db-cloud-sql-connection-name=$$CLOUD_SQL_CONNECTION_NAME --db-user=$(shell gcloud config get-value account) --db-host=127.0.0.1:$$CLOUD_SQL_DB_PORT | go run cmd/logfmt/main.go
