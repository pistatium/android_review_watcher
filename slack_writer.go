package android_review_watcher

import "github.com/operando/golack"

type SlackWriter struct {
	webhook golack.Webhook
	conf    golack.Slack
}

func (s *SlackWriter) Write(p []byte) (n int, err error) {
	payload := golack.Payload{
		Slack: s.conf,
	}
	payload.Slack.Text = string(p)
	golack.Post(payload, s.webhook)
	return len(p), nil
}

func NewSlackWriter(webhook golack.Webhook, conf golack.Slack) *SlackWriter {
	return &SlackWriter{
		webhook: webhook,
		conf:    conf,
	}
}
