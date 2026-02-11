vet: ## Run go vet
	./tool/go vet ./...

tidy: ## Run go mod tidy
	./tool/go mod tidy

updatedeps: ## Update depaware deps for tailscaled
	# depaware (via x/tools/go/packages) shells back to "go", so make sure the "go"
	# it finds in its $$PATH is the right one.
	PATH="$$(./tool/go env GOROOT)/bin:$$PATH" ./tool/go run github.com/tailscale/depaware --update --vendor --internal \
		tailscale.com/cmd/tailscaled
	PATH="$$(./tool/go env GOROOT)/bin:$$PATH" ./tool/go run github.com/tailscale/depaware --update --file=depaware-min.txt --goos=linux --tags="$$(./tool/go run ./cmd/featuretags --min)" --vendor --internal \
		tailscale.com/cmd/tailscaled

depaware: ## Run depaware checks for tailscaled
	# depaware (via x/tools/go/packages) shells back to "go", so make sure the "go"
	# it finds in its $$PATH is the right one.
	PATH="$$(./tool/go env GOROOT)/bin:$$PATH" ./tool/go run github.com/tailscale/depaware --check --vendor --internal \
		tailscale.com/cmd/tailscaled
	PATH="$$(./tool/go env GOROOT)/bin:$$PATH" ./tool/go run github.com/tailscale/depaware --check --file=depaware-min.txt --goos=linux --tags="$$(./tool/go run ./cmd/featuretags --min)" --vendor --internal \
		tailscale.com/cmd/tailscaled

build: ## Build tailscaled for current platform
	./build_dist.sh ./cmd/tailscaled

build-mesh: ## Build single mesh binary (daemon + CLI)
	./tool/go build -tags ts_include_cli -o mesh ./cmd/tailscaled

buildwindows: ## Build tailscaled for windows/amd64
	GOOS=windows GOARCH=amd64 ./build_dist.sh ./cmd/tailscaled

build386: ## Build tailscaled for linux/386
	GOOS=linux GOARCH=386 ./build_dist.sh ./cmd/tailscaled

buildlinuxarm: ## Build tailscaled for linux/arm
	GOOS=linux GOARCH=arm ./build_dist.sh ./cmd/tailscaled

buildplan9: ## Build tailscaled for plan9/amd64
	GOOS=plan9 GOARCH=amd64 ./build_dist.sh ./cmd/tailscaled

buildlinuxloong64: ## Build tailscaled for linux/loong64
	GOOS=linux GOARCH=loong64 ./build_dist.sh ./cmd/tailscaled

check: vet depaware build ## Perform basic checks and compilation tests



.PHONY: generate
generate: ## Generate code
	./tool/go generate ./...

help: ## Show this help
	@echo ""
	@echo "Specify a command. The choices are:"
	@echo ""
	@grep -hE '^[0-9a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[0;36m%-20s\033[m %s\n", $$1, $$2}'
	@echo ""
.PHONY: help

.DEFAULT_GOAL := help
