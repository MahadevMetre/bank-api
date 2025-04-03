include .env.development

# Export all variables from .env
export $(shell sed 's/=.*//' .env.development)

DB_STRING = postgres://$(POSTGRES_USER_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DATABASE)?sslmode=disable

run:
	go run main.go
	
build:
	GOOS=linux GOARCH=amd64 go build -tags=jsoniter -ldflags="-w -s" -o bankapi

# go build -gcflags="-m"

# $(MAKE) build-upload
format:
	@echo "Running formater..."
	gofmt -w .

build-upload:
	sudo scp -i ~/Desktop/Keypairs/paydoh-key.pem ./bankapi ubuntu@172.31.4.28:~/bankapi

mac-build:
	GOOS=darwin GOARCH=arm64 go build -o main main.go

upload-migration:
	sudo scp -i ~/Desktop/Workspace/paydoh/aws-ec2/paydoh-key.pem -r ./migrations ubuntu@172.31.4.28:~/bankapi/

upload-templates:
	sudo scp -i ~/Desktop/Workspace/paydoh/aws-ec2/paydoh-key.pem -r ./templates ubuntu@172.31.4.28:~/bankapi/

git-submodule-update:
	git submodule update --recursive

swagger-generate:
	swag init

# goose migrations
goose-create-new:
	goose -dir ./migrations create $(MIGRATION_NAME) sql

create_migration:
	@read -p "Enter migration name: " name; \
	goose -dir ./migrations create $$name sql

goose-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://$(POSTGRES_USER_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DATABASE)?sslmode=disable" goose -dir='./migrations' up

goose-up-by-one:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://$(POSTGRES_USER_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DATABASE)?sslmode=disable" goose -dir='./migrations' up-by-one

goose-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://$(POSTGRES_USER_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DATABASE)?sslmode=disable" goose -dir='./migrations' down

goose-down-to:
	goose -dir=./migrations postgres "$(DB_STRING)" down-to $(VERSION)

goose-status:
	goose -dir=./migrations postgres "$(DB_STRING)" status

goose-version:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="postgres://$(POSTGRES_USER_NAME):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DATABASE)?sslmode=disable" goose -dir='./migrations' version

test_run:
	GIN_MODE=test go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

test_coverage:
	@go test -cover ./...
	@go tool cover -html=coverage.out

lint:
	golangci-lint run

import-ifsc-code:
	go run data/import_ifsc_code.go