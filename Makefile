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

.PHONY: keys
keys: 
	mkdir keys
	openssl genrsa -out ./keys/app.rsa 1024
	openssl rsa -in app.rsa -pubout > ./keys/app.rsa.pub

.PHONY: deps
deps:
	go mod tidy
	go mod download

.PHONY: dev
dev: 
	dev_appserver.py app_dev.yaml --clear_datastore=yes --port=9999

.PHONY: web-dev
web-dev: 
	cd web && yarn start

.PHONY: web-deps
web-deps: 
	cd web && yarn install

.PHONY: web-build
web-build: 
	cd web && rm -rf build && yarn build

.PHONY: deploy
deploy: 
	make lint
	gcloud app deploy --verbosity=info

.PHONY: browse
browse: 
	gcloud app browse

.PHONY: datastore-index
datastore-index:
	gcloud datastore indexes create $(PWD)/index.yaml