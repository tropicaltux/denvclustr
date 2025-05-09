package schema

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

// Infrastructure describes a single infrastructure backend.
type Infrastructure struct {
	Id       TrimmedString      `json:"id" jsonschema:"required,minLength=1,pattern=^[_a-zA-Z][a-zA-Z0-9-]*[a-zA-Z0-9]$" jsonschema_description:"Unique identifier for this infrastructure provider within the cluster. Must start with a letter or underscore, can contain alphanumeric characters, underscores, and hyphens. Cannot end with a hyphen."`
	Kind     InfrastructureKind `json:"kind" jsonschema:"required,enum=vm" jsonschema_description:"Type of infrastructure. Currently only 'vm' is supported."`
	Provider Provider           `json:"provider" jsonschema:"required,enum=aws" jsonschema_description:"Name of the platform. Currently only 'aws' is supported."`
	Region   TrimmedString      `json:"region" jsonschema:"required,minLength=1" jsonschema_description:"Geographic location where resources will be deployed (e.g., 'us-west-2' for AWS). Must be a valid region identifier for the specified provider."`
}
