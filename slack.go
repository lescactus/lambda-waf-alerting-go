package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

// Represents the info we are interested in
type slackMessage struct {
	Title            string
	AlarmName        string
	AlarmDescription string
	Trigger          string
	OldState         string
	NewState         string
	NewStateReason   string
	AWSAccountID     string
	AWSRegion        string
	WebACL           string
}

// FormatTrigger will format the trigger message in a human readable format
// It returns the full string printed in Slack
func (s *slackMessage) FormatTrigger(t events.CloudWatchAlarmTrigger) string {
	return fmt.Sprintf("%s %s %s %.3f for %d period(s) of %d seconds in Namespace %s", t.Statistic, t.MetricName, t.ComparisonOperator, t.Threshold, t.EvaluationPeriods, t.Period, t.Namespace)
}

// FormatAttachment will create a Slack attachment to wrap all the alarm info we want in Slack
// Slack attachments supports formatting as per https://api.slack.com/docs/formatting
// It returns a slack.Attachment with the correct formatting ready to be sent to Slack
func (s *slackMessage) FormatAttachment() slack.Attachment {
	return slack.Attachment{
		Color:     "danger",
		Title:     s.Title,
		TitleLink: cloudWatchAlertLink,
		Pretext:   ":fire: *WAF Alert - WebACL: " + s.WebACL + "* :fire:",
		Text:      "Automatic alert\n",
		Fields: []slack.AttachmentField{
			{
				Title: "Alarm Name",
				Value: "_" + s.AlarmName + "_",
			},
			{
				Title: "Alarm Description",
				Value: "_" + s.AlarmDescription + "_",
			},
			{
				Title: "Alarm Trigger",
				Value: "_" + s.Trigger + "_",
			},
			{
				Title: "Alarm Old State",
				Value: "_" + s.OldState + "_",
				Short: true,
			},
			{
				Title: "Alarm New State",
				Value: "_" + s.NewState + "_",
				Short: true,
			},
			{
				Title: "Alarm New State Reason",
				Value: "_" + s.NewStateReason + "_",
			},
			{
				Title: "AWS Account ID",
				Value: "_" + s.AWSAccountID + "_",
				Short: true,
			},
			{
				Title: "AWS Region",
				Value: "_" + s.AWSRegion + "_",
				Short: true,
			},
		},
	}
}
