ANDROID_CLIENT_DIR := android-client

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: generate
generate: ## Generate code
	go generate ./...

.PHONY: analyst
analyst: ## Build MESH analyst client
	./analyst/build.sh

.PHONY: clean
clean: ## Clean up build artifacts
	cd $(ANDROID_CLIENT_DIR) && $(MAKE) clean
	rm -f analyst/mesh

.PHONY: help
help: ## Show this help
	@echo ""
	@echo "Specify a command. The choices are:"
	@echo ""
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""

.DEFAULT_GOAL := help
