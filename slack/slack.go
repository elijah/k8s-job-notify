package slack

import (
	"os"
	"fmt"
	"github.com/achilles-git/k8s-job-notify/env"
	"github.com/ashwanthkumar/slack-go-webhook"
)

type requestBody struct {
	Text string `json:"text"`
	Channel string `json:"channel,omitempty"`
}

func SendSlackMessage(message string, jobUrl string) error {
	//slackBody, _ := json.Marshal(requestBody{Text: message,Channel: os.Getenv("channel")})
	slackWebHookURL, err := env.GetSlackWebHookURL()
	if err != nil {
		return err
	}
	attachment1 := slack.Attachment {}
	attachment1.AddField(slack.Field { Title: "Status", Value: "Failed" })
	attachment1.AddAction(slack.Action { Type: "button", Text: "View logs", Url: jobUrl, Style: "primary" })
	payload := slack.Payload {
      Text: message,
      Username: "robot",
      Channel: os.Getenv("channel"),
      IconEmoji: ":monkey_face:",
      Attachments: []slack.Attachment{attachment1},
    }
    errr := slack.Send(slackWebHookURL, "", payload)
    if len(errr) > 0 {
      fmt.Printf("error: %s\n", errr)
    }


	return nil
}
