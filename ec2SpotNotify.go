package ec2spotnotify

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	requestNotFoundError                = "404 Not Found"
	requestFound                        = "200 OK"
	timeFormat                          = "2006-01-02T15:04:05Z07:00" // RFC 3339
	timeThresholdInterval time.Duration = 3 * time.Second
)

var (
	defaultUrlTerminationNotification               = "http://169.254.169.254/latest/meta-data/spot/termination-time"
	defaultUrlInstanceDetails                       = "http://169.254.169.254/latest/dynamic/instance-identity/document"
	RequestTimeoutSeconds             time.Duration = 150 * time.Millisecond
)

// GetNotificationTime is a public function that returns a time based channel and error
// Returned channel should be read to identify when Spot Instance will be terminated so actions can be taken
func GetNotificationTime() (chan time.Time, string, error) {
	notifyChan := make(chan time.Time, 1)

	// quick error check (timeout error if not an EC2 instance)
	_, instanceDetails, err := lookupInstanceMetadata()
	if err != nil {
		log.Fatal(err)
	}

	// run goroutine and keep it running until data is available
	// Ticker from time package ensures it runs for X time until we stop it
	// once done with processing, stop ticker and close channel to stop receiving messages
	go func() {
		ticker := time.NewTicker(timeThresholdInterval)
		defer ticker.Stop()
		defer close(notifyChan)
		// listen to "ticks" and do something about it
		for _ = range ticker.C {
			// find out if data is available and store it on notification
			notification, _, err := lookupInstanceMetadata()

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			// data is ready! Send it over to notifyChan channel and clean up everything
			if !notification.IsZero() {
				notifyChan <- notification
				return
			}
		}
	}()

	return notifyChan, instanceDetails, nil
}

// lookupInstanceMetadata looks at EC2 Instance Metadata for URLTerminationNotification to extract when Spot Instance will be terminated
// Returns timestamp and error so it can be worked by GetNotificationTime
// While the instance is not marked for termination EC2 will return HTTP 404
func lookupInstanceMetadata() (timestamp time.Time, instance string, err error) {

	// return a Zero timestamp if termination notification is not ready
	// ref: ref: https://golang.org/pkg/time/#Time.IsZero
	ZeroTimestamp, _ := time.Parse(timeFormat, "")

	var config Config
	config.parseURLEndpoints()

	// set shorter timeout as default is way too high for this operation
	req := gorequest.
		New().
		Timeout(RequestTimeoutSeconds)
	resp, body, errs := req.Get(config.URL.InstanceTermination).End()

	// return error if request times out
	if errs != nil {
		err = fmt.Errorf("[!] Is this running on an EC2 Instance? Error found: %v", errs)
		return
	}

	// make another request only for instance details (eg AZ, IP, AMI, etc)
	req2 := gorequest.
		New().
		Timeout(RequestTimeoutSeconds)

	_, instanceDetails, _ := req2.Get(config.URL.InstanceMetadata).End()

	switch resp.Status {
	case requestNotFoundError:
		return ZeroTimestamp, "", nil
		fallthrough
	case requestFound:
		notification, _ := time.Parse(timeFormat, body)
		return notification, instanceDetails, nil
		fallthrough
	default:
		fmt.Errorf("[!] Received a non-compliant status: %s ", resp.Status)
		return ZeroTimestamp, "", nil
	}
	return
}
