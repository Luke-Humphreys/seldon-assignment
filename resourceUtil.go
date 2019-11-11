package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// structs used to marshall json response from apiserver
type CreateResourceResponse struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

type GetAvailabilityResponse struct {
	Status struct {
		State string `json:"state"`
	} `json:"status"`
}

// creates resource by curling apiserver, returns the name of created resource
func createResource(filepath, namespace string) (string, error) {
	// open file
	f, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to open manifest file: %s", err.Error())
	}
	defer f.Close()

	log.Println("Creating resource")
	// construct api call
	url := fmt.Sprintf("http://127.0.0.1:8080/apis/machinelearning.seldon.io/v1alpha2/namespaces/%s/seldondeployments", namespace)
	req, err := http.NewRequest("POST", url, f)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %s", err.Error())
	}

	// get file type from filename (yaml or json)
	contentType := fmt.Sprintf("application/%s", getFiletype(filepath))
	req.Header.Set("Content-Type", contentType)

	// do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to curl mainifest to api: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return "", fmt.Errorf(resp.Status)
	}

	name, err := getResourceName(resp)
	if err != nil {
		return "", err
	}
	log.Printf("Resource '%s' created with status: %d", name, resp.StatusCode)
	return name, nil
}

// returns the name of the created resource
func getResourceName(resp *http.Response) (string, error) {
	createdResponse := CreateResourceResponse{}
	err := json.NewDecoder(resp.Body).Decode(&createdResponse)
	if err != nil {
		err = fmt.Errorf("failed to decode json response: %s", err.Error())
		return "", err
	}
	return createdResponse.Metadata.Name, nil
}

// monitors resource until it becomes available, then returns
func monitorResource(resourceName, namespace string) error {
	log.Printf("Waiting for '%s' to become available", resourceName)
	url := fmt.Sprintf("http://127.0.0.1:8080/apis/machinelearning.seldon.io/v1alpha2/namespaces/%s/seldondeployments/%s", namespace, resourceName)
	available := false
	var err error
	// infinite loop
	for {
		if available == true {
			break
		}
		// poll every second
		time.Sleep(1 * time.Second)
		available, err = getAvailable(url)
		if err != nil {
			return err
		}
	}
	log.Printf("Resource '%s' now available", resourceName)
	return nil
}

// curl -X GET on the resource
func getAvailable(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to check state of resource: %s", err.Error())
	}
	defer resp.Body.Close()
	available, err := decodeAvailability(resp)
	if err != nil {
		return false, fmt.Errorf("failed to get state of resouce: %s", err.Error())
	}
	return available, nil
}

// decode getAvailable response and returns bool of whether resource is available
func decodeAvailability(resp *http.Response) (bool, error) {
	availabilityResponse := GetAvailabilityResponse{}
	err := json.NewDecoder(resp.Body).Decode(&availabilityResponse)
	if err != nil {
		err = fmt.Errorf("failed to decode json response: %s", err.Error())
		return false, err
	}
	return availabilityResponse.Status.State == "Available", nil
}

// delete given resource
func deleteResource(resourceName, namespace string) error {
	url := fmt.Sprintf("http://127.0.0.1:8080/apis/machinelearning.seldon.io/v1alpha2/namespaces/%s/seldondeployments/%s", namespace, resourceName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %s", err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("faield to delete resource: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}
	log.Printf("Resource '%s' successfully deleted", resourceName)
	return nil
}

// gets filetype
func getFiletype(filePath string) string {
	x := strings.Split(filePath, ".")
	return x[1]
}
