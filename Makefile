.DEFAULT_GOAL := build-snapshot

.PHONY: lint
lint:
	golangci-lint run
.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -rf ./$(BUILD_DIR)

GORELEASER_PARALLELISM ?= $(shell nproc --ignore=1)
GORELEASER_DEBUG ?= false

export GORELEASER_CURRENT_TAG=$(GIT_TAG)

.PHONY: build-snapshot
build-snapshot: ## Builds a snapshot with goreleaser
build-snapshot:
	goreleaser --debug=$(GORELEASER_DEBUG) \
		build \
		--snapshot \
		--rm-dist \
		--parallelism=$(GORELEASER_PARALLELISM) \
		--single-target \
		--skip-post-hooks

.PHONY: release
release: ## Builds a release with goreleaser
release:
	goreleaser --debug=$(GORELEASER_DEBUG) \
		release \
		--rm-dist \
		--parallelism=$(GORELEASER_PARALLELISM)

.PHONY: release-snapshot
release-snapshot: ## Builds a snapshot release with goreleaser
release-snapshot:
	goreleaser --debug=$(GORELEASER_DEBUG) \
		release \
		--snapshot \
		--skip-publish \
		--rm-dist \
		--parallelism=$(GORELEASER_PARALLELISM)

.PHONY: mod-tidy
mod-tidy: ## Run go mod tidy
	go mod tidy -v -compat=1.17
	go mod verify
