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

.PHONY: testacc
testacc: install testacc-awsv4upgrade-simple testacc-awsv4upgrade-full

.PHONY: testacc-awsv4upgrade-simple
testacc-awsv4upgrade-simple: install
	scripts/testacc/awsv4upgrade.sh run simple

.PHONY: testacc-awsv4upgrade-full
testacc-awsv4upgrade-full: install
	scripts/testacc/awsv4upgrade.sh run full

.PHONY: testacc-awsv4upgrade-debug
testacc-awsv4upgrade-debug: install
	scripts/testacc/awsv4upgrade.sh $(ARG)

.PHONY: testacc-all
testacc-all: install
	scripts/testacc/all.sh

.PHONY: check
check: lint test
