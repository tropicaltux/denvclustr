package main

import (
	"fmt"
	"log"

	"github.com/tropicaltux/denvclustr/pkg/schema"
)

func main() {
	// Example JSON configuration for a simple cluster
	jsonData := `{
		"name": "example-cluster",
		"infrastructure": [
			{
				"id": "infra1",
				"kind": "vm",
				"provider": "aws",
				"region": "us-west-2"
			}
		],
		"nodes": [
			{
				"id": "node1",
				"infrastructure_id": "infra1",
				"properties": {
					"instance_type": "t2.micro"
				},
				"remote_access": {
					"public_ssh_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC3ExamplePublicKey"
				}
			}
		],
		"devcontainers": [
			{
				"id": "dev1",
				"node_id": "node1",
				"source": {
					"url": "https://github.com/example/repo.git"
				}
			}
		]
	}`

	// Parse the JSON data
	config, err := schema.Parse([]byte(jsonData))
	if err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}

	// Successfully parsed - print basic information about the parsed configuration
	fmt.Printf("Successfully parsed cluster: %s\n", config.Name)
	fmt.Printf("Infrastructure entries: %d\n", len(config.Infrastructure))
	fmt.Printf("Node entries: %d\n", len(config.Nodes))
	fmt.Printf("Devcontainer entries: %d\n", len(config.Devcontainers))

	// Access some specific fields from the configuration
	fmt.Printf("\nDetails:\n")
	fmt.Printf("- Infrastructure provider: %s\n", config.Infrastructure[0].Provider)
	fmt.Printf("- Node instance type: %s\n", config.Nodes[0].Properties.InstanceType)
	fmt.Printf("- Devcontainer source URL: %s\n", config.Devcontainers[0].Source.URL)
}
