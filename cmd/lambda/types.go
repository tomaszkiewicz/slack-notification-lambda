package main

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
	Color   string `json:"Color""`
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
