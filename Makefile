# Image URL to use all building/pushing image targets
IMG ?= dpejcev/kube-webhook-certgen:latest

lint: # lint code using golangci-lint
	golangci-lint run

test: # run tests using gotestsum
	gotestsum ./...

docker-build: test ## Run tests and build docker image.
	docker build -t ${IMG} .

docker-build-q: ## Build docker image.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}
