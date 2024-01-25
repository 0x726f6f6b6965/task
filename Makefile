PROJECTNAME := $(shell basename "$(PWD)")
include .env
export $(shell sed 's/=.*//' .env)

# Protobuf
## proto-lint: Check protobuf rule
.PHONY: proto-lint
proto-lint:
	@buf lint

## proto-gen: Generate golang files based on protobuf
.PHONY: proto-gen
proto-gen:
	@buf generate

## proto-check-breaking: Check protobuf breaking
.PHONY: proto-check-breaking
proto-check-breaking:
	@buf breaking --against '.git#branch=main' --error-format=json | jq .

## proto-clean: Clean the golang files which are generated based on protobuf
.PHONY: proto-clean
proto-clean: 
	@find protos -type f -name "*.go" -delete

# Dockerfile
## gen-images: Generate serivces' image
.PHONY: gen-images
gen-images:
	@docker build --tag task-svc:$(shell git describe --tags --abbrev=0) -f ./build/Dockerfile .

## service-up: Run the all components by deployment/compose.yaml
.PHONY: service-up
service-up:
	@docker-compose  -f ./deployment/compose.yaml --project-directory . up

## service-down: Docker-compose down
.PHONY: service-down
service-down:
	@docker-compose -f ./deployment/compose.yaml --project-directory . down

## test-go: Test go file and show the coverage
.PHONY: test-go
test-go:
	@go test --coverprofile=coverage.out ./... 
	@go tool cover -html=coverage.out  

## help: Print usage information
.PHONY: help
help: Makefile
	@echo
	@echo "Choose a command to run in $(PROJECTNAME)"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo