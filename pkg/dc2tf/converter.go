package dc2tf

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	denvclustr "github.com/tropicaltux/denvclustr/pkg/denvclustr_specification"
	"github.com/zclconf/go-cty/cty"
)

// Converter handles the conversion from denvclustr specification to Terraform HCL
type Converter struct {
	spec *denvclustr.DenvclustrRoot
}

// NewConverter creates a new converter instance
func NewConverter(spec *denvclustr.DenvclustrRoot) *Converter {
	return &Converter{
		spec: spec,
	}
}

// ToTerraform converts the denvclustr specification to Terraform HCL
func (c *Converter) ToTerraform() (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	// Add provider blocks
	if err := c.addProviderBlocks(rootBody); err != nil {
		return nil, fmt.Errorf("failed to add provider blocks: %w", err)
	}

	// Add module blocks for each node
	if err := c.addModuleBlocks(rootBody); err != nil {
		return nil, fmt.Errorf("failed to add module blocks: %w", err)
	}

	// Add outputs
	if err := c.addOutputs(rootBody); err != nil {
		return nil, fmt.Errorf("failed to add outputs: %w", err)
	}

	return f, nil
}

// buildInfraMap creates a map of infrastructure providers by ID
func (c *Converter) buildInfraMap() map[string]*denvclustr.InfrastructureProvider {
	infraMap := make(map[string]*denvclustr.InfrastructureProvider)
	for _, infra := range c.spec.Infrastructure {
		infraMap[infra.Id] = infra
	}
	return infraMap
}

// buildNodeDevcontainersMap groups devcontainers by node ID
func (c *Converter) buildNodeDevcontainersMap() map[string][]*denvclustr.Devcontainer {
	nodeDevcontainers := make(map[string][]*denvclustr.Devcontainer)
	for _, dc := range c.spec.Devcontainers {
		nodeDevcontainers[dc.NodeId] = append(nodeDevcontainers[dc.NodeId], dc)
	}
	return nodeDevcontainers
}

// addProviderBlocks adds provider blocks for each infrastructure provider
func (c *Converter) addProviderBlocks(body *hclwrite.Body) error {
	for _, infra := range c.spec.Infrastructure {
		if infra.Provider == denvclustr.ProviderAws {
			// Add the provider block
			providerBlock := body.AppendNewBlock("provider", []string{"aws"})
			providerBody := providerBlock.Body()
			providerBody.SetAttributeValue("region", cty.StringVal(infra.Region))
			providerBody.SetAttributeValue("alias", cty.StringVal(infra.Id))
		} else {
			return fmt.Errorf("unsupported provider: %s", infra.Provider)
		}
	}

	return nil
}

// addModuleBlocks adds a module block for each node
func (c *Converter) addModuleBlocks(body *hclwrite.Body) error {
	// Get maps for infrastructure and devcontainers
	infraMap := c.buildInfraMap()
	nodeDevcontainers := c.buildNodeDevcontainersMap()

	// Create a module for each node
	for _, node := range c.spec.Nodes {
		// Check if node has devcontainers
		dcs, exists := nodeDevcontainers[node.Id]
		if !exists || len(dcs) == 0 {
			return fmt.Errorf("node %s has no devcontainers assigned to it", node.Id)
		}

		// Get the infrastructure for this node
		infra, exists := infraMap[node.InfrastructureId]
		if !exists {
			return fmt.Errorf("node %s references non-existent infrastructure %s", node.Id, node.InfrastructureId)
		}

		moduleBlock := body.AppendNewBlock("module", []string{node.Id})
		moduleBody := moduleBlock.Body()

		moduleBody.SetAttributeValue("source", cty.StringVal("github.com/tropicaltux/terraform-devcontainers"))
		moduleBody.SetAttributeValue("name", cty.StringVal(node.Id))
		moduleBody.SetAttributeValue("instance_type", cty.StringVal(node.Properties.InstanceType))
		moduleBody.SetAttributeValue("provider", cty.StringVal(fmt.Sprintf("aws.%s", infra.Id)))

		// Add devcontainers for this node
		if err := c.addDevcontainersForNode(moduleBody, dcs, node); err != nil {
			return fmt.Errorf("failed to add devcontainers for node %s: %w", node.Id, err)
		}

		// Add public SSH key for this node
		if err := c.addPublicSshKeyForNode(moduleBody, node); err != nil {
			return fmt.Errorf("failed to add SSH key for node %s: %w", node.Id, err)
		}

		// Add DNS if configured for this node
		if err := c.addDNSConfigForNode(moduleBody, node); err != nil {
			return fmt.Errorf("failed to add DNS config for node %s: %w", node.Id, err)
		}
	}

	return nil
}

// addDevcontainersForNode adds the devcontainers list for a specific node
func (c *Converter) addDevcontainersForNode(body *hclwrite.Body, devcontainers []*denvclustr.Devcontainer, node *denvclustr.Node) error {
	devcontainersBlock := body.AppendNewBlock("devcontainers", nil)
	devcontainersBody := devcontainersBlock.Body()

	for _, dc := range devcontainers {
		dcBlock := devcontainersBody.AppendNewBlock("", nil)
		dcBody := dcBlock.Body()

		dcBody.SetAttributeValue("id", cty.StringVal(dc.Id))

		// Add source block
		sourceBlock := dcBody.AppendNewBlock("source", nil)
		sourceBody := sourceBlock.Body()
		sourceBody.SetAttributeValue("url", cty.StringVal(dc.Source.URL))

		if dc.Source.Branch != "" {
			sourceBody.SetAttributeValue("branch", cty.StringVal(dc.Source.Branch))
		}

		if dc.Source.DevcontainerPath != "" {
			sourceBody.SetAttributeValue("devcontainer_path", cty.StringVal(dc.Source.DevcontainerPath))
		}

		// Add SSH key if needed
		if dc.Source.SshKey != nil {
			sshKeyBlock := sourceBody.AppendNewBlock("ssh_key", nil)
			sshKeyBody := sshKeyBlock.Body()
			sshKeyBody.SetAttributeValue("ref", cty.StringVal(dc.Source.SshKey.Reference))
			sshKeyBody.SetAttributeValue("src", cty.StringVal(string(dc.Source.SshKey.Source)))
		}

		// Add remote access configuration
		if dc.RemoteAccess != nil {
			remoteAccessBlock := dcBody.AppendNewBlock("remote_access", nil)
			remoteAccessBody := remoteAccessBlock.Body()

			// Add OpenVSCode Server config
			if dc.RemoteAccess.OpenVsCodeServer != nil {
				vsCodeBlock := remoteAccessBody.AppendNewBlock("openvscode_server", nil)
				vsCodeBody := vsCodeBlock.Body()

				if dc.RemoteAccess.OpenVsCodeServer.Port != nil {
					vsCodeBody.SetAttributeValue("port", cty.NumberIntVal(int64(*dc.RemoteAccess.OpenVsCodeServer.Port)))
				}
			}

			// Add SSH config
			if dc.RemoteAccess.Ssh != nil {
				sshBlock := remoteAccessBody.AppendNewBlock("ssh", nil)
				sshBody := sshBlock.Body()

				if dc.RemoteAccess.Ssh.Port != nil {
					sshBody.SetAttributeValue("port", cty.NumberIntVal(int64(*dc.RemoteAccess.Ssh.Port)))
				}

				// Add public SSH key if different from node's key
				if dc.RemoteAccess.Ssh.PublicSshKey != node.RemoteAccess.PublicSSHKey {
					pubKeyBlock := sshBody.AppendNewBlock("public_ssh_key", nil)
					pubKeyBody := pubKeyBlock.Body()
					pubKeyBody.SetAttributeValue("local_key_path", cty.StringVal(dc.RemoteAccess.Ssh.PublicSshKey))
				}
			}
		}
	}

	return nil
}

// addPublicSshKeyForNode adds the public SSH key configuration for a specific node
func (c *Converter) addPublicSshKeyForNode(body *hclwrite.Body, node *denvclustr.Node) error {
	sshKeyPath := node.RemoteAccess.PublicSSHKey
	if sshKeyPath == "" {
		return fmt.Errorf("node %s has no public SSH key configured", node.Id)
	}

	pubKeyBlock := body.AppendNewBlock("public_ssh_key", nil)
	pubKeyBody := pubKeyBlock.Body()
	pubKeyBody.SetAttributeValue("local_key_path", cty.StringVal(sshKeyPath))

	return nil
}

// addDNSConfigForNode adds DNS configuration for a specific node if available
func (c *Converter) addDNSConfigForNode(body *hclwrite.Body, node *denvclustr.Node) error {
	if node.DNS != nil && node.DNS.HighLevelDomain != "" {
		dnsBlock := body.AppendNewBlock("dns", nil)
		dnsBody := dnsBlock.Body()
		dnsBody.SetAttributeValue("high_level_domain", cty.StringVal(node.DNS.HighLevelDomain))
	}
	return nil
}

// addOutputs adds output blocks for the Terraform configuration
func (c *Converter) addOutputs(body *hclwrite.Body) error {
	// Get map for devcontainers
	nodeDevcontainers := c.buildNodeDevcontainersMap()

	// Add an output for each node with devcontainers
	for _, node := range c.spec.Nodes {
		// Skip nodes without devcontainers
		dcs, exists := nodeDevcontainers[node.Id]
		if !exists || len(dcs) == 0 {
			continue
		}

		// Create output name based on node ID
		outputName := fmt.Sprintf("devcontainers_%s_output", node.Id)

		outputBlock := body.AppendNewBlock("output", []string{outputName})
		outputBody := outputBlock.Body()
		outputBody.SetAttributeValue("value", cty.ObjectVal(map[string]cty.Value{
			"module": cty.StringVal(fmt.Sprintf("module.%s", node.Id)),
		}))
		outputBody.SetAttributeValue("sensitive", cty.BoolVal(true))
	}

	return nil
}
