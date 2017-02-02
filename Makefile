

CMD_DIR=./cmd/android_review_watcher

deps:
	go get -v ./...


build-artifacts: deps
	cd $(CMD_DIR); go build main.go config.go cursor.go
	cd $(CMD_DIR); mkdir -p android_review_watcher
	cd $(CMD_DIR); cp -r main .cursor templates android_review_watcher
	cd $(CMD_DIR); zip -r artifacts.zip android_review_watcher
	cd $(CMD_DIR); rm -r android_review_watcher