package ec2SpotNotify

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
)

const (
	RequestNotFoundError       string        = "404 Not Found"
	RequestTimeoutError        string        = "i/o timeout"
	RequestFound               string        = "200 OK"
	RequestTimeoutSeconds      time.Duration = 150 * time.Millisecond
	URLTerminationNotification string        = "http://169.254.169.254/latest/meta-data/spot/termination-time"
	TimeFormat                 string        = "2015-01-05T18:02:00Z"
	TimeThresholdInterval      time.Duration = 3 * time.Second
)

// GetNotificationTime is a public function that returns a time based channel and error
// Returned channel should be read to identify when Spot Instance will be terminated so actions can be taken
func GetNotificationTime() (chan time.Time, error) {
	notifyChan := make(chan time.Time, 1)

	// quick error check (timeout error if not an EC2 instance)
	if _, err := lookupInstanceMetadata(); err != nil {
		return nil, err
	}

	// run goroutine and keep it running until data is available
	// Ticker from time package ensures it runs for X time until we stop it
	// once done with processing, stop ticker and close channel to stop receiving messages
	go func() {
		ticker := time.NewTicker(TimeThresholdInterval)
		defer ticker.Stop()
		defer close(notifyChan)
		// listen to "ticks" and do something about it
		for _ = range ticker.C {
			// find out if data is available and store it on notification
			notification, err := lookupInstanceMetadata()

			if err != nil {
				return
			}

			// data is ready! Send it over to notifyChan channel and clean up everything
			if !notification.IsZero() {
				notifyChan <- notification
				return
			}
		}
	}()

	return notifyChan, nil
}

// lookupInstanceMetadata looks at EC2 Instance Metadata for URLTerminationNotification to extract when Spot Instance will be terminated
// Returns timestamp and error so it can be worked by GetNotificationTime
// While the instance is not marked for termination EC2 will return HTTP 404
func lookupInstanceMetadata() (timestamp time.Time, err error) {

	// set shorter timeout as default is way too high for this operation
	request := gorequest.
		New().
		Timeout(RequestTimeoutSeconds)
	resp, body, errs := request.Get(URLTerminationNotification).End()

	// return error if request times out
	if errs != nil {
		err = fmt.Errorf("[!] Is this running on an EC2 Instance? Error found: %v", errs)
		return
	}

	switch resp.Status {
	case RequestNotFoundError:
		// return a Zero timestamp if notification is not ready
		// Zero timestamp ref: https://golang.org/pkg/time/#Time.IsZero
		notification, _ := time.Parse(TimeFormat, string(""))
		return notification, nil
	case RequestFound:
		fmt.Println("[+] Found it!!")
		notification, _ := time.Parse(TimeFormat, string(body))
		return notification, nil
	default:
		fmt.Errorf("[!] Received a non-compliant status: %s ", resp.Status)
		return
	}
	return
}
