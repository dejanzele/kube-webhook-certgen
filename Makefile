# Image URL to use all building/pushing image targets
IMG ?= dpejcev/kube-webhook-certgen:latest

lint: # lint code using golangci-lint
	golangci-lint run --max-issues-per-linter=0 --sort-results ./...

test: # run tests using gotestsum
	gotestsum ./...

build: # build executable
	go build -o kube-webhook-certgen

docker-build: test ## Run tests and build docker image.
	docker build -t ${IMG} .

docker-build-q: ## Build docker image.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}
