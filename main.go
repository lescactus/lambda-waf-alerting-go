package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/slack-go/slack"
)

const (
	envSlackChannel                  = "SLACK_CHANNEL"
	envAWSSecretsManagerName         = "AWS_SECRETS_MANAGER_NAME"
	envCloudWatchAlertLink           = "CLOUDWATCH_ALERT_LINK"
	envWebACL                        = "WEB_ACL"
	jsonAWSSecretsManagerNameKeyName = "slackToken"
)

var (
	slackToken            string // Slack token. Stored in AWS Secrets Manager
	slackChannel          string // Name of the Slack channel where to send the alert
	awsSecretsManagerName string // AWS Secrets Manager name where the Slack token is stored into
	cloudWatchAlertLink   string // URL of the CloudWatch alarm
	webACL                string // WebACL name for which the lambda is set as a trigger
)

func init() {
	if slackChannel = os.Getenv(envSlackChannel); slackChannel == "" {
		log.Fatalln("Environment variable " + envSlackChannel + " is not set but is required! Exiting!")
	}
	if awsSecretsManagerName = os.Getenv(envAWSSecretsManagerName); awsSecretsManagerName == "" {
		log.Fatalln("Environment variable " + envAWSSecretsManagerName + " is not set but is required! Exiting!")
	}
	if cloudWatchAlertLink = os.Getenv(envCloudWatchAlertLink); cloudWatchAlertLink == "" {
		log.Fatalln("Environment variable " + envCloudWatchAlertLink + " is not set but is required! Exiting!")
	}
	if webACL = os.Getenv(envWebACL); webACL == "" {
		log.Fatalln("Environment variable " + envWebACL + " is not set but is required! Exiting!")
	}
}

func handler(ctx context.Context, snsEvent events.SNSEvent) {
	// Create a new AWS SecretsManager client
	awssmClient, err := New()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Retrieve the Slack token stored in AWS Secrets Manager
	slackToken, err := awssmClient.GetSlackToken()
	if err != nil {
		log.Println(err)
		return
	}

	// New Slack client
	api := slack.New(slackToken)

	// https://github.com/aws/aws-lambda-go/blob/master/events/sns.go#L34
	var snsPayload events.CloudWatchAlarmSNSPayload

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		// Log in CloudWatch
		log.Printf("[%s %s] Message = %s \n", record.EventSource, snsRecord.Timestamp, snsRecord.Message)

		// Unmarshal the SNSEntity.Message as it is a json string representing a CloudWatchAlarmSNSPayload
		// https://github.com/aws/aws-lambda-go/blob/master/events/sns.go#L29
		if err := json.Unmarshal([]byte(snsRecord.Message), &snsPayload); err != nil {
			log.Fatalln("Unmarshaling json failed! Exiting!")
			return
		}

		slackMessage := slackMessage{
			Title:            ":rotating_light: AWS CloudWatch Notification :rotating_light:\n",
			AlarmName:        snsPayload.AlarmName,
			AlarmDescription: snsPayload.AlarmDescription,
			OldState:         snsPayload.OldStateValue,
			NewState:         snsPayload.NewStateValue,
			NewStateReason:   snsPayload.NewStateReason,
			AWSAccountID:     snsPayload.AWSAccountID,
			AWSRegion:        snsPayload.Region,
			WebACL:           webACL,
		}
		slackMessage.Trigger = slackMessage.FormatTrigger(snsPayload.Trigger)

		// Actual send of the message to Slack
		channelID, timestamp, err := api.PostMessage(
			slackChannel,
			slack.MsgOptionAttachments(slackMessage.FormatAttachment()),
			slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
		)
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
}

func main() {
	lambda.Start(handler)
}
