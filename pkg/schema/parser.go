package schema

// Parse deserializes the provided raw data into a DenvclustrRoot structure,
// and validates the deserialized data. It returns the validated DenvclustrRoot or an error.
func Parse(data []byte) (*DenvclustrRoot, error) {
	// Deserialize the data
	root, err := deserializeDenvclustrFile(data)
	if err != nil {
		return nil, err
	}

	// Validate the deserialized data
	if err := Validate(root); err != nil {
		return nil, err
	}

	return root, nil
}
