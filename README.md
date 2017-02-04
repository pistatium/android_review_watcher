# android_review_watcher

[![CircleCI](https://circleci.com/gh/pistatium/android_review_watcher/tree/master.svg?style=svg)](https://circleci.com/gh/pistatium/android_review_watcher/tree/master)

This application is a tool to get reviews of your Android apps and notify to Slack.
It is using a Google API official API to get reviews.


![capture.png](https://raw.githubusercontent.com/pistatium/android_review_watcher/master/resources/capture.png)

## HowToUse

### Require files

* `config.toml`
  * This file is configuration of your Slack's WebHook, the package names of the application you want to obtain, etc
  * You can set multiple applications at once.
  * Please copy config.toml.sample and custom it.
  
* `service_account.json`
  * This file is OAuth key of GooglePlayAPIs.
  * Please download according to the procedure below.
     * https://developers.google.com/android-publisher/getting_started

### Execute

```bash
# Download Zip (for Linux)
./android_review_watcher -c config.toml -o service_account.json
```

or

```bash
git clone https://github.com/pistatium/android_review_watcher.git
cd android_review_watcher/cmd/android_review_watcher
go get
go build
./android_review_watcher -c config.toml -o service_account.json
```


### Run periodically
This script is not a daemon.
Please run the command at the timing you want.
The same review is controlled not to be reposted.

If you want to use crontab, configure like this

```
0 */3 * * * cd /path/to/android_review_watcher; ./android_review_watcher
```

## References
* https://developers.google.com/android-publisher/api-ref/reviews
