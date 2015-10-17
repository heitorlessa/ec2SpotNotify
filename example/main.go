package main

import (
	"ec2SpotNotify"
	"fmt"
)

func main() {

	fmt.Println("Looking up Instance Metadata....")
	notification, err := ec2SpotNotify.GetNotificationTime()
	if err != nil {
		fmt.Errorf("Ooops! Something terrible happened dude: ", err)
	}
	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided
	fmt.Println(<-notification)

	// do something about it ;)
}
