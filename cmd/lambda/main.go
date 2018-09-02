package main

// based on: http://www.blog.labouardy.com/slack-notification-cloudwatch-alarms-lambda/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Records []struct {
		SNS struct {
			Type      string `json:"Type"`
			Timestamp string `json:"Timestamp"`
			Message   string `json:"Message"`
		} `json:"Sns"`
	} `json:"Records"`

	// For simplified alerts
	Message string `json:"Message"`
	Color   string `json:Color`
}

type SNSMessage struct {
	// Cloudwatch Alerts
	AlarmName      string `json:"AlarmName"`
	NewStateValue  string `json:"NewStateValue"`
	NewStateReason string `json:"NewStateReason"`

	// Bounce and complaints for SES
	NotificationType string `json:"NotificationType"`
}

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text  string `json:"text"`
	Color string `json:"color"`
	Title string `json:"title"`
}

func handler(request Request) error {
	if request.Message != "" {
		// handle simple message

		// TODO

		return nil
	}

	if len(request.Records) > 0 {
		// handle messages from SNS

		var snsMessage SNSMessage

		err := json.Unmarshal([]byte(request.Records[0].SNS.Message), &snsMessage)

		if err != nil {
			return err
		}

		if snsMessage.NotificationType == "Bounce" || snsMessage.NotificationType == "Complaint" {
			// handle message from SES

			// TODO

			return nil
		}

		log.Printf("New alarm: %s - Reason: %s", snsMessage.AlarmName, snsMessage.NewStateReason)

		slackMessage := buildCloudwatchAlertSlackMessage(snsMessage)
		postToSlack(slackMessage)

		log.Println("Notification has been sent")
	}

	return nil
}

func buildCloudwatchAlertSlackMessage(message SNSMessage) SlackMessage {
	color := "gray"

	switch message.NewStateValue {
	case "ALARM":
		color = "danger"
	case "OK":
		color = "good"
	}

	return SlackMessage{
		Text: fmt.Sprintf("`%s`", message.AlarmName),
		Attachments: []Attachment{
			{
				Text:  message.NewStateReason,
				Color: color,
				Title: "Reason",
			},
		},
	}
}

func postToSlack(message SlackMessage) error {
	client := &http.Client{}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("SLACK_WEBHOOK"), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
