# Purpose

**Work-In-Progress -- Missing Test**

ec2SpotNotify is aimed to monitor EC2 Instance Metadata for EC2 Spot Termination Notices every 3 seconds. Once termination notice has been provided by EC2, you can run an arbitrary command or send a message to AWS SNS in which you can invoke a Lambda function thereafter (i.e resize AutoScaling Group for on-demand instances, send to SQS queue and count it later, etc).

**Env variables to be configured:**

* **EC2SPOT_SNS_TOPIC**="<SNS Topic ARN (i.e. arn:aws:sns:<region>:<accountNumber>:<topicName>)>"
    * ***Required***
* **EC2SPOT_REGION**="<region (i.e. eu-west-1>"
    * ***Required***
* EC2SPOT_RUN_COMMAND="date 
    * *Linux:* any command as a wrapper 'sh -c' is called
    * *Windows:* must be a Powershell script though

Simple example that would use this package and look up for these env variables - **[Here]**(https://github.com/heitorlessa/ec2SpotNotify/blob/master/example/main.go)

## Why yet another EC2 Spot SNS small program?

I'm actually using it for the sake of learning Go more appropriately, so nothing better than having small projects as I have bigger ones to come :)

# Platform supported

This program should be able to run in both Linux and Windows platforms. 

# Deployment

[!] TODO

# TODO
* Create test for package

## Improvements
* Make comments 'godoc' compatible so it can generate HTML if needed
* Read TimeThresholdTime from Config file instead of const for more flexibility
* Generate builds for Linux and Windows
* Write quick start guide for both Linux and Windows including IAM SNS Permissions for EC2 IAM Role
