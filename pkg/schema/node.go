package schema

type NodeProperties struct {
	InstanceType TrimmedString `json:"instance_type" jsonschema:"required,minLength=1" jsonschema_description:"The machine type or class used to provision this node, specific to the target infrastructure."`
}

type NodeRemoteAccess struct {
	PublicSSHKey TrimmedString `json:"public_ssh_key" jsonschema:"required,minLength=1" jsonschema_description:"Path to local public SSH key."`
}

type NodeDNS struct {
	HighLevelDomain TrimmedString `json:"high_level_domain" jsonschema:"required,minLength=1" jsonschema_description:"Top-level domain or subdomain that will be used for public devcontainer."`
}

type Node struct {
	Id               TrimmedString    `json:"id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Unique identifier for this node. Must start with a letter or underscore, can contain alphanumeric characters, underscores, and hyphens. Cannot end with a hyphen."`
	InfrastructureId TrimmedString    `json:"infrastructure_id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Reference to an entry in the 'infrastructure' array."`
	Properties       NodeProperties   `json:"properties" jsonschema:"required" jsonschema_description:"General technical configuration of the node."`
	RemoteAccess     NodeRemoteAccess `json:"remote_access" jsonschema:"required" jsonschema_description:"Access configuration for the node."`
	DNS              *NodeDNS         `json:"dns,omitempty" jsonschema_description:"DNS configuration for this node."`
}
