## Prerequisites

- kubectl is configured on local machine to the correct cluster
    - check this by running `kubectl cluster-info`
- go version 1.11 or higher (written using 1.13.4)

## Build

- standard golang build, ensure it is on your gopath and run `go build`

## Running

- Run the binary as you would normally e.g. `./app`
- On startup it will prompt you for some input:
    - path to manifest file
    - namespace to deploy to
    - delay between availability and deletion of resource
    
## How?

- Starts a kubectl proxy on your machine to connect to the cluster
- Runs various curl commands to create, check availability and delete resource
- Kills kubectl proxy

## Alternative approach
Another approach I considered was extending the kubenetes/client-go library to include
the seldonDeployment custom resource definition. However decided against this approach as the client-go library 
requires more cluster specific security in order to access a cluster. The current approach is more cluster-agnostic.

Please let me know if you have any questions.