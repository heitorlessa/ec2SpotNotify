package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/heitorlessa/ec2spotnotify"
)

const (
	instanceDetailsURL = "http://169.254.169.254/latest/dynamic/instance-identity/document" // dumps ec2 metadata JSON
)

var (
	requestTimeout = 300 * time.Millisecond // how long before we give up on reaching ec2 metadata
)

func main() {

	log.Println("[*] Loading configuration...")
	var config ec2spotnotify.Config
	config.LoadConfig()

	if config.Err != nil {
		log.Fatalf("[*] Cowardly quitting due to: %s", config.Err)
	}

	log.Println("[*] Looking up Instance Metadata....")
	notification, err := ec2spotnotify.GetNotificationTime()

	if err != nil {
		log.Fatalf("[!] Error found while trying to query Instance Metadata: %s ", err)
	}

	// as notification may take a while to be injected on EC2 Metadata - Read from channel provided with range
	for timestamp := range notification {
		log.Println("[*] Received notification -> ", timestamp)
		instance, err := collectInstanceDetails()
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("[*] Publishing termination notification to SNS")
		config.SNS.Message = instance
		if err := ec2spotnotify.PublishSNS(config); err != nil {
			log.Fatalln(err)
		}
	}

	// run command and its arguments via sh -c or powershell if specified
	if err := ec2spotnotify.RunCommand(); err != nil {
		log.Fatalln(err)
	}
}

// Look up for instance metadata to gather more details about this EC2 instance
func collectInstanceDetails() (instance string, err error) {

	req := http.Client{Timeout: requestTimeout}

	resp, errs := req.Get(instanceDetailsURL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if errs != nil {
		err = fmt.Errorf("[!] An error occurred while retrieving instance details: %v", errs)
		return "", err
	}

	instanceDump, errs := ioutil.ReadAll(resp.Body)
	if errs != nil {
		err = fmt.Errorf("[!] An error occurred while reading response: %v", errs)
		return
	}

	return string(instanceDump), err
}
