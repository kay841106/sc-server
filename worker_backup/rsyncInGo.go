package main

import (
	"fmt"
	"log"
	"os/exec"
)

// /home/ubuntu/bimo
func rsyncCommand() {
	cmd := exec.Command("python")
	// cmd := exec.Command("/usr/bin/rsync", "-av", "/home/ubuntu/bimo", "/opt/nas/meter")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	a, _ := cmd.Output()
	fmt.Println(a)
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}
func main() {
	// c := cron.New()
	// c.AddFunc("@daily", rsyncCommand)

	rsyncCommand()
}
