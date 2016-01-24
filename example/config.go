package main

import (
	"fmt"
	"log"
	"os"
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
