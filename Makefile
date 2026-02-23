SHELL := /bin/bash

API_URL ?= https://paas.null.stackit.run/api
UI_DIR := week_5/ui
API_DIR := week_4/api
GO_CACHE := $(CURDIR)/.gocache

.PHONY: ui-dev ui-build ui-release api-build api-test api-release

ui-dev:
	cd $(UI_DIR) && VITE_API_URL="$(API_URL)" npm run dev -- --host

ui-build:
	cd $(UI_DIR) && npm run build

ui-release:
	./week_5/ui/scripts/rebuild_push_ui.sh --update-kustomize

api-build:
	cd $(API_DIR) && GOCACHE="$(GO_CACHE)" go build ./...

api-test:
	cd $(API_DIR) && GOCACHE="$(GO_CACHE)" go test ./...

api-release:
	./week_4/api/scripts/rebuild_push_api.sh --update-kustomize
