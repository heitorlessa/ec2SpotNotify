package ec2spotnotify

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

// executes a given command defined in EC2SPOT_RUN_SCRIPT
// Optional step so no need to include in Config
func RunCommand() {

	script := os.Getenv("EC2SPOT_RUN_COMMAND")

	if err := isEmpty(script); err != nil {
		return
	}

	// uses powershell to execute a .ps1 script instead of CMD if it's Windows otherwise fall back to old and good sh -c that accepts both scripts and commands + arguments
	if runtime.GOOS == "windows" {

		if out, err := exec.Command("powershell.exe", "-File", script).Output(); err != nil {
			fmt.Errorf("Error: %s", err)
		} else {
			log.Printf("Command result: %s", out)
		}
	}

	if out, err := exec.Command("sh", "-c", script).Output(); err != nil {
		fmt.Errorf("Error: ", err)
	} else {
		log.Printf("Command result: %s", out)
	}
}
