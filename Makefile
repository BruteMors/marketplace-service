.PHONY: build-all
build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build

.PHONY: run-all
run-all:
	docker-compose up --force-recreate --build -d

.PHONY: stop-all
stop-all:
	docker-compose down

.PHONY: run-e2e-tests
run-e2e-tests: run-all
	cd ./test/e2e && go test ./...
	make stop-all


MONITORING_COMPOSE_FILE_PATH = ./docker-compose.monitoring.yml

.PHONY: run-all
run-all:
	docker-compose -f $(MONITORING_COMPOSE_FILE_PATH) up --force-recreate --build -d

.PHONY: stop-all
stop-all:
	docker-compose -f $(MONITORING_COMPOSE_FILE_PATH) down
