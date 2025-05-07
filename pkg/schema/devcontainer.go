package schema

// Enum of supported SSH key sources.
type SshKeySource string

const (
	SshKeySourceSecretsManager    SshKeySource = "secrets_manager"
	SshKeySourceSsmParameterStore SshKeySource = "ssm_parameter_store"
)

type DevcontainerOpenVSCodeServer struct {
	Port *int `json:"port,omitempty" jsonschema:"minimum=1024,maximum=65535" jsonschema_description:"TCP port used to expose the OpenVSCode Server interface. If not specified, an available port will be selected automatically. Must be omitted if DNS is configured at the node level."`
}

type DevcontainerSSH struct {
	Port         *int          `json:"port,omitempty" jsonschema:"minimum=1024,maximum=65535" jsonschema_description:"TCP port used for remote SSH access to the devcontainer. If not specified, an available port will be selected automatically."`
	PublicSshKey TrimmedString `json:"public_ssh_key,omitempty" jsonschema_description:"Path to the local public SSH key used for authentication. If omitted, the public SSH key configured at the node level will be used."`
}

type DevcontainerRemoteAccess struct {
	OpenVsCodeServer *DevcontainerOpenVSCodeServer `json:"openvscode_server,omitempty" jsonschema_description:"Optional web-based IDE access to the devcontainer via OpenVSCode Server."`
	Ssh              *DevcontainerSSH              `json:"ssh,omitempty" jsonschema_description:"Optional access to the devcontainer via Secure Shell (SSH)."`
}

type DevcontainerSourceSSHKey struct {
	Reference TrimmedString `json:"reference" jsonschema:"required,minLength=1" jsonschema_description:"Reference identifier for the SSH key in the specified secret backend. Used to authenticate with private Git repositories."`
	Source    SshKeySource  `json:"source" jsonschema:"required,enum=secrets_manager,enum=ssm_parameter_store" jsonschema_description:"Secret backend service where the private SSH key is stored. Must be either 'secrets_manager' or 'ssm_parameter_store'."`
}

type DevcontainerSource struct {
	URL              TrimmedString             `json:"url" jsonschema:"required,format=uri,minLength=1" jsonschema_description:"Git repository URL containing the devcontainer definition. For SSH URLs (starting with 'ssh://' or 'git@'), an SSH key must be provided."`
	Branch           TrimmedString             `json:"branch,omitempty" jsonschema_description:"Git branch to checkout. If not specified, the repository's default branch will be used."`
	DevcontainerPath TrimmedString             `json:"devcontainer_path,omitempty" jsonschema_description:"Relative path to the devcontainer definition within the repository. If not specified, the root directory will be used."`
	SshKey           *DevcontainerSourceSSHKey `json:"ssh_key,omitempty" jsonschema_description:"SSH key configuration for cloning from private Git repositories. Required for SSH URLs, must be omitted for HTTPS URLs."`
}

type Devcontainer struct {
	Id           TrimmedString             `json:"id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Unique identifier of this devcontainer within the cluster."`
	NodeId       TrimmedString             `json:"node_id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9_]$" jsonschema_description:"Identifier of the node that will host this devcontainer (must match an entry in the top‑level nodes list)."`
	Source       *DevcontainerSource       `json:"source" jsonschema:"required" jsonschema_description:"Reference to the source location containing the devcontainer definition and related files."`
	RemoteAccess *DevcontainerRemoteAccess `json:"remote_access,omitempty" jsonschema_description:"Configuration for accessing the devcontainer remotely via SSH or a web-based IDE. OpenVSCode Server is enabled by default."`
}
