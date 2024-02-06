package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"os"

	"github.com/achilles-git/k8s-job-notify/env"
)

type requestBody struct {
	Text string `json:"text"`
	Channel string `json:"channel,omitempty"`
}

func SendSlackMessage(message string) error {
	slackBody, _ := json.Marshal(requestBody{Text: message,Channel: os.Getenv("channel")})
	slackWebHookURL, err := env.GetSlackWebHookURL()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, slackWebHookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	if buf.String() != "ok" {
		return errors.New("non ok response from Slack")
	}
	return nil
}
