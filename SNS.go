package ec2spotnotify

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

/*
PublishSNS publishes message to a SNS topic (EC2SPOT_SNS_TOPIC) and returns Error
Refer to SNS struct for required properties
*/
func PublishSNS(c Config) (err error) {

	// AWS initialization
	client := sns.New(session.New(), &aws.Config{Region: aws.String(c.Region)})

	/*
		    SNS params prior to publish message to topic
			ref: http://docs.aws.amazon.com/sdk-for-go/api/service/sns/SNS.html#Publish-instance_method
	*/
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

	_, errs := client.Publish(params)
	if errs != nil {
		err = fmt.Errorf("[!] Error found while trying to publish instance details via SNS: %s", errs)
	}

	return
}
