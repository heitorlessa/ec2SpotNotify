package main

import (
	"log"

	"github.com/heitorlessa/ec2spotnotify"
)

func main() {

	var c *ec2spotnotify.Config
	config := c.LoadConfig()

	log.Println("Looking up Instance Metadata....")
	notification, instance, err := ec2spotnotify.GetNotificationTime()
	if err != nil {
		log.Fatalln("Ooops! Something went terribly wrong: ", err)
	}

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided with range
	for timestamp := range notification {
		log.Println("Received notification -> ", timestamp)
		config.SNS.Message = instance
		ec2spotnotify.PublishSNS(config)
	}

	// run command and its arguments via sh -c or powershell if specified
	ec2spotnotify.RunCommand()
}
