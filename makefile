SHELL := /usr/bin/bash

ENV ?= dev
ENV_FILE := env/.env.$(ENV)

include $(ENV_FILE)
export

CURDIR := $(shell pwd)
APPNAME := via
COVERAGE_THRESHOLD := 30# coverage lower threshold
BROWSERDIR := "/c/Program Files/Google/Chrome/Application/chrome.exe"

GOBUILD := $(CURDIR)/build
TEMP := $(CURDIR)/tmp
TMP := $(CURDIR)/tmp

# Test execution
test:
	@mkdir -p $(GOBUILD) $(TEMP)
	go test ./... -coverprofile=$(TEMP)/coverage.out

# Will test and verify coverage in order to build
build: test
	@coverage=$(shell go tool cover -func=$(TEMP)/coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$coverage" ]; then \
		echo "Error: Unable to get coverage."; \
		exit 1; \
	fi; \
	coverage_int=$$(echo "$$coverage" | awk '{print int($$1)}'); \
	if [ "$$coverage_int" -ge $(COVERAGE_THRESHOLD) ]; then \
		echo "Coverage successfully achieved, got $$coverage_int% of $(COVERAGE_THRESHOLD)% required."; \
		go build -o $(GOBUILD)/$(APPNAME) ./cmd; \
	else \
		echo "Insufficient coverage, got $$coverage_int% of $(COVERAGE_THRESHOLD)% required."; \
		exit 1; \
	fi;

coverage: test
	go tool cover -html=$(TEMP)/coverage.out -o $(TEMP)/coverage.html
	$(BROWSERDIR) "$(TEMP)/coverage.html"

coverage-no-test:
	"$(BROWSERDIR)" $(TEMP)/coverage.html

clean-local:
	rm -rf $(GOBUILD) $(TEMP)

up:
	@echo "🔼 Levantando servicio con entorno $(ENV_FILE)..."
	docker-compose --env-file $(ENV_FILE) up --build -d
	
# Detiene contenedor
down:
	docker-compose --env-file $(ENV_FILE) down

# Logs del servicio
logs:
	docker-compose --env-file $(ENV_FILE) logs -f

# Rebuild limpio
rebuild: down clean-local up

# Shell dentro del contenedor
sh:
	docker-compose --env-file $(ENV_FILE) exec $(SERVICE_NAME) sh

# Verifica variables cargadas
env:
	@echo "🧪 Variables desde $(ENV_FILE):"
	@cat $(ENV_FILE)