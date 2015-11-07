package main

import (
	"ec2SpotNotify"
	"fmt"
	"log"
)

func main() {

	log.Println("Looking up Instance Metadata....")
	notification, err := ec2SpotNotify.GetNotificationTime()
	if err != nil {
		fmt.Errorf("Ooops! Something went terribly wrong: ", err)
	}

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided
	for {
		select {
		case <-notification:
			log.Println("received termination -> ", notification)
		default:
			log.Println("Not yet...")
		}
	}
	log.Println("Doing something else...")
	// do something about it ;)
}
