.DEFAULT_GOAL := help

.PHONY: help
help: ## List targets in this Makefile
	@awk '\
		BEGIN { FS = ":$$|:[^#]+|:.*?## "; OFS="\t" }; \
		/^[0-9a-zA-Z_-]+?:/ { print $$1, $$2 } \
	' $(MAKEFILE_LIST) \
		| sort --dictionary-order \
		| column --separator $$'\t' --table --table-wrap 2 --output-separator '    '

.PHONY: fmt
fmt: ## Format code
	goimports -w -local 'github.com/metagram-net/drift' .

.PHONY: build
build: ## Build everything
	go build ./cmd/drift

.PHONY: lint
lint: ## Run linters
	golangci-lint run

.PHONY: test
test: ## Run Go tests
	go test ./...

.PHONY: licensed ## Cache and check dependency licenses
licensed: licensed-cache licensed-check

.PHONY: licensed-check
licensed-check:
	go mod tidy
	licensed status

.PHONY:
licensed-cache:
	go mod tidy
	licensed cache
