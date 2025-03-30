package codecluster

import (
	"sync"

	"github.com/invopop/jsonschema"
)

var (
	codeclusterSchema     *jsonschema.Schema
	codeclusterSchemaOnce sync.Once
)

// CodeclusterConfiguration describes the top-level JSON structure:
type CodeclusterConfiguration struct {
	// "devcontainers" is a map whose keys are devcontainer id in cluster,
	// and values are Devcontainer objects.
	Devcontainers map[string]*Devcontainer `json:"devcontainers,omitempty"`
}

// Devcontainer describes each named devcontainer entry.
//
// - Exactly ONE of "repository_url" or "repository_path" must be set.
// - "branch" defaults to "main" and will ignored if "repository_path" is set.
type Devcontainer struct {
	// If not set, defaults to the devcontainer id.
	Name string `json:"name,omitempty" jsonschema_description:"Name of the devcontainer in cluster. This parameter intended for a clearer, more human-friendly name. If not set, defaults to the devcontainer id."`

	// Not required.
	Description string `json:"description,omitempty" jsonschema_description:"Description of the devcontainer in cluster. This parameter intended for a more detailed description of the devcontainer."`

	// "oneof_required=url" means this field is required for sub-schema "url"
	RepositoryURL string `json:"repository_url,omitempty" jsonschema:"oneof_required=url" jsonschema_description:"URL to the repository. Cannot coexist with repository_path."`

	// "oneof_required=path" means this field is required for sub-schema "path"
	RepositoryPath string `json:"repository_path,omitempty" jsonschema:"oneof_required=path" jsonschema_description:"Local path to the repository. Cannot coexist with repository_url."`

	// If not set, defaults to "main". Ignored if repository_path is set.
	Branch string `json:"branch,omitempty" jsonschema:"default=main" jsonschema_description:"Branch name (default 'main'). Optional, ignored if repository_path is set."`
}

// GetCodeclusterSchema lazily builds and returns the JSON Schema for codecluster configuration file.
func GetCodeclusterSchema() *jsonschema.Schema {
	codeclusterSchemaOnce.Do(func() {
		r := &jsonschema.Reflector{}
		codeclusterSchema = r.Reflect(&CodeclusterConfiguration{})
	})
	return codeclusterSchema
}
