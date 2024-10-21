.PHONY: toolchain-migrate
toolchain-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: toolchain
toolchain: toolchain-migrate
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install go.uber.org/mock/mockgen@latest

.PHONY: lint
lint:
	cd sportify && \
	golangci-lint run --fix --timeout=3m ./...

.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir ./sportify/db/migrations $(name)

.PHONY: docker-compose-up
docker-compose-up:
	docker compose -f deploy/docker-compose.yaml up -d

.PHONY: docker-compose-down
docker-compose-down:
	docker compose -f deploy/docker-compose.yaml down

.PHONY: docker-compose-build
docker-compose-build:
	docker compose -f deploy/docker-compose.yaml build
