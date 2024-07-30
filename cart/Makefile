GOTOOLCHAIN := $(shell which go)

.PHONY: build
build:
	@echo Using Go at: $(GOTOOLCHAIN)
	@echo "Building the project for ${GOOS}/${GOARCH}..."
	GOOS=${GOOS} GOARCH=${GOARCH} $(GOTOOLCHAIN) build -o build/cart cmd/cart/main.go

.PHONY: mocks
mocks:
	@echo "Generating mocks..."
	@minimock -i 'github.com/BruteMors/marketplace-service/cart/internal/service/cart.ProductService' -o './internal/service/cart/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/cart/internal/service/cart.CartRepository' -o './internal/service/cart/mock' -s '_mock.go'
	@minimock -i 'github.com/BruteMors/marketplace-service/cart/internal/service/cart.LomsService' -o './internal/service/cart/mock' -s '_mock.go'
	@echo "Done!"

.PHONY: coverage
coverage:
	@echo "Running tests with race detection and coverage reporting..."
	@FILES=$$(go list ./... | grep -v /mock | grep -v /test); if [ -z "$$FILES" ]; then \
		echo "No Go files to test after filtering"; \
	else \
		go test -race -covermode=atomic -coverprofile=./coverage/coverage.out -count=5 $$FILES; \
		go tool cover -func=./coverage/coverage.out; \
		echo "Generating HTML coverage report..."; \
		go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html; \
	fi
	@echo "Done!"

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "Done!"

LOCAL_BIN:=$(CURDIR)/bin

PROTOC = PATH="$$PATH:$(LOCAL_BIN)" protoc

LOMS_PROTO_PATH:="../loms/api/proto/loms/v1"
LOMS_GRPC_PATH:="pkg/api/grpc/loms/v1"

.PHONY: .protoc-generate
.protoc-generate: .bin-deps .vendor-proto
#.protoc-generate:
	protoc \
	-I ${LOMS_PROTO_PATH} \
	-I vendor-proto \
	--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
	--go_out ${LOMS_GRPC_PATH} \
	--go_opt paths=source_relative \
	--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
	--go-grpc_out ${LOMS_GRPC_PATH} \
	--go-grpc_opt paths=source_relative \
	../loms/api/proto/loms/v1/loms.proto
	go mod tidy

.PHONY: vendor-proto
.vendor-proto: .vendor-rm  vendor-proto/google/protobuf vendor-proto/validate vendor-proto/google/api vendor-proto/protoc-gen-openapiv2/options

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
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@v0.30.5

.install-proto:
	make vendor-proto/google/protobuf && vendor-proto/validate && vendor-proto/google/api

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
