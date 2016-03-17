/*
Package ec2spotnotify implements an interface to access Amazon EC2 Spot Termination Notification via EC2 Metadata
With that information handy, it also publishes a message to an Amazon SNS Topic
In addition to that, one may want to run a set of commands after such instance has been marked for termination
Therefore that is also supported by the package but it is entirely up to the client to opt in for that

This package is broken down into the following files:

    SNS.go          - Logic to publish message to a given AWS SNS Topic
    config.go       - Lookup for environment variables that configures this package and provide default values for a few adjustable ones
    runCommand.go   - Implements exec.Command wrapper for OSes supported
    example/main.go - Client that demonstrates how to use this package altogether

The following are the core functions of its functionality:

    GetNotificationTime     - Logic to poll EC2 Instance Termination Time every 3 seconds and exposes a channel that client can look at when timestamp is ready
    lookupInstanceMetadata  - Lookup Endpoint provided in 'URL.InstanceTermination'
*/
package ec2spotnotify

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	requestNotFoundError                = "404 Not Found"
	requestFound                        = "200 OK"
	timeFormat                          = "2006-01-02T15:04:05Z07:00" // RFC 3339
	timeThresholdInterval time.Duration = 3 * time.Second
)

var (
	defaultURLTerminationNotification = "http://169.254.169.254/latest/meta-data/spot/termination-time"
	requestTimeoutThreshold           = 150 * time.Millisecond
	errEmptySNSTopic                  = errors.New("[!] SNS Topic is not defined. Please ensure 'EC2SPOT_SNS_TOPIC' env is not empty")
	errEmptyAWSRegion                 = errors.New("[!] Region is not defined. Please ensure 'EC2SPOT_REGION' env is not empty")
	errItsEmpty                       = errors.New("[!] Value is empty")
	errEmptyCommand                   = errors.New("[!] Command is not defined. Please ensure 'EC2SPOT_RUN_COMMAND' env is not empty")
)

// GetNotificationTime returns a time based channel and Error
// Returned channel should be read to identify when Spot Instance will be terminated so actions can be taken
//
//    Example:
//
// notification, err := ec2spotnotify.GetNotificationTime()
//
// if err != nil {
//     log.Fatalf("[!] Cannot continue due to: %s", err)
// }
//  for timestamp := range notification {
// 	    log.Printf("[*] Notification received at: %s ", timestamp)
//  }
func GetNotificationTime() (timestamp chan time.Time, err error) {
	notifyChan := make(chan time.Time, 1)

	// quick error check (timeout error if not an EC2 instance)
	_, errs := lookupInstanceMetadata()
	if errs != nil {
		err = errs
		return nil, err
	}

	// run goroutine and keep it running until data is available
	// Ticker from time package ensures it runs for X time until we stop it
	// once done with processing, stops ticker and closes notify channel to stop receiving messages
	go func() {
		ticker := time.NewTicker(timeThresholdInterval)
		defer ticker.Stop()
		defer close(notifyChan)

		for _ = range ticker.C {
			notification, errs := lookupInstanceMetadata()

			if errs != nil {
				return
			}

			if !notification.IsZero() {
				notifyChan <- notification
				return
			}
		}
	}()

	return notifyChan, nil
}

// lookupInstanceMetadata looks up at EC2 Instance Metadata and returns termination notification and Error
func lookupInstanceMetadata() (timestamp time.Time, err error) {

	// return a Zero timestamp if termination notification is not set
	// ref: https://golang.org/pkg/time/#Time.IsZero
	ZeroTimestamp, _ := time.Parse(timeFormat, "")

	var config Config
	config.parseURLEndpoint()

	// set shorter timeout as default is way too high for this operation
	req := http.Client{Timeout: requestTimeoutThreshold}
	resp, errs := req.Get(config.URL.InstanceTermination)

	if resp != nil {
		defer resp.Body.Close()
	}

	// return error if request times out
	if errs != nil {
		err = fmt.Errorf("[!] Is this running on an EC2 Instance? Details: %v", errs)
		return
	}

	notification, errs := ioutil.ReadAll(resp.Body)
	if errs != nil {
		err = fmt.Errorf("[!] An error occurred while reading response: %v", errs)
		return
	}

	// While instance is not marked for termination EC2 should keep returning HTTP 404
	switch resp.Status {
	case requestNotFoundError:
		return ZeroTimestamp, nil
	case requestFound:
		timestamp, _ = time.Parse(timeFormat, string(notification))
		return timestamp, nil
	default:
		return ZeroTimestamp, nil
	}
}
