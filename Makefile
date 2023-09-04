.PHONY: compose_up compose_down build run test lint swag

.DEFAULT_GOAL := compose_up

APP := user-segmentation

compose_up:
	docker compose up -d --build --remove-orphans
	docker logs -f ${APP}

compose_down:
	docker compose down --volumes --remove-orphans
	docker image prune -f --filter="dangling=true"

build:
	go mod vendor
	go build -mod=vendor -o ./bin/${APP} ./cmd/${APP}

run: build
	./bin/${APP}

test:
	go test -v -coverpkg=./... -coverprofile=./test_coverage/cover.out ./...
	go tool cover -html=./test_coverage/cover.out -o ./test_coverage/cover.html

lint:
	golangci-lint -c golangci.yml run -v ./...

swag:
	swag init -g internal/app/app.go --parseInternal