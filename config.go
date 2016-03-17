package ec2spotnotify

import (
	"os"
	"strings"
)

// Config contains AWS Region (used by SDK session), Error, SNS and URL additional structs
type Config struct {
	Region string
	SNS
	URL
	Err error
}

// SNS contains properties required to be fullfiled within AWS SNS Publish API
type SNS struct {
	Topic   string
	Subject string
	Message string
}

// URL instance termination endpoint that should return a timestamp value
type URL struct {
	InstanceTermination string
}

// Helper function to check whether a string is empty or not and it returns an erorr type to be evaluated against
func isEmpty(s string) (err error) {

	if len(strings.TrimSpace(s)) == 0 {
		err = errItsEmpty
	}

	return
}

// Parses environment variable for Notification Endpoint and return default EC2 URLs if not set
// For custom endpoints (testing purposes) fill in EC2SPOT_NOTIFICATION_ENDPOINT
func (c *Config) parseURLEndpoint() {

	c.URL.InstanceTermination = os.Getenv("EC2SPOT_NOTIFICATION_ENDPOINT")

	if err := isEmpty(c.URL.InstanceTermination); err != nil {
		c.URL.InstanceTermination = defaultURLTerminationNotification
	}

}

// LoadConfig looks up for EC2SPOT_REGION and EC2SPOT_SNS_TOPIC to ensure it can proceed without issues
// Should these environment variables be empty it exists with a non-0 status
func (c *Config) LoadConfig() {

	c.Region = os.Getenv("EC2SPOT_REGION")
	c.SNS.Topic = os.Getenv("EC2SPOT_SNS_TOPIC")
	c.SNS.Subject = os.Getenv("EC2SPOT_SNS_SUBJECT")

	if err := isEmpty(c.Region); err != nil {
		c.Err = errEmptyAWSRegion
	}

	if err := isEmpty(c.SNS.Topic); err != nil {
		c.Err = errEmptySNSTopic
	}

	if err := isEmpty(c.SNS.Subject); err != nil {
		c.SNS.Subject = "Spot Termination notification"
	}
}
