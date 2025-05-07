package schema

// DenvclustrRoot describes the top-level JSON structure:
type DenvclustrRoot struct {
	Name           TrimmedString     `json:"name" jsonschema:"required,minLength=1" jsonschema_description:"Unique identifier for the cluster."`
	Infrastructure []*Infrastructure `json:"infrastructure" jsonschema:"required,minItems=1" jsonschema_description:"List of infrastructure backends where nodes may be deployed."`
	Nodes          []*Node           `json:"nodes" jsonschema:"required,minItems=1" jsonschema_description:"List of nodes where devcontainers will be deployed."`
	Devcontainers  []*Devcontainer   `json:"devcontainers" jsonschema:"required,minItems=1" jsonschema_description:"List of devcontainers that will be deployed on nodes."`
}
