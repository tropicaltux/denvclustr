package schema

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	validatorOnce sync.Once
	validator     *jsonschema.Schema
	validatorErr  error
)

func getValidator() (*jsonschema.Schema, error) {
	validatorOnce.Do(func() {
		schemaBytes, err := json.Marshal(GetSchema())
		if err != nil {
			validatorErr = fmt.Errorf("marshal schema: %w", err)
			return
		}

		var compilerFormatSchema any
		if err := json.Unmarshal(schemaBytes, &compilerFormatSchema); err != nil {
			validatorErr = fmt.Errorf("unmarshal schema: %w", err)
			return
		}

		c := jsonschema.NewCompiler()

		if err := c.AddResource("schema.go", compilerFormatSchema); err != nil {
			validatorErr = fmt.Errorf("add resource: %w", err)
			return
		}

		validator, err = c.Compile("schema.go")
		if err != nil {
			validatorErr = fmt.Errorf("compile schema: %w", err)
		}
	})
	return validator, validatorErr
}

func validateSchema(data any) error {
	validator, err := getValidator()
	if err != nil {
		return fmt.Errorf("get validator: %w", err)
	}

	if err := validator.Validate(data); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	return nil
}

// Validate a fully‑deserialized spec.
func validateDeserialized(root *DenvclustrRoot) error {
	if err := validateInfrastructure(root); err != nil {
		return err
	}
	if err := validateNodes(root); err != nil {
		return err
	}
	if err := validateDevcontainers(root); err != nil {
		return err
	}
	return nil
}

func validateInfrastructure(root *DenvclustrRoot) error {
	seenIds := make(map[string]struct{})
	type tuple struct {
		Kind     InfrastructureKind
		Provider Provider
		Region   string
	}
	seenTuple := make(map[tuple]string)

	for i, infrastructure := range root.Infrastructure {
		id := string(infrastructure.Id)

		if id == "" {
			return fmt.Errorf("infrastructure[%d]: id is missing", i)
		}

		if _, exists := seenIds[id]; exists {
			return fmt.Errorf("infrastructure id %q is duplicated", id)
		}
		seenIds[id] = struct{}{}

		if infrastructure.Region == "" {
			return fmt.Errorf("infrastructure %q: region is missing", id)
		}

		t := tuple{infrastructure.Kind, infrastructure.Provider, string(infrastructure.Region)}
		if previousId, exists := seenTuple[t]; exists {
			return fmt.Errorf(
				"infrastructure %q: duplicate combination of kind %q, provider %q, region %q (also defined by %q)",
				id, infrastructure.Kind, infrastructure.Provider, infrastructure.Region, previousId,
			)
		}
		seenTuple[t] = id
	}

	// Ensure each infrastructure is referenced by at least one node
	referenced := make(map[string]struct{})
	for _, node := range root.Nodes {
		referenced[string(node.InfrastructureId)] = struct{}{}
	}
	for id := range seenIds {
		if _, ok := referenced[id]; !ok {
			return fmt.Errorf("infrastructure %q is not referenced by any node", id)
		}
	}
	return nil
}

func validateNodes(root *DenvclustrRoot) error {
	seen := make(map[string]struct{})

	for i, node := range root.Nodes {
		id := string(node.Id)

		if id == "" {
			return fmt.Errorf("node[%d]: id is missing", i)
		}

		if _, exists := seen[id]; exists {
			return fmt.Errorf("node id %q is duplicated", id)
		}
		seen[id] = struct{}{}

		if node.InfrastructureId == "" {
			return fmt.Errorf("node %q: infrastructure_id is missing", id)
		}

		if node.Properties.InstanceType == "" {
			return fmt.Errorf("node %q: instance_type is missing", id)
		}

		if node.RemoteAccess.PublicSSHKey == "" {
			return fmt.Errorf("node %q: public_ssh_key is missing", id)
		}

		if node.DNS != nil && node.DNS.HighLevelDomain == "" {
			return fmt.Errorf(
				"node %q: dns.high_level_domain must be provided when DNS settings exist", id,
			)
		}
	}

	// Ensure each node is referenced by at least one devcontainer
	referenced := make(map[string]struct{})
	for _, devcontainer := range root.Devcontainers {
		referenced[string(devcontainer.NodeId)] = struct{}{}
	}
	for id := range seen {
		if _, ok := referenced[id]; !ok {
			return fmt.Errorf("node %q is not referenced by any devcontainer", id)
		}
	}
	return nil
}

func validateDevcontainers(root *DenvclustrRoot) error {
	seen := make(map[string]struct{})
	nodeMap := collectNodeMap(root)

	for i, devcontainer := range root.Devcontainers {
		id := string(devcontainer.Id)

		if id == "" {
			return fmt.Errorf("devcontainer[%d]: id is missing", i)
		}

		if _, exists := seen[id]; exists {
			return fmt.Errorf("devcontainer id %q is duplicated", id)
		}
		seen[id] = struct{}{}

		if devcontainer.NodeId == "" {
			return fmt.Errorf("devcontainer %q: node_id is missing", id)
		}

		if _, ok := nodeMap[string(devcontainer.NodeId)]; !ok {
			return fmt.Errorf("devcontainer %q: refers to unknown node_id %q", id, devcontainer.NodeId)
		}

		if devcontainer.Source == nil || devcontainer.Source.URL == "" {
			return fmt.Errorf("devcontainer %q: source.url is required and must be valid", id)
		}

		isSSH := strings.HasPrefix(string(devcontainer.Source.URL), "ssh://") || strings.HasPrefix(string(devcontainer.Source.URL), "git@")

		if isSSH && devcontainer.Source.SshKey == nil {
			return fmt.Errorf("devcontainer %q: ssh_key must be provided for SSH-based URLs", id)
		}

		if !isSSH && devcontainer.Source.SshKey != nil {
			return fmt.Errorf("devcontainer %q: ssh_key must not be used with non-SSH URLs", id)
		}
	}
	return nil
}
