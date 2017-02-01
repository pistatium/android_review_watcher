# android_review_watcher

This application is a tool to periodically get a review of Android and notify Slack.
It is using a Google API official API to get reviews.


## HowToUse

### Prepare files
* config.toml
  * A file that specifies Slack's WebHook, the package name of the application you want to obtain, how to display Slack, etc.
  * Copy config.toml.sample and custom it.
  
* client_secret.json
  * OAuth key file to authorize GooglePlayAPI.
  * Please download according to the procedure below.
  * https://developers.google.com/android-publisher/getting_started

### Execute

```bash
android_review_watcher -c config.toml -o client_secret.json
```

or

```bash
cd cmd/android_review_watcher
go get
go run *.go -c config.toml -o client_secret.json
```

## References
* https://developers.google.com/android-publisher/api-ref/reviews
