# Purpose

ec2SpotNotify is aimed to monitor EC2 Instance Metadata for EC2 Spot Termination Notices every 3 seconds. Once termination notice has been provided by EC2, you can run an arbitrary command or send a message to AWS SNS in which you can invoke a Lambda function thereafter (i.e resize AutoScaling Group for on-demand instances, send to SQS queue and count it later, etc).

**Environment variables that can be used to configure it:**

* **EC2SPOT_SNS_TOPIC**="SNS Topic ARN (i.e. arn:aws:sns:region:accountNumber:topicName)"
    * ***Required***
* **EC2SPOT_REGION**="region (i.e. eu-west-1"
    * ***Required***
* **EC2SPOT_RUN_COMMAND**="logger 'Command ran successfully........................$(date)'" 
    * *Linux:* any command or script as that can be invoked by 'sh -c' (e.g. EC2SPOT_RUN_COMMAND="saveStateS3.sh"
    * *Windows:* must be a Powershell script (e.g. EC2SPOT_RUN_COMMAND="deregisterELB.ps1"
* **EC2SPOT_NOTIFICATION_ENDPOINT**="http://dockerhost/fakeTimestamp"
    * Used for testing/dev purposes in case you want to run it locally before using an EC2 instance

Example as to how to consume this package can be found  **[Here]**(https://github.com/heitorlessa/ec2SpotNotify/blob/master/example/main.go)

### Why yet another EC2 Spot SNS small program?

Mostly for 2 reasons:

1. Wanted to get my hands dirty and cause some invalid memory by doing something meaningful
2. Rewrite [Spotterm project](https://github.com/rlmcpherson/spoterm) and add some cool things to get exposed to other libraries (AWS SNS, Command execution, structs, etc)

# Platform supported

This program should be able to run in both Linux and Windows platforms. But nothing stops anyone to build a binary for FreeBSD that can also run on EC2. 


![TODO](https://img.shields.io/badge/pending-actions-orange.svg)
* Create test for package
* Upload binary to Github releases and add Userdata & Launch-specification samples
