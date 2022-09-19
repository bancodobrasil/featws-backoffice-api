
lint:
	go install golang.org/x/lint/golint@latest
	@make run-on-our-code-directories ARGS="golint"

up-services:
	docker-compose up -d adminer

build:generate-swagger
	go build -o api

up:up-services
	make run

run:build
	./api

test:
	@make run-on-our-code-directories ARGS="go test -v -coverprofile=coverage.out"

run-on-our-code-directories:
	@echo "${ARGS} <our-code-directories>"
	@make our-code-directories | xargs -n 1 $(ARGS)
our-code-directories:
	@go list ./... | grep -v /vendor | grep -v /docs
verify:test
	make lint
generate-mocks:
#	Install testify on https://github.com/stretchr/testify
#	Install mockery on https://github.com/vektra/mockery
	mockery --disable-version-string --output ./mocks/services --dir services --all
	mockery --disable-version-string --output ./mocks/repository --dir repository --all

generate-swagger:
#   Install swag on https://github.com/swaggo/swag
	swag i