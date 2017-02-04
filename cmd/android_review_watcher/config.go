package main

import (
	"github.com/BurntSushi/toml"
	"github.com/operando/golack"
	. "github.com/pistatium/android_review_watcher"
)

type Config struct {
	SlackWebHook golack.Webhook `toml:"slack_webhook"`
	TargetApps   []TargetApp    `toml:"target_app"`
}

type TargetApp struct {
	PackageName string       `toml:"package_name"`
	SlackConf   golack.Slack `toml:"slack_conf"`
}

func LoadApps(configPath string) ([]App, error) {
	var config Config
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		return nil, err
	}
	apps := make([]App, len(config.TargetApps))
	for i, target := range config.TargetApps {
		apps[i] = App{
			PackageName: target.PackageName,
			Writer:      NewSlackWriter(config.SlackWebHook, target.SlackConf),
		}
	}

	return apps, nil
}
