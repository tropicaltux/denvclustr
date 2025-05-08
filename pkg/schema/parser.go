package schema

import (
	"encoding/json"
)

// Parse deserializes the provided raw data into a DenvclustrRoot structure,
// and validates the deserialized data. It returns the validated DenvclustrRoot or an error.
func Parse(data []byte) (*DenvclustrRoot, error) {

	// First validate against JSON schema
	var jsonObj map[string]any
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, err
	}

	if err := validateSchema(jsonObj); err != nil {
		return nil, err
	}

	// Deserialize the data
	root, err := deserializeDenvclustrFile(data)
	if err != nil {
		return nil, err
	}

	// Validate the deserialized data
	if err := validateDeserialized(root); err != nil {
		return nil, err
	}

	return root, nil
}
