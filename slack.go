package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

//const slackWebhook = "your_slack_webhook_url"

func sendSlackNotification(message string) {
	api := slack.NewWebHook(slackWebhook)

	_, _, err := api.PostMessage("your_slack_channel", slack.MsgOptionText(message, false))
	if err != nil {
		fmt.Printf("Failed to post message to Slack: %v", err)
		return
	}

	fmt.Println("Message successfully sent to Slack")
}
