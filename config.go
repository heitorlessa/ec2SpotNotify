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
	InstanceMetadata    string
	InstanceTermination string
}

// Helper function to check whether a string is empty or not and it returns an erorr type to be evaluated against
func isEmpty(s string) (err error) {

	if len(strings.TrimSpace(s)) == 0 {
		err = fmt.Errorf("Value is empty")
	}

	return
}

// Parses environment variables for both Metadata and Notification Endpoints and return default EC2 URLs if not set
func (c *Config) parseURLEndpoints() {

	c.URL.InstanceMetadata = os.Getenv("EC2SPOT_METADATA")
	c.URL.InstanceTermination = os.Getenv("EC2SPOT_NOTIFICATION")

	// set default EC2 Metadata Endpoint URL if env is not set
	if err := isEmpty(c.URL.InstanceMetadata); err != nil {
		c.URL.InstanceMetadata = defaultUrlInstanceDetails
	}

	if err := isEmpty(c.URL.InstanceTermination); err != nil {
		c.URL.InstanceTermination = defaultUrlTerminationNotification
	}

}

// LoadConfig looks up for EC2SPOT_REGION and EC2SPOT_SNS_TOPIC to ensure it can proceed without issues
// Should these environment variables be empty it exists with a non-0 status
func (c *Config) LoadConfig() *Config {

	config := &Config{
		SNS: SNS{
			Subject: "Spot Termination notification",
		},
	}

	config.Region = os.Getenv("EC2SPOT_REGION")
	config.SNS.Topic = os.Getenv("EC2SPOT_SNS_TOPIC")

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
