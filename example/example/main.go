package main

import (
	"fmt"
	"github.com/heitorlessa/ec2spotnotify"
	"log"
)

func main() {

	config := ec2spotnotify.LoadConfig()

	log.Println("Looking up Instance Metadata....")
	notification, instance, err := ec2spotnotify.GetNotificationTime()
	if err != nil {
		fmt.Errorf("Ooops! Something went terribly wrong: ", err)
	} else {
		config.SNS.Message = instance
	}

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided with range
	for timestamp := range notification {
		log.Println("Received notification -> ", timestamp)
		ec2spotnotify.PublishSNS(config)
	}

	// run command and its arguments via sh -c or powershell if specified
	ec2spotnotify.RunCommand()
}
