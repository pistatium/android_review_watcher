.PHONY: help deps build-artifacts upload-artifacts test

CMD_DIR=./cmd/android_review_watcher

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Install requirements
	go get -v ./...
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr

build-artifacts: deps ## Build command tool package
	cd $(CMD_DIR); mkdir -p android_review_watcher
	cd $(CMD_DIR); go build main.go config.go
	cd $(CMD_DIR); cp -r .cursor templates android_review_watcher
	cd $(CMD_DIR); cp main android_review_watcher/android_review_watcher
	cp config.toml.sample $(CMD_DIR)/android_review_watcher/config.toml
	cd $(CMD_DIR); zip -r artifacts.zip android_review_watcher
	cd $(CMD_DIR); rm -r android_review_watcher

upload-artifacts: ## Upload artifacts to github
	cd $(CMD_DIR); ghr -t $(GITHUB_TOKEN) -u $(USERNAME) -r $(CIRCLE_PROJECT_REPONAME) $(CIRCLE_TAG) artifacts.zip

test: deps ## Run test
	go test -cover