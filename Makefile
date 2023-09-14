include ./config/.env
export

.DEFAULT_GOAL := compose_run

.PHONY: compose_run compose_down psql build run test test_cover lint swag

compose_run:
	docker compose run -it --rm --build --remove-orphans --service-ports --name $(APP) $(APP)

compose_down:
	docker compose down --volumes --rmi local

psql:
	docker exec -it $(APP) \
	psql postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@\
	$(POSTGRES_ADDR)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

build:
	go build -mod=vendor -o ./bin/$(APP) ./cmd/$(APP)

run: build
	./bin/$(APP)

test:
	go test -v -race ./...

test_cover:
	go test -v -coverpkg=./... -coverprofile=./test_coverage/cover.out ./...
	go tool cover -html=./test_coverage/cover.out -o ./test_coverage/cover.html

lint:
	golangci-lint -c golangci.yml run -v ./...

swag:
	swag init -g internal/app/app.go --parseInternal