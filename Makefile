NAME := tfedit

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go mod download

.PHONY: build
build: deps
	go build -o bin/$(NAME)

.PHONY: install
install: deps
	go install

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test: build
	go test ./...

.PHONY: testacc-awsv4upgrade
testacc-awsv4upgrade: install
	scripts/testacc/awsv4upgrade.sh $(ARG)

.PHONY: testacc-all
testacc-all: install
	scripts/testacc/all.sh

.PHONY: check
check: lint test
