
lint:
	go install golang.org/x/lint/golint@latest
	golint ./...

up-services:
	docker-compose up -d adminer

build:
	go build -o api

up:
	make up-services
	make generate-swagger
	make run

run:
	make build
	./api

test:
	go test -v ./...

generate-mocks:
	mockery --disable-version-string --output ./mocks/services --dir services --all
	mockery --disable-version-string --output ./mocks/repository --dir repository --all

generate-swagger:
	swag i