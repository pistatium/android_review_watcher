package main

import (
	"github.com/BurntSushi/toml"
	"github.com/operando/golack"
)

type Config struct {
	SlackWebHook slackWebhook `toml:"slack_webhook"`
	TargetApps   []targetApp `toml:"target_app"`
}

type slackWebhook struct {
	Url string
}

type targetApp struct {
	PackageName string `toml:"package_name"`
	SlackConf   golack.Slack `toml:"slack_conf"`
}

func LoadConfig(configPath string, config *Config) (*Config, error) {
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return config, err
	}
	return config, nil
}
