SHELL=/bin/bash

PROJECT_ROOT := $(shell pwd)
include .env
export

.PHONY: up
up: build network run

.PHONY: build
build:
	docker build \
			--target run \
			-f ${PROJECT_ROOT}/deployments/prod/Dockerfile \
			-t infura-image .

.PHONY: network
network:
	# exit code 0 means network exists
	docker network inspect infura-net &> /dev/null; \
	if [ $$? -ne 0 ]; then docker network create infura-net ; fi

.PHONY: run
run:
	docker run -d \
		  --name infura-rest \
		  --network infura-net \
		  --env-file ${PROJECT_ROOT}/.env \
		  -p 8080:8080 infura-image

.PHONY: tests
tests:
	go clean -testcache && go test -v ./...

.PHONY: bench
bench:
	go clean -testcache && go test -bench . ./... -count 10

.PHONY: start
start:
	docker start infura-rest

.PHONY: stop
stop:
	docker stop infura-rest

.PHONY: down
down:
	make stop
	docker rm infura-rest
	docker network rm infura-net

.PHONY: rebuild
rebuild:
	docker stop infura-rest
	docker rm infura-rest
	make build
	make run