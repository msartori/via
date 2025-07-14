MAKEFLAGS += --no-print-directory
ENV ?= local
ENV_FILE := env/.env.$(ENV)
LT_API_LOG=lt_api.log
LT_FRONT_LOG=lt_front.log

include $(ENV_FILE)
export

CURDIR := $(shell pwd)/api
API_IMAGE_NAME := $(APPNAME)-$(ENV)-api-img
WEB_IMAGE_NAME := $(APPNAME)-$(ENV)-web-img
COVERAGE_THRESHOLD := 40# coverage lower threshold
BROWSERDIR := open -a "Google Chrome"

GOBUILD := $(CURDIR)/build
TEMP := $(CURDIR)/tmp

export GOCACHE=${TEMP}

# Test execution
test:
	@mkdir -p $(GOBUILD) $(TEMP)
	@cd api && \
	PKG_LIST=$$(go list ./... | grep -v -E '/mock|/internal/ent') && \
	go test $$PKG_LIST -coverpkg=$$(echo $$PKG_LIST | tr ' ' ',') -coverprofile=$(TEMP)/coverage.out
	@echo "Total test coverage:"
	cd api && go tool cover -func=$(TEMP)/coverage.out | grep total
	cd api && go tool cover -html=$(TEMP)/coverage.out -o $(TEMP)/coverage.html


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
	@echo "$$coverage"

coverage: test
	$(BROWSERDIR) "$(TEMP)/coverage.html"

coverage-no-test:
	$(BROWSERDIR) "$(TEMP)/coverage.html"

clean-local:
	rm -rf $(GOBUILD) $(TEMP)

localtunnel-start:
	@echo "Starting LocalTunnel API ($(API_REST_PORT))..."
	@nohup npx localtunnel --subdomain via-$(ENV)-api --port $(API_REST_PORT) &

	@echo "Starting LocalTunnel API ($(API_SSE_PORT))..."
	@nohup npx localtunnel --subdomain via-$(ENV)-sse-api --port $(API_SSE_PORT) &

	@echo "Starting LocalTunnel WEB ($(WEB_PORT))..."
	@nohup npx localtunnel --subdomain via-$(ENV)-web --port $(WEB_PORT) &

	@sed -i '' -e '/^VITE_API_URL=/d' web/env/.env.$(ENV)
	@echo "VITE_API_URL=https://via-$(ENV)-api.loca.lt" >> web/env/.env.$(ENV)

	@sed -i '' -e '/^VITE_WEB_URL=/d' web/env/.env.$(ENV)
	@echo "VITE_WEB_URL=https://via-$(ENV)-web.loca.lt" >> web/env/.env.$(ENV)

	@sed -i '' -e '/^VITE_API_SSE_URL=/d' web/env/.env.$(ENV)
	@echo "VITE_API_SSE_URL=https://via-$(ENV)-sse-api.loca.lt" >> web/env/.env.$(ENV)

localtunnel-stop:
	@echo "Killing LocalTunnel process..."
	@pkill -f "localtunnel" || true

start-api:
	docker-compose --env-file $(ENV_FILE) up --build -d api

start-web:
	docker-compose --env-file $(ENV_FILE) up --build -d web

up:
	echo "🔼 Starting service on ENV: $(ENV_FILE)..."
	#docker-compose --env-file $(ENV_FILE) up --build -d
	make start-api
ifeq ($(SKIP_TUNNEL),)
	make localtunnel-start
endif
	make start-web
	
# Stop container
down:
	@echo "🔽 Stopping service on ENV: $(ENV_FILE)..."
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