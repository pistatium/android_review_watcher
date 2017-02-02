.PHONY: help deps build-artifacts test

CMD_DIR=./cmd/android_review_watcher

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install requirements
	go get -v ./...
	go get github.com/mitchellh/gox

build-artifacts: deps ## Build command tool package
	cd $(CMD_DIR); mkdir -p android_review_watcher
	cd $(CMD_DIR); cp -r main .cursor templates android_review_watcher
	cd $(CMD_DIR); zip -r artifacts.zip android_review_watcher
	cd $(CMD_DIR); rm -r android_review_watcher

test: deps ## Run test
	go test