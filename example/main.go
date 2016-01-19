package main

import (
	"ec2SpotNotify"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Config struct {
	Region string
	SNS
	Err error
}

type SNS struct {
	Topic   string
	Subject string
	Message string
}

func main() {

	config := loadConfig()

	log.Println("Looking up Instance Metadata....")
	notification, instance, err := ec2SpotNotify.GetNotificationTime()
	if err != nil {
		fmt.Errorf("Ooops! Something went terribly wrong: ", err)
	} else {
		config.SNS.Message = instance
	}

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided with range
	for timestamp := range notification {
		log.Println("Received notification -> ", timestamp)
		publishSNS(config)
	}

	// run command and its arguments via sh -c or powershell if specified
	runCommand()
}

// Helper function to check whether a string is empty or not and it returns an erorr type to be evaluated against
func isEmpty(s string) (err error) {

	if len(strings.TrimSpace(s)) == 0 {
		err = fmt.Errorf("Value is empty")
	}

	return
}

// loadConfig looks up for EC2SPOT_REGION and EC2SPOT_SNS_TOPIC to ensure it can proceed without issues
// Should these environment variables be empty it exists with a non-0 status
func loadConfig() *Config {

	config := &Config{
		Region: "",
		SNS: SNS{
			Topic:   "",
			Subject: "Spot Termination notification",
			Message: "",
		},
		Err: nil,
	}

	log.Println("Gathering configuration....")
	config.Region = os.Getenv("EC2SPOT_REGION")
	config.SNS.Topic = os.Getenv("EC2SPOT_SNS_TOPIC")

	if err := isEmpty(config.Region); err != nil {
		config.Err = fmt.Errorf("Region cannot be empty")
		log.Println(config.Err)
	}

	if err := isEmpty(config.SNS.Topic); err != nil {
		config.Err = fmt.Errorf("SNS Topic cannot be empty")
		log.Println(config.Err)
	}

	if config.Err != nil {
		log.Fatalf("Cowardly quitting due to: %s", config.Err)
	} else {
		return &Config{
			Region: config.Region,
			SNS: SNS{
				Topic:   config.SNS.Topic,
				Subject: config.SNS.Subject,
				Message: config.SNS.Message,
			},
			Err: config.Err,
		}
	}

	return config
}

// publish message stored in Message under SNS struct to a SNS topic defined in EC2SPOT_SNS_TOPIC
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

// executes a given command defined in EC2SPOT_RUN_SCRIPT
// Optional step so no need to include in Config
func runCommand() {

	script := os.Getenv("EC2SPOT_RUN_COMMAND")

	if err := isEmpty(script); err != nil {
		return
	}

	// uses powershell to execute a .ps1 script instead of CMD if it's Windows otherwise fall back to old and good sh -c that accepts both scripts and commands + arguments
	if runtime.GOOS == "windows" {

		if out, err := exec.Command("powershell.exe", "-File", script).Output(); err != nil {
			fmt.Errorf("Error: %s", err)
		} else {
			log.Printf("Command result: %s", out)
		}
	}

	if out, err := exec.Command("sh", "-c", script).Output(); err != nil {
		fmt.Errorf("Error: ", err)
	} else {
		log.Printf("Command result: %s", out)
	}
}
