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
	emoji := ":question:"

	switch message.NewStateValue {
	case "ALARM":
		color = "danger"
		emoji = ":bomb:"
	case "OK":
		color = "good"
		emoji = ":beer:"
	}

	return SlackMessage{
		Text: fmt.Sprintf("%s `%s`", emoji, message.AlarmName),
		Attachments: []Attachment{
			{
				Text:  message.NewStateReason,
				Color: color,
				Title: "Details",
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
