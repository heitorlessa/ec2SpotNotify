package ec2spotnotify

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

// publish message stored in Message under SNS struct to a SNS topic defined in EC2SPOT_SNS_TOPIC
func PublishSNS(c *Config) {

	// AWS initialization
	client := sns.New(session.New(), &aws.Config{Region: aws.String(c.Region)})

	// SNS params prior to publish message to topic
	// ref: http://docs.aws.amazon.com/sdk-for-go/api/service/sns/SNS.html#Publish-instance_method
	params := &sns.PublishInput{
		Message: aws.String(string(c.SNS.Message)),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Key": {
				DataType:    aws.String("String"),
				StringValue: aws.String("String"),
			},
		},
		Subject:  aws.String(c.SNS.Subject),
		TopicArn: aws.String(c.SNS.Topic),
	}

	_, err := client.Publish(params)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
