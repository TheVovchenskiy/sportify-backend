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

.PHONY: migration-up
migration-up:
	docker exec -it sportify-backend-backend_sportify-1 ./migrate -database postgres://postgres:postgres@postgres:5432/sportify?sslmode=disable -path ./migrations up

.PHONY: migration-down
migration-up:
	docker exec -it sportify-backend-backend_sportify-1 ./migrate -database postgres://postgres:postgres@postgres:5432/sportify?sslmode=disable -path ./migrations down

.PHONY: migration-up-reserve
migration-up-reserve:
	docker exec -it sportify-backend-backend_sportify-1 ./migrate -database postgres://postgres:postgres@localhost:5432/sportify?sslmode=disable -path ./migrations up

.PHONY: fill-db
fill-db:
	docker cp sportify/db/fill.sql sportify-backend-postgres-1:fill.sql
	docker exec -it sportify-backend-postgres-1 psql -U postgres -d sportify -f fill.sql

up_names=
.PHONY: docker-compose-up
docker-compose-up:
	docker compose --project-directory . -f deploy/docker-compose.yaml up -d $(up_names)

.PHONY: docker-compose-down
docker-compose-down:
	docker compose --project-directory . -f deploy/docker-compose.yaml down

.PHONY: docker-compose-build
docker-compose-build:
	docker compose --project-directory . -f deploy/docker-compose.yaml build

names=backend_sportify
.PHONY: docker-compose-logs
docker-compose-logs:
	docker compose --project-directory . -f deploy/docker-compose.yaml logs -f $(names)

.PHONY: gen-cert
gen-cert:
	openssl req -x509 -sha256 -nodes -subj "/C=RU/ST=Moscow/L=Moscow/O=Sportify/CN=localhost" -newkey rsa:2048 -days 365 -keyout config/localhost.key -out config/localhost.crt
