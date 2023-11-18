GOPATH := $(shell go env GOPATH)
BINARY_NAME=technoStore
MIGRATE_BIN := $(GOPATH)/bin/migrate

DB_NAME := $(or $(DB_NAME),technoStore)
DB_USER := $(or $(DB_USER),postgres)
DB_PASS := $(or $(DB_PASSWORD),docker)
DB_HOST := $(or $(DB_HOST),localhost)
DB_PORT := $(or $(DB_PORT),5432)

$(MIGRATE_BIN):
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

$(MOCK_BIN):
	go install go.uber.org/mock/mockgen@latest

run:
	go run ./cmd/server/main.go

lru:
	go run ./cmd/problem_solving/main.go

test:
	go test -v -cover -short ./...

swag:
	swag init -g ./cmd/server/main.go -o ./docs

clean:
	go clean

mock-db: $(MOCK_BIN)
	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/brand.go techno-store/internal/domain/definition BrandRepository
	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/category.go techno-store/internal/domain/definition CategoryRepository
	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/supplier.go techno-store/internal/domain/definition SupplierRepository
	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/product.go techno-store/internal/domain/definition ProductRepository
	mockgen -package mockdb -destination internal/infrastructure/datastores/mockdb/productStock.go techno-store/internal/domain/definition ProductStockRepository

migrate-up: $(MIGRATE_BIN)
	migrate -source file://db/migrations -database postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable -verbose up

migrate-down: $(MIGRATE_BIN)
	migrate -source file://db/migrations -database postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable -verbose down

.PHONY: run lru test swag clean mock migrate-up migrate-down mock-db
