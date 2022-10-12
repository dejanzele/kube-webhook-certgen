# Image URL to use all building/pushing image targets
IMG ?= dpejcev/kube-webhook-certgen:latest
# Cockroach build rules.
GO ?= go
# Allow setting of go build flags from the command line.
GOFLAGS :=
# Set to 1 to use static linking for all builds (including tests).
STATIC :=

ifeq ($(STATIC),1)
LDFLAGS += -s -w -extldflags "-static"
endif

lint: # lint code using golangci-lint
	golangci-lint run --max-issues-per-linter=0 --sort-results ./...

test: # run tests using gotestsum
	gotestsum ./...

build-macos-apple: # build executable
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags '$(LDFLAGS)' -v -o kube-webhook-certgen

build-macos-intel: # build executable
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -v -o kube-webhook-certgen

build-intel: # build executable for amd64
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -v -o kube-webhook-certgen

docker-build: test ## Run tests and build docker image.
	docker build -t ${IMG} .

docker-build-quick: ## Build docker image.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}
