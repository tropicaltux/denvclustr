package dc2tf

import (
	"fmt"

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
		moduleBody.SetAttributeValue("provider", cty.StringVal(fmt.Sprintf("aws.%s", infrastructureById[string(node.InfrastructureId)].Id)))

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
	parent := body.AppendNewBlock("devcontainers", nil).Body()
	for _, devcontainer := range devcontainers {
		devcontainerBody := parent.AppendNewBlock("", nil).Body()
		devcontainerBody.SetAttributeValue("id", cty.StringVal(string(devcontainer.Id)))

		srcBody := devcontainerBody.AppendNewBlock("source", nil).Body()
		srcBody.SetAttributeValue("url", cty.StringVal(string(devcontainer.Source.URL)))
		if devcontainer.Source.Branch != "" {
			srcBody.SetAttributeValue("branch", cty.StringVal(string(devcontainer.Source.Branch)))
		}
		if devcontainer.Source.DevcontainerPath != "" {
			srcBody.SetAttributeValue("devcontainer_path", cty.StringVal(string(devcontainer.Source.DevcontainerPath)))
		}
		if devcontainer.Source.SshKey != nil {
			sshKeyBody := srcBody.AppendNewBlock("ssh_key", nil).Body()
			sshKeyBody.SetAttributeValue("ref", cty.StringVal(string(devcontainer.Source.SshKey.Reference)))
			sshKeyBody.SetAttributeValue("src", cty.StringVal(string(devcontainer.Source.SshKey.Source)))
		}

		remoteAccessBody := devcontainerBody.AppendNewBlock("remote_access", nil).Body()
		if devcontainer.RemoteAccess.OpenVsCodeServer != nil {
			vscodeServerBody := remoteAccessBody.AppendNewBlock("openvscode_server", nil).Body()
			if devcontainer.RemoteAccess.OpenVsCodeServer.Port != nil {
				vscodeServerBody.SetAttributeValue("port", cty.NumberIntVal(int64(*devcontainer.RemoteAccess.OpenVsCodeServer.Port)))
			}
		}
		if devcontainer.RemoteAccess.Ssh != nil {
			sshBody := remoteAccessBody.AppendNewBlock("ssh", nil).Body()
			if devcontainer.RemoteAccess.Ssh.Port != nil {
				sshBody.SetAttributeValue("port", cty.NumberIntVal(int64(*devcontainer.RemoteAccess.Ssh.Port)))
			}
			if devcontainer.RemoteAccess.Ssh.PublicSshKey != node.RemoteAccess.PublicSSHKey {
				pk := sshBody.AppendNewBlock("public_ssh_key", nil).Body()
				pk.SetAttributeValue("local_key_path", cty.StringVal(string(devcontainer.RemoteAccess.Ssh.PublicSshKey)))
			}
		}
	}
	return nil
}

func (c *converter) writeNodeSSH(body *hclwrite.Body, node *schema.Node) error {
	publicSshKeyBody := body.AppendNewBlock("public_ssh_key", nil).Body()
	publicSshKeyBody.SetAttributeValue("local_key_path", cty.StringVal(string(node.RemoteAccess.PublicSSHKey)))
	return nil
}

func (c *converter) writeDNS(body *hclwrite.Body, node *schema.Node) {
	if node.DNS != nil && node.DNS.HighLevelDomain != "" {
		dnsBody := body.AppendNewBlock("dns", nil).Body()
		dnsBody.SetAttributeValue("high_level_domain", cty.StringVal(string(node.DNS.HighLevelDomain)))
	}
}
