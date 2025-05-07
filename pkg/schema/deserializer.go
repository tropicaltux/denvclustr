package schema

import (
	"encoding/json"
)

// DeserializeDenvclustrFile parses JSON into DenvclustrRoot.
func deserializeDenvclustrFile(data []byte) (*DenvclustrRoot, error) {
	var root DenvclustrRoot
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	return &root, nil
}
