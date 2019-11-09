package main

import (
	"log"
	"os/exec"
)

// start kubectl proxy on port 8080
func startProxy() (*exec.Cmd, error) {
	cmd := exec.Command("kubectl", "proxy", "--port=8080")
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	log.Println("Kubectl proxy started on port 8080")
	return cmd, nil
}

// stop running proxy
func killProxy(cmd *exec.Cmd) error {
	if err := cmd.Process.Kill(); err != nil {
		return err
	}
	log.Println("Killed kubectl proxy")
	return nil
}
