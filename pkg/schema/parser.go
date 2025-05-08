package schema

import (
	"encoding/json"
	"fmt"
)

// Parse deserializes the provided raw data into a DenvclustrRoot structure,
// and validates the deserialized data. It returns the validated DenvclustrRoot or an error.
func Parse(data []byte) (*DenvclustrRoot, error) {

	var rawRoot any
	if err := json.Unmarshal(data, &rawRoot); err != nil {
		return nil, fmt.Errorf("unmarshal schema: %w", err)
	}

	if err := validateSchema(rawRoot); err != nil {
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
