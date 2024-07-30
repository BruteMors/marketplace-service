include .env

GOTOOLCHAIN := $(shell which go)

.PHONY: build
build:
	@echo Using Go at: $(GOTOOLCHAIN)
	@echo "Building the project for ${GOOS}/${GOARCH}..."
	GOOS=${GOOS} GOARCH=${GOARCH} $(GOTOOLCHAIN) build -o build/loms cmd/loms/main.go

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	mkdir -p ./internal/service/order/mock
	mkdir -p ./internal/service/stock/mock
	@minimock -i 'github.com/BruteMors/marketplace-service/loms/internal/service/order.Repository' -o './internal/service/order/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/loms/internal/service/order.StockService' -o './internal/service/order/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/loms/internal/service/order.TxManager' -o './internal/service/order/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/loms/internal/service/order.StatusOutboxRepository' -o './internal/service/order/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/loms/internal/service/stock.Repository' -o './internal/service/stock/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/loms/internal/service/stock.TxManager' -o './internal/service/stock/mock' -s '_mock.go'
	@echo "Done!"

.PHONY: coverage
coverage:
	@echo "Running tests with race detection and coverage reporting..."
	@mkdir -p ./coverage
	@FILES=$$(go list ./... | grep -v /mock | grep -v /test); if [ -z "$$FILES" ]; then \
		echo "No Go files to test after filtering"; \
	else \
		go test -race -covermode=atomic -coverprofile=./coverage/coverage.out -count=5 $$FILES; \
		go tool cover -func=./coverage/coverage.out; \
		echo "Generating HTML coverage report..."; \
		go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html; \
	fi
	@echo "Done!"

.PHONY: intergration_tests
intergration_tests:
	@echo "Running intergration tests..."
	@go test -tags=integration ./...

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "Done!"

LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

LOMS_PROTO_PATH:="api/proto/loms/v1"
LOMS_GRPC_PATH:="pkg/api/grpc/loms/v1"

.PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto
#.protoc-generate:
	mkdir -p api/openapiv2
	protoc \
	-I ${LOMS_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out ${LOMS_GRPC_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out ${LOMS_GRPC_PATH} \
	--go-grpc_opt paths=source_relative \
	--plugin=protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
	--validate_out="lang=go,paths=source_relative:pkg/api/grpc/loms/v1" \
	--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
	--grpc-gateway_out ${LOMS_GRPC_PATH} \
	--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
	--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
	--openapiv2_out api/openapiv2 \
	--openapiv2_opt logtostderr=true \
	api/proto/loms/v1/loms.proto
	go mod tidy

.vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

.PHONY: .serve-swagger
.serve-swagger:
	bin/swagger serve api/openapiv2/loms.swagger.json

.PHONY: .vendor-rm
.vendor-rm:
	rm -rf vendor-proto

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1 && \
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5 && \
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0 && \
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: .install-proto
.install-proto:
	make vendor-proto/google/protobuf && vendor-proto/validate && vendor-proto/google/api

.PHONY: vendor-proto
vendor-proto/google/protobuf:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/protocolbuffers/protobuf vendor-proto/protobuf &&\
	cd vendor-proto/protobuf &&\
	git sparse-checkout set --no-cone src/google/protobuf &&\
	git checkout
	mkdir -p vendor-proto/google
	mv vendor-proto/protobuf/src/google/protobuf vendor-proto/google
	rm -rf vendor-proto/protobuf

vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor-proto/tmp && \
		cd vendor-proto/tmp && \
		git sparse-checkout set --no-cone validate &&\
		git checkout
		mkdir -p vendor-proto/validate
		mv vendor-proto/tmp/validate vendor-proto/
		rm -rf vendor-proto/tmp


vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor-proto/googleapis && \
 	cd vendor-proto/googleapis && \
	git sparse-checkout set --no-cone google/api && \
	git checkout
	mkdir -p  vendor-proto/google
	mv vendor-proto/googleapis/google/api vendor-proto/google
	rm -rf vendor-proto/googleapis

vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem


SWAGGER_UI_VERSION = 3.52.5
SWAGGER_UI_URL = https://github.com/swagger-api/swagger-ui/archive/v$(SWAGGER_UI_VERSION).zip

SWAGGER_UI_DIR = ./public/swagger-ui

.PHONY: swagger-ui
swagger-ui:
	@echo "Downloading Swagger UI $(SWAGGER_UI_VERSION)..."
	curl -L $(SWAGGER_UI_URL) -o swagger-ui.zip
	@echo "Unzipping Swagger UI..."
	unzip swagger-ui.zip
	@echo "Copying Swagger UI to target directory..."
	rm -rf $(SWAGGER_UI_DIR)
	mkdir -p $(SWAGGER_UI_DIR)
	cp -r swagger-ui-$(SWAGGER_UI_VERSION)/dist/* $(SWAGGER_UI_DIR)
	rm -rf swagger-ui-$(SWAGGER_UI_VERSION)
	rm swagger-ui.zip
	@echo "Swagger UI has been set up in $(SWAGGER_UI_DIR)"

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

.PHONY: local-migration-status
local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

.PHONY: local-migration-up
local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

.PHONY: local-migration-down
local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

.PHONY: sqlc-generate
sqlc-generate:
	cd ./internal/repository/postgres/order/sqlc && $(LOCAL_BIN)/sqlc generate
	cd ./internal/repository/postgres/stock/sqlc && $(LOCAL_BIN)/sqlc generate
	cd ./internal/repository/postgres/outbox/sqlc && $(LOCAL_BIN)/sqlc generate
