package main

import (
	"ec2SpotNotify"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

type Config struct {
	Region string
	SNS
}

type SNS struct {
	Topic   string
	Subject string
	Message string
}

func main() {

	log.Println("Looking up Instance Metadata....")
	notification, instance, err := ec2SpotNotify.GetNotificationTime()
	if err != nil {
		fmt.Errorf("Ooops! Something went terribly wrong: ", err)
	}

	config := loadConfig(&Config{
		Region: "<AWS Region>",
		SNS: SNS{
			Topic:   "arn:aws:sns:<region>:<accountID>:<topic>",
			Subject: "Spot Termination notification",
			Message: instance,
		},
	})

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided with range
	for timestamp := range notification {
		log.Println("Received notification -> ", timestamp)
		publishSNS(config)
	}

	log.Println("Doing something else...")
	// do something about it ;)
}

// Need to implement config check + read from env if file is not present or if 'nil' passed as argument
func loadConfig(c *Config) *Config {

	log.Println("Gathering configuration....")

	return &Config{
		Region: c.Region,
		SNS: SNS{
			Topic:   c.SNS.Topic,
			Subject: c.SNS.Subject,
			Message: c.SNS.Message,
		},
	}
}

func publishSNS(c *Config) {

	// AWS initialization
	config := aws.NewConfig().WithRegion(c.Region)
	client := sns.New(config)

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
