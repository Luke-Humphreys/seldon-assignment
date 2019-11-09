package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	// get information from user
	filepath, namespace, deletionDelay, err := userInput()
	if err != nil {
		log.Fatalf("failed to parse user input: %s", err)
	}

	// start kubectl proxy on port 8080
	// sleep for 5s to allow it to become ready
	cmd, err := startProxy()
	time.Sleep(5 * time.Second)
	if err != nil {
		log.Fatalf("failed to start kubectl proxy: %s", err.Error())
	}

	// create the resource in the defined ns
	resourceName, err := createResource(filepath, namespace)
	if err != nil {
		_ = killProxy(cmd)
		log.Fatalf("failed to create resource: %s", err.Error())
	}

	// watch until resource becomes available
	err = monitorResource(resourceName, namespace)
	if err != nil {
		_ = killProxy(cmd)
		log.Fatalf("failed to monitor resource: %s", err.Error())
	}

	// delete resource after delay period
	log.Printf("Deleting resource in %d seconds", deletionDelay)
	time.Sleep(time.Duration(deletionDelay) * time.Second)
	err = deleteResource(resourceName, namespace)
	if err != nil {
		_ = killProxy(cmd)
		log.Fatalf("failed to delete resource: %s", err.Error())
	}

	// kill kubectl proxy to clean up afterwards
	err = killProxy(cmd)
	if err != nil {
		log.Fatalf("failed to kill proxy: %s", err.Error())
	}
}

// gathers input from stdin
func userInput() (string, string, int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the filepath for the manifest file: ")
	filePath, err := reader.ReadString('\n')
	if err != nil {
		return "", "", 0, fmt.Errorf("could not read user input: %s", err.Error())
	}
	// remove trailing newline char
	filePath = strings.Trim(filePath, "\n")

	fmt.Print("Enter the namespace to deploy to: ")
	namespace, err := reader.ReadString('\n')
	if err != nil {
		return "", "", 0, fmt.Errorf("could not read user input: %s", err.Error())
	}
	// remove trailing newline char
	namespace = strings.Trim(namespace, "\n")

	var deletionDelay int
	fmt.Print("Enter the time delay in seconds between resource becoming available and being deleted: ")
	_, err = fmt.Scanf("%d", &deletionDelay)
	if err != nil {
		return "", "", 0, fmt.Errorf("could not read user input: %s", err.Error())
	}
	fmt.Println(deletionDelay)

	return filePath, namespace, deletionDelay, nil
}
