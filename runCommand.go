package ec2spotnotify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

/*
RunCommand executes a given command defined in EC2SPOT_RUN_COMMAND and returns an Error
It executes any powershell script via powershell.exe if underlying system is Windows
It runs any command including its arguments via 'sh -c'
OS detection relies on GOOS defined at compile time

    Example:

        EC2SPOT_RUN_COMMAND="logger 'ec2spotnotify executed successfully at $(date)"
        EC2SPOT_RUN_COMMAND="saveStateS3.sh"
        EC2SPOT_RUN_COMMAND="deregisterELB.ps1"
*/
func RunCommand() (err error) {

	script := os.Getenv("EC2SPOT_RUN_COMMAND")

	if errs := isEmpty(script); errs != nil {
		err = errEmptyCommand
		return
	}

	switch runtime.GOOS {
	case "windows":
		if _, errs := exec.Command("powershell.exe", "-File", script).Output(); errs != nil {
			err = fmt.Errorf("[!] An error occurred while executing this command: %s. Details: %s", script, errs)
		}
	default:
		if _, errs := exec.Command("sh", "-c", script).Output(); errs != nil {
			err = fmt.Errorf("[!] An error occurred while executing this command: %s. Details: %s", script, errs)
		}
	}

	return
}
