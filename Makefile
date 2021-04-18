NAME=split-specs
ROOT=github.com/Shelex/${NAME}
GO111MODULE=on
SHELL=/bin/bash

.PHONY: gql
gql:
	cd api && rm -f generated.go models/*_gen.go &&\
	go run github.com/99designs/gqlgen generate

.PHONY: build
build: 
	make lint
	go build -o ./cmd/client ./

.PHONY: api
api:
	make gql
	make build
	cmd/client

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: deps
deps:
	go mod tidy
	go mod download

.PHONY: dev
dev: 
	dev_appserver.py .

.PHONY: deploy
deploy: 
	make lint
	gcloud app deploy

.PHONY: browse
browse: 
	gcloud app browse