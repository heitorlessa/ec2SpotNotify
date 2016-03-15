package main

import (
	"fmt"
	"github.com/heitorlessa/ec2spotnotify"
	"github.com/parnurzeal/gorequest"
	"log"
	"time"
)

const (
	defaultURLInstanceDetails = "http://169.254.169.254/latest/dynamic/instance-identity/document" // dumps ec2 metadata JSON
)

var (
	requestTimeout = 300 * time.Millisecond // how long before we give up on reaching ec2 metadata
)

func main() {

	var c *ec2spotnotify.Config
	config := c.LoadConfig()

	log.Println("Looking up Instance Metadata....")
	notification, err := ec2spotnotify.GetNotificationTime()
	if err != nil {
		log.Fatalln("Ooops! Something went terribly wrong: ", err)
	}

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided with range
	for timestamp := range notification {
		log.Println("Received notification -> ", timestamp)
		instance, err := collectInstanceDetails()
		if err != nil {
			log.Fatalln("Error found while trying to get instance metadata: ", err)
		}

		config.SNS.Message = instance
		ec2spotnotify.PublishSNS(config)
	}

	// run command and its arguments via sh -c or powershell if specified
	ec2spotnotify.RunCommand()
}

// Look up for instance metadata to gather more details about this EC2 instance
func collectInstanceDetails() (instance string, err error) {

	req := gorequest.
		New().
		Timeout(requestTimeout)

	_, instance, errs := req.Get(defaultURLInstanceDetails).End()

	if errs != nil {
		err = fmt.Errorf("Failed to retrieve instance details: %v", errs)
		return
	}

	return instance, err
}
