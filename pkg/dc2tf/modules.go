package dc2tf

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/tropicaltux/denvclustr/pkg/schema"
	"github.com/zclconf/go-cty/cty"
)

func (c *converter) addModules(body *hclwrite.Body) error {
	infrastructureById := map[string]*schema.Infrastructure{}
	for _, infrastructure := range c.root.Infrastructure {
		infrastructureById[string(infrastructure.Id)] = infrastructure
	}

	devcontainerByNode := map[string][]*schema.Devcontainer{}
	for _, devcontainer := range c.root.Devcontainers {
		devcontainerByNode[string(devcontainer.NodeId)] = append(devcontainerByNode[string(devcontainer.NodeId)], devcontainer)
	}

	for _, node := range c.root.Nodes {
		devcontainers := devcontainerByNode[string(node.Id)]

		moduleBody := body.AppendNewBlock("module", []string{string(node.Id)}).Body()
		moduleBody.SetAttributeValue("source", cty.StringVal("github.com/tropicaltux/terraform-devcontainers"))
		moduleBody.SetAttributeValue("name", cty.StringVal(string(node.Id)))
		moduleBody.SetAttributeValue("instance_type", cty.StringVal(string(node.Properties.InstanceType)))

		// Replace provider attribute with providers attribute using raw tokens
		providerTokens := hclwrite.Tokens{}

		// Opening brace
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrace,
			Bytes: []byte("{"),
		})

		// Key "aws"
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenIdent,
			Bytes:        []byte("aws"),
			SpacesBefore: 1,
		})

		// Equal sign
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenEqual,
			Bytes:        []byte("="),
			SpacesBefore: 1,
		})

		// Reference to aws provider without quotes
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenIdent,
			Bytes:        []byte("aws"),
			SpacesBefore: 1,
		})
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:  hclsyntax.TokenDot,
			Bytes: []byte("."),
		})
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(string(infrastructureById[string(node.InfrastructureId)].Id)),
		})

		// Closing brace
		providerTokens = append(providerTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenCBrace,
			Bytes:        []byte("}"),
			SpacesBefore: 1,
		})

		moduleBody.SetAttributeRaw("providers", providerTokens)

		if err := c.writeDevcontainers(moduleBody, devcontainers, node); err != nil {
			return err
		}
		if err := c.writeNodeSSH(moduleBody, node); err != nil {
			return err
		}
		c.writeDNS(moduleBody, node)
	}
	return nil
}

func (c *converter) writeDevcontainers(body *hclwrite.Body, devcontainers []*schema.Devcontainer, node *schema.Node) error {
	var devcontainerItems []cty.Value

	for _, devcontainer := range devcontainers {
		devcontainerMap := map[string]cty.Value{}
		devcontainerMap["id"] = cty.StringVal(string(devcontainer.Id))

		sourceMap := map[string]cty.Value{}
		sourceMap["url"] = cty.StringVal(string(devcontainer.Source.URL))
		if devcontainer.Source.Branch != "" {
			sourceMap["branch"] = cty.StringVal(string(devcontainer.Source.Branch))
		}
		if devcontainer.Source.DevcontainerPath != "" {
			sourceMap["devcontainer_path"] = cty.StringVal(string(devcontainer.Source.DevcontainerPath))
		}
		if devcontainer.Source.SshKey != nil {
			sshKeyMap := map[string]cty.Value{}
			sshKeyMap["ref"] = cty.StringVal(string(devcontainer.Source.SshKey.Reference))
			sshKeyMap["src"] = cty.StringVal(string(devcontainer.Source.SshKey.Source))
			sourceMap["ssh_key"] = cty.ObjectVal(sshKeyMap)
		}

		// Create remote_access block - adding source before remote_access to match order in etalon
		devcontainerMap["source"] = cty.ObjectVal(sourceMap)

		if devcontainer.RemoteAccess != nil {
			remoteAccessMap := map[string]cty.Value{}
			if devcontainer.RemoteAccess.OpenVsCodeServer != nil {
				vscodeServerMap := map[string]cty.Value{}
				if devcontainer.RemoteAccess.OpenVsCodeServer.Port != nil {
					vscodeServerMap["port"] = cty.NumberIntVal(int64(*devcontainer.RemoteAccess.OpenVsCodeServer.Port))
				}
				remoteAccessMap["openvscode_server"] = cty.ObjectVal(vscodeServerMap)
			}
			if devcontainer.RemoteAccess.Ssh != nil {
				sshMap := map[string]cty.Value{}
				if devcontainer.RemoteAccess.Ssh.Port != nil {
					sshMap["port"] = cty.NumberIntVal(int64(*devcontainer.RemoteAccess.Ssh.Port))
				}
				if devcontainer.RemoteAccess.Ssh.PublicSshKey != "" && devcontainer.RemoteAccess.Ssh.PublicSshKey != node.RemoteAccess.PublicSSHKey {
					sshMap["public_ssh_key"] = cty.ObjectVal(map[string]cty.Value{
						"local_key_path": cty.StringVal(string(devcontainer.RemoteAccess.Ssh.PublicSshKey)),
					})
				}
				remoteAccessMap["ssh"] = cty.ObjectVal(sshMap)
			}
			devcontainerMap["remote_access"] = cty.ObjectVal(remoteAccessMap)
		}

		devcontainerItems = append(devcontainerItems, cty.ObjectVal(devcontainerMap))
	}

	// Handle empty devcontainers case
	if len(devcontainerItems) == 0 {
		body.SetAttributeValue("devcontainers", cty.ListValEmpty(cty.DynamicPseudoType))
	} else {
		body.SetAttributeValue("devcontainers", cty.ListVal(devcontainerItems))
	}

	return nil
}

func (c *converter) writeNodeSSH(body *hclwrite.Body, node *schema.Node) error {
	body.SetAttributeValue("public_ssh_key", cty.ObjectVal(map[string]cty.Value{
		"local_key_path": cty.StringVal(string(node.RemoteAccess.PublicSSHKey)),
	}))
	return nil
}

func (c *converter) writeDNS(body *hclwrite.Body, node *schema.Node) {
	if node.DNS != nil && node.DNS.HighLevelDomain != "" {
		body.SetAttributeValue("dns", cty.ObjectVal(map[string]cty.Value{
			"high_level_domain": cty.StringVal(string(node.DNS.HighLevelDomain)),
		}))
	}
}
