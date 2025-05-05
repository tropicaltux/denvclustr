package denvclustr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// collectInfrastructureIds returns a map of all infrastructure IDs
func collectInfrastructureIds(infrastructure []*InfrastructureProvider) map[string]struct{} {
	infrastructureIds := make(map[string]struct{})
	for _, infra := range infrastructure {
		infrastructureIds[infra.Id] = struct{}{}
	}
	return infrastructureIds
}

// collectNodes returns a map of nodes by ID
func collectNodes(nodes []*Node) map[string]*Node {
	nodeMap := make(map[string]*Node)
	for _, node := range nodes {
		nodeMap[node.Id] = node
	}
	return nodeMap
}

func processInfrastructure(infrastructure []*InfrastructureProvider) error {
	// Track Ids to ensure uniqueness using a set
	seenIds := make(map[string]struct{})

	type infrastructureTuple struct {
		Kind     InfrastructureKind
		Provider Provider
		Region   string
	}

	// Track (Kind, Provider, Region) tuples to ensure uniqueness
	seenTuples := make(map[infrastructureTuple]string)

	for i, infra := range infrastructure {
		// Trim string fields
		infra.Id = strings.TrimSpace(infra.Id)
		infra.Region = strings.TrimSpace(infra.Region)

		// Validate Id is not empty
		if infra.Id == "" {
			return fmt.Errorf("infrastructure provider at index %d has empty Id; Id must be specified", i)
		}

		// Check if region is empty
		if infra.Region == "" {
			return fmt.Errorf("infrastructure provider at index %d has empty region; region must be specified", i)
		}

		// Check for duplicate Ids
		if _, exists := seenIds[infra.Id]; exists {
			return fmt.Errorf("duplicate infrastructure provider Id '%s'; Ids must be unique", infra.Id)
		}
		seenIds[infra.Id] = struct{}{}

		// Create a tuple for (Kind, Provider, Region)
		tuple := infrastructureTuple{
			Kind:     infra.Kind,
			Provider: infra.Provider,
			Region:   infra.Region,
		}
		// Check for duplicate tuples
		if existingId, exists := seenTuples[tuple]; exists {
			return fmt.Errorf("infrastructure provider '%s' has the same combination of Kind=%s, Provider=%s, Region=%s as provider '%s'; these values must be unique as a group",
				infra.Id, infra.Kind, infra.Provider, infra.Region, existingId)
		}
		seenTuples[tuple] = infra.Id
	}

	return nil
}

func processNodes(nodes []*Node, infrastructureIds map[string]struct{}) error {
	// Track node IDs to ensure uniqueness
	seenNodeIds := make(map[string]struct{})

	for i, node := range nodes {
		// Trim string fields
		node.Id = strings.TrimSpace(node.Id)
		node.InfrastructureId = strings.TrimSpace(node.InfrastructureId)
		node.RemoteAccess.PublicSSHKey = strings.TrimSpace(node.RemoteAccess.PublicSSHKey)
		node.Properties.InstanceType = strings.TrimSpace(node.Properties.InstanceType)

		if node.DNS != nil {
			node.DNS.HighLevelDomain = strings.TrimSpace(node.DNS.HighLevelDomain)
		}

		// Validate Id is not empty
		if node.Id == "" {
			return fmt.Errorf("node at index %d has empty Id; Id must be specified", i)
		}

		// Check for duplicate node Ids
		if _, exists := seenNodeIds[node.Id]; exists {
			return fmt.Errorf("duplicate node Id '%s'; node Ids must be unique", node.Id)
		}
		seenNodeIds[node.Id] = struct{}{}

		// Validate InfrastructureId references a valid infrastructure provider
		if node.InfrastructureId == "" {
			return fmt.Errorf("node '%s' has empty infrastructure_id; infrastructure_id must be specified", node.Id)
		}

		if _, exists := infrastructureIds[node.InfrastructureId]; !exists {
			return fmt.Errorf("node '%s' references non-existent infrastructure provider '%s'", node.Id, node.InfrastructureId)
		}

		// Validate RemoteAccess.PublicSSHKey is not empty
		if node.RemoteAccess.PublicSSHKey == "" {
			return fmt.Errorf("node '%s' has empty public_ssh_key in remote_access; public_ssh_key must be specified", node.Id)
		}

		// Validate Properties.InstanceType is not empty
		if node.Properties.InstanceType == "" {
			return fmt.Errorf("node '%s' has empty instance_type in properties; instance_type must be specified", node.Id)
		}

		// Validate DNS if present
		if node.DNS != nil && node.DNS.HighLevelDomain == "" {
			return fmt.Errorf("node '%s' has empty high_level_domain in DNS configuration", node.Id)
		}
	}

	return nil
}

// processSource validates the Source configuration of a devcontainer
func processSource(dc *Devcontainer) error {
	// If Source is nil, nothing to validate
	if dc.Source == nil {
		return fmt.Errorf("devcontainer '%s' has nil source; source is required", dc.Id)
	}

	// Trim string fields
	dc.Source.URL = strings.TrimSpace(dc.Source.URL)
	dc.Source.Branch = strings.TrimSpace(dc.Source.Branch)
	dc.Source.DevcontainerPath = strings.TrimSpace(dc.Source.DevcontainerPath)

	// URL is required
	if dc.Source.URL == "" {
		return fmt.Errorf("devcontainer '%s' has empty source URL; a valid Git repository URL is required", dc.Id)
	}

	// Check if URL is SSH or HTTPS
	isSSHUrl := strings.HasPrefix(dc.Source.URL, "ssh://") || strings.HasPrefix(dc.Source.URL, "git@")

	// Validate SSH key requirements based on URL type
	if isSSHUrl {
		// SSH URLs require SSH key
		if dc.Source.SshKey == nil {
			return fmt.Errorf("devcontainer '%s' uses SSH URL but does not specify ssh_key; SSH authentication is required for SSH repository URLs", dc.Id)
		}

		dc.Source.SshKey.Reference = strings.TrimSpace(dc.Source.SshKey.Reference)

		// Validate Reference is not empty
		if dc.Source.SshKey.Reference == "" {
			return fmt.Errorf("devcontainer '%s' has empty ssh_key.reference; a valid reference to the SSH key is required", dc.Id)
		}
	} else {
		// HTTPS URLs must not have SSH key
		if dc.Source.SshKey != nil {
			return fmt.Errorf("devcontainer '%s' uses HTTPS URL but specifies ssh_key; SSH authentication is not applicable for HTTPS repository URLs", dc.Id)
		}
	}

	return nil
}

func processRemoteAccess(dc *Devcontainer, nodes map[string]*Node) error {
	// RemoteAccess must not be nil
	if dc.RemoteAccess == nil {
		return fmt.Errorf("devcontainer '%s' has nil RemoteAccess; RemoteAccess is required", dc.Id)
	}

	node, exists := nodes[dc.NodeId]
	if !exists {
		return fmt.Errorf("devcontainer '%s' references non-existent node '%s'", dc.Id, dc.NodeId)
	}

	if dc.RemoteAccess.Ssh != nil {
		dc.RemoteAccess.Ssh.PublicSshKey = strings.TrimSpace(dc.RemoteAccess.Ssh.PublicSshKey)
	}

	// Validate OpenVSCode Server configuration
	if dc.RemoteAccess.OpenVsCodeServer != nil {
		// If node has DNS configured, OpenVSCode Server shouldn't have a custom port
		if node.DNS != nil && dc.RemoteAccess.OpenVsCodeServer.Port != nil {
			return fmt.Errorf("devcontainer '%s' specifies OpenVSCode Server port but its node has DNS configured; port must be omitted when DNS is configured at the node level", dc.Id)
		}
	}

	// If SSH is specified but no key is provided, use node's key as fallback
	if dc.RemoteAccess.Ssh != nil && dc.RemoteAccess.Ssh.PublicSshKey == "" {
		dc.RemoteAccess.Ssh.PublicSshKey = node.RemoteAccess.PublicSSHKey
	}

	return nil
}

func processDevcontainers(devcontainers []*Devcontainer, nodes map[string]*Node) error {
	// Track devcontainer IDs to ensure uniqueness
	seenIds := make(map[string]struct{})

	for i, dc := range devcontainers {
		// Trim string fields
		dc.Id = strings.TrimSpace(dc.Id)
		dc.NodeId = strings.TrimSpace(dc.NodeId)

		// Validate Id is not empty
		if dc.Id == "" {
			return fmt.Errorf("devcontainer at index %d has empty Id; Id must be specified", i)
		}

		// Check for duplicate devcontainer Ids
		if _, exists := seenIds[dc.Id]; exists {
			return fmt.Errorf("duplicate devcontainer Id '%s'; devcontainer Ids must be unique", dc.Id)
		}
		seenIds[dc.Id] = struct{}{}

		// Validate NodeId references a valid node
		if dc.NodeId == "" {
			return fmt.Errorf("devcontainer '%s' has empty node_id; node_id must be specified", dc.Id)
		}

		if _, exists := nodes[dc.NodeId]; !exists {
			return fmt.Errorf("devcontainer '%s' references non-existent node '%s'", dc.Id, dc.NodeId)
		}

		// Set default RemoteAccess if nil
		if dc.RemoteAccess == nil {
			dc.RemoteAccess = &DevcontainerRemoteAccess{
				OpenVsCodeServer: &DevcontainerOpenVSCodeServer{},
			}
		}

		// Process and validate Source configuration
		if err := processSource(dc); err != nil {
			return err
		}

		// Process and validate RemoteAccess configuration
		if err := processRemoteAccess(dc, nodes); err != nil {
			return err
		}
	}

	return nil
}

// DeserializeDenvclustrFile parses the JSON file data and returns
// the deserialized denvclustr file, informational messages, warnings, and any errors.
func DeserializeDenvclustrFile(data []byte) (*DenvclustrRoot, []string, []string, error) {
	var root DenvclustrRoot
	var warnings []string
	var infos []string

	// Deserialize the JSON data
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse denvclustr file: %w", err)
	}

	// Validate and process infrastructure
	if err := processInfrastructure(root.Infrastructure); err != nil {
		return nil, nil, nil, err
	}

	// Validate and process nodes
	if err := processNodes(root.Nodes, collectInfrastructureIds(root.Infrastructure)); err != nil {
		return nil, nil, nil, err
	}

	// Collect node data for devcontainer validation
	nodes := collectNodes(root.Nodes)

	// Validate and process devcontainers
	if err := processDevcontainers(root.Devcontainers, nodes); err != nil {
		return nil, nil, nil, err
	}

	return &root, infos, warnings, nil
}
