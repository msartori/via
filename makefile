SHELL := /usr/bin/bash
MAKEFLAGS += --no-print-directory
ENV ?= dev
ENV_FILE := env/.env.$(ENV)

include $(ENV_FILE)
export

CURDIR := $(shell pwd)/api
API_IMAGE_NAME := $(APPNAME)-$(ENV)-api-img
WEB_IMAGE_NAME := $(APPNAME)-$(ENV)-web-img
COVERAGE_THRESHOLD := 20# coverage lower threshold
BROWSERDIR := "/c/Program Files/Google/Chrome/Application/chrome.exe"

GOBUILD := $(CURDIR)/build
TEMP := $(CURDIR)/tmp

export GOCACHE=${TEMP}

# Test execution
test:
	@mkdir -p $(GOBUILD) $(TEMP)
	cd api && go test ./... -coverprofile=$(TEMP)/coverage.out


# Will test and verify coverage in order to build
build: test
	@coverage=$(shell cd api && go tool cover -func=$(TEMP)/coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ -z "$$coverage" ]; then \
		echo "Error: Unable to get coverage."; \
		exit 1; \
	fi; \
	coverage_int=$$(echo "$$coverage" | awk '{print int($$1)}'); \
	if [ "$$coverage_int" -ge $(COVERAGE_THRESHOLD) ]; then \
		echo "Coverage successfully achieved, got $$coverage_int% of $(COVERAGE_THRESHOLD)% required."; \
		cd api && go build -o $(GOBUILD)/$(APPNAME) ./cmd; \
	else \
		echo "Insufficient coverage, got $$coverage_int% of $(COVERAGE_THRESHOLD)% required."; \
		exit 1; \
	fi;

coverage: test
	cd api && GOCACHE=/tmp/gocache go tool cover -html=$(TEMP)/coverage.out -o $(TEMP)/coverage.html
	$(BROWSERDIR) "$(TEMP)/coverage.html"

coverage-no-test:
	$(BROWSERDIR) "$(TEMP)/coverage.html"

clean-local:
	rm -rf $(GOBUILD) $(TEMP)

up:
	echo "🔼 Starting service on ENV: $(ENV_FILE)..."
	docker-compose --env-file $(ENV_FILE) up --build -d
	
# Stop container
down:
	docker-compose --env-file $(ENV_FILE) down

# Chech service logs
logs:
	docker-compose --env-file $(ENV_FILE) logs -f

# Rebuild 
rebuild: down clean-local remove-image up

# Shell in container
sh:
	docker-compose --env-file $(ENV_FILE) exec api sh

# Verify vars in
env:
	@echo "🧪 Vars in $(ENV_FILE):"
	@cat $(ENV_FILE)

# clean image
remove-image:
	@echo "🗑️  Deleting image $(API_IMAGE_NAME)..."
	docker rmi -f $(API_IMAGE_NAME) || true
	@echo "🗑️  Deleting image $(WEB_IMAGE_NAME)..."
	docker rmi -f $(WEB_IMAGE_NAME) || true

.PHONY: env logs