package denvclustr

import (
	"sync"

	"github.com/invopop/jsonschema"
)

var (
	denvclustrSchema     *jsonschema.Schema
	denvclustrSchemaOnce sync.Once
)

// Enum of supported infrastructure kinds.
type InfrastructureKind string

const (
	KindVm InfrastructureKind = "vm"
)

// Enum of supported infrastructure providers.
type Provider string

const (
	ProviderAws Provider = "aws"
)

// Enum of supported SSH key sources.
type SshKeySource string

const (
	SshKeySourceSecretsManager    SshKeySource = "secrets_manager"
	SshKeySourceSsmParameterStore SshKeySource = "ssm_parameter_store"
)

// DenvclustrRoot describes the top-level JSON structure:
type DenvclustrRoot struct {
	Name string `json:"name" jsonschema:"required,minLength=1" jsonschema_description:"Unique identifier for the cluster."`

	Infrastructure []*InfrastructureProvider `json:"infrastructure" jsonschema:"required,minItems=1" jsonschema_description:"List of infrastructure backends where nodes may be deployed."`

	Nodes []*Node `json:"nodes" jsonschema:"required,minItems=1" jsonschema_description:"List of nodes where devcontainers will be deployed."`

	Devcontainers []*Devcontainer `json:"devcontainers" jsonschema:"required,minItems=1" jsonschema_description:"List of devcontainers that will be deployed on nodes."`
}

// Infrastructure describes a single infrastructure backend.
type InfrastructureProvider struct {
	Id string `json:"id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Unique identifier for this infrastructure provider within the cluster. Must start with a letter or underscore, can contain alphanumeric characters, underscores, and hyphens. Cannot end with a hyphen."`

	Kind InfrastructureKind `json:"kind" jsonschema:"required,enum=vm" jsonschema_description:"Type of infrastructure. Currently only 'vm' is supported."`

	Provider Provider `json:"provider" jsonschema:"required,enum=aws" jsonschema_description:"Name of the platform. Currently only 'aws' is supported."`

	Region string `json:"region" jsonschema:"required,minLength=1" jsonschema_description:"Geographic location where resources will be deployed (e.g., 'us-west-2' for AWS). Must be a valid region identifier for the specified provider."`
}

type NodeRemoteAccess struct {
	PublicSSHKey string `json:"public_ssh_key" jsonschema:"required,minLength=1" jsonschema_description:"Path to local public SSH key."`
}

type NodeDNS struct {
	HighLevelDomain string `json:"high_level_domain" jsonschema:"required,minLength=1" jsonschema_description:"Top-level domain or subdomain what will be used for public devcontainer."`
}

type NodeProperties struct {
	InstanceType string `json:"instance_type" jsonschema:"required,minLength=1" jsonschema_description:"The machine type or class used to provision this node, specific to the target infrastructure."`
}

type DevcontainerOpenVSCodeServer struct {
	Port *int `json:"port,omitempty" jsonschema:"minimum=1024,maximum=65535" jsonschema_description:"TCP port used to expose the OpenVSCode Server interface. If not specified, an available port will be selected automatically. Must be omitted if DNS is configured at the node level."`
}

type DevcontainerSSH struct {
	Port         *int   `json:"port,omitempty" jsonschema:"minimum=1024,maximum=65535" jsonschema_description:"TCP port used for remote SSH access to the devcontainer. If not specified, an available port will be selected automatically."`
	PublicSshKey string `json:"public_ssh_key,omitempty" jsonschema_description:"Path to the local public SSH key used for authentication. If omitted, the public SSH key configured at the node level will be used."`
}

type DevcontainerRemoteAccess struct {
	OpenVsCodeServer *DevcontainerOpenVSCodeServer `json:"openvscode_server,omitempty" jsonschema_description:"Optional web-based IDE access to the devcontainer via OpenVSCode Server."`
	Ssh              *DevcontainerSSH              `json:"ssh,omitempty" jsonschema_description:"Optional access to the devcontainer via Secure Shell (SSH)."`
}

type DevcontainerSourceSSHKey struct {
	Reference string       `json:"reference" jsonschema:"required,minLength=1" jsonschema_description:"Reference identifier for the SSH key in the specified secret backend. Used to authenticate with private Git repositories."`
	Source    SshKeySource `json:"source" jsonschema:"required,enum=secrets_manager,enum=ssm_parameter_store" jsonschema_description:"Secret backend service where the private SSH key is stored. Must be either 'secrets_manager' or 'ssm_parameter_store'."`
}

type DevcontainerSource struct {
	URL              string                    `json:"url" jsonschema:"required,format=uri,minLength=1" jsonschema_description:"Git repository URL containing the devcontainer definition. For SSH URLs (starting with 'ssh://' or 'git@'), an SSH key must be provided."`
	Branch           string                    `json:"branch,omitempty" jsonschema_description:"Git branch to checkout from the repository. If not specified, the repository's default branch will be used."`
	DevcontainerPath string                    `json:"devcontainer_path,omitempty" jsonschema_description:"Relative path to the devcontainer definition within the repository. If not specified, the root directory will be used."`
	SshKey           *DevcontainerSourceSSHKey `json:"ssh_key,omitempty" jsonschema_description:"SSH key configuration for cloning from private Git repositories. Required for SSH URLs, must be omitted for HTTPS URLs."`
}

type Devcontainer struct {
	Id string `json:"id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Unique identifier of this devcontainer within the cluster."`

	NodeId string `json:"node_id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Identifier of the node that will host this devcontainer (must match an entry in the top‑level nodes list)."`

	Source *DevcontainerSource `json:"source" jsonschema:"required" jsonschema_description:"Reference to the source location containing the devcontainer definition and related files."`

	RemoteAccess *DevcontainerRemoteAccess `json:"remote_access,omitempty" jsonschema:"default={\"openvscode_server\":{}}" jsonschema_description:"Configuration for accessing the devcontainer remotely via SSH or a web-based IDE. OpenVSCode Server is enabled by default."`
}

type Node struct {
	Id string `json:"id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Unique identifier for this node. Must start with a letter or underscore, can contain alphanumeric characters, underscores, and hyphens. Cannot end with a hyphen."`

	// Reference to an entry in the 'infrastructure' array
	InfrastructureId string `json:"infrastructure_id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Reference to an entry in the 'infrastructure' array."`

	Properties NodeProperties `json:"properties" jsonschema:"required" jsonschema_description:"General technical configuration of the node."`

	RemoteAccess NodeRemoteAccess `json:"remote_access" jsonschema:"required" jsonschema_description:"Access configuration for the node."`

	DNS *NodeDNS `json:"dns,omitempty" jsonschema_description:"DNS configuration for this node."`
}

// GetSchema lazily builds and returns the JSON Schema for denvclustr file.
func GetSchema() *jsonschema.Schema {
	denvclustrSchemaOnce.Do(func() {
		r := &jsonschema.Reflector{
			DoNotReference: true,
		}
		denvclustrSchema = r.Reflect(&DenvclustrRoot{})
	})
	return denvclustrSchema
}
