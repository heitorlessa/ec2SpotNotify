package ec2spotnotify

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Config contains AWS Region, Error interface, SNS config and Instance endpoints
type Config struct {
	Region string
	SNS
    URL
	Err error
}

// SNS properties required to be fullfiled within publish API
type SNS struct {
	Topic   string
	Subject string
	Message string
}

// URL defines Endpoints for both Instance metadata and termination notification timestamp
type URL struct {
    InstanceDetails string
    InstanceTermination string
}

// Helper function to check whether a string is empty or not and it returns an erorr type to be evaluated against
func isEmpty(s string) (err error) {

	if len(strings.TrimSpace(s)) == 0 {
		err = fmt.Errorf("Value is empty")
	}

	return
}

// LoadConfig looks up for EC2SPOT_REGION and EC2SPOT_SNS_TOPIC to ensure it can proceed without issues
// Should these environment variables be empty it exists with a non-0 status
func LoadConfig() *Config {

	config := &Config{
		SNS: SNS{
			Subject: "Spot Termination notification",
		},
	}

	config.Region = os.Getenv("EC2SPOT_REGION")
	config.SNS.Topic = os.Getenv("EC2SPOT_SNS_TOPIC")
    
    // optional as GetNotificationTime will use default if not informed
	config.URL.InstanceDetails = os.Getenv("EC2SPOT_METADATA")
	config.URL.InstanceTermination = os.Getenv("EC2SPOT_NOTIFICATION")

	if err := isEmpty(config.Region); err != nil {
		config.Err = fmt.Errorf("Region cannot be empty")
	}

	if err := isEmpty(config.SNS.Topic); err != nil {
		config.Err = fmt.Errorf("SNS Topic cannot be empty")
	}

	if config.Err != nil {
		log.Fatalf("Cowardly quitting due to: %s", config.Err)
	} 

	return config
}
