package schema

import (
	"sync"

	"github.com/invopop/jsonschema"
)

var (
	denvclustrSchema     *jsonschema.Schema
	denvclustrSchemaOnce sync.Once
)

func addCustomValidations(schema *jsonschema.Schema) {
	properties := schema.Properties
	if property, ok := properties.Get("infrastructure"); ok && property != nil {
		property.UniqueItems = true
	}
	if property, ok := properties.Get("nodes"); ok && property != nil {
		property.UniqueItems = true
	}
	if property, ok := properties.Get("devcontainers"); ok && property != nil {
		property.UniqueItems = true
	}
}

// GetSchema returns the JSON Schema for DenvclustrRoot.
func GetSchema() *jsonschema.Schema {
	denvclustrSchemaOnce.Do(func() {
		reflector := &jsonschema.Reflector{
			DoNotReference:             true,
			RequiredFromJSONSchemaTags: true,
			ExpandedStruct:             true,
			AllowAdditionalProperties:  false,
		}
		denvclustrSchema = reflector.Reflect(&DenvclustrRoot{})
		addCustomValidations(denvclustrSchema)
	})
	return denvclustrSchema
}
