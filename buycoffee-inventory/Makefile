BUILD_ID := $(shell git rev-parse --short HEAD 2>/dev/null || echo no-commit-id)
IMAGE := anubhavmishra/buycoffee
GOOGLE_PROJECT_ID := serverless-253203

.DEFAULT_GOAL := help
help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean the build
	rm -rf ./build
	packr clean

get: ## Get now cli
	npm install -g now
	go get -u github.com/gobuffalo/packr/packr
	go get .

build-docker-image: ## Build a docker container for the Go service
	mkdir -p ./build/linux/amd64
	packr
	GOOS=linux GOARCH=amd64 go build -v -o ./build/linux/amd64/buycoffee .
	docker build -t $(IMAGE):$(BUILD_ID) .
	docker tag $(IMAGE):$(BUILD_ID) $(IMAGE):latest
	docker tag $(IMAGE):$(BUILD_ID) gcr.io/$(GOOGLE_PROJECT_ID)/$(IMAGE):latest

push: ## push docker image to GCR
	docker push gcr.io/$(GOOGLE_PROJECT_ID)/$(IMAGE):latest

run: ## Build and run the project locally
	mkdir -p ./build
	go build -o ./build/buycoffee && ./build/buycoffee

describe: ## Describe Google cloud run services and domain mappings
	gcloud beta run services list --platform managed
	gcloud beta run domain-mappings list --platform managed

deploy: ## Deploy image to Google cloud run
	gcloud config set project $(GOOGLE_PROJECT_ID)
	gcloud beta run deploy --image gcr.io/$(GOOGLE_PROJECT_ID)/$(IMAGE):latest --platform managed

configure-docker-registry: ## Configure docker to authenticate with GCR.
	gcloud auth configure-docker