package schema

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/*.json
var testDataFS embed.FS

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		filename      string
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid minimal config",
			filename:    "valid_minimal.json",
			expectError: false,
		},
		{
			name:          "Invalid JSON syntax",
			filename:      "invalid_json_syntax.json",
			expectError:   true,
			errorContains: "invalid",
		},
		{
			name:          "Missing required field - nodes",
			filename:      "missing_required.json",
			expectError:   true,
			errorContains: "missing property 'nodes'",
		},
		{
			name:          "Empty arrays",
			filename:      "empty_arrays.json",
			expectError:   true,
			errorContains: "minItems",
		},
		{
			name:          "Duplicate infrastructure ID",
			filename:      "duplicate_infrastructure_id.json",
			expectError:   true,
			errorContains: "infrastructure id \"infrastructure1\" is duplicated",
		},
		{
			name:          "Unreferenced infrastructure",
			filename:      "unreferenced_infrastructure.json",
			expectError:   true,
			errorContains: "not referenced by any node",
		},
		{
			name:          "Duplicate node ID",
			filename:      "duplicate_node_id.json",
			expectError:   true,
			errorContains: "node id \"node1\" is duplicated",
		},
		{
			name:          "Unreferenced node",
			filename:      "unreferenced_node.json",
			expectError:   true,
			errorContains: "not referenced by any devcontainer",
		},
		{
			name:          "Duplicate devcontainer ID",
			filename:      "duplicate_dev_id.json",
			expectError:   true,
			errorContains: "devcontainer id \"devcontainer1\" is duplicated",
		},
		{
			name:          "SSH URL without SSH key",
			filename:      "ssh_without_key.json",
			expectError:   true,
			errorContains: "ssh_key must be provided for SSH-based URLs",
		},
		{
			name:          "Non-SSH URL with SSH key",
			filename:      "nonssh_with_key.json",
			expectError:   true,
			errorContains: "devcontainer",
		},
		{
			name:          "Empty file",
			filename:      "empty_file.json",
			expectError:   true,
			errorContains: "unexpected end of JSON input",
		},
		{
			name:        "Valid complete config",
			filename:    "valid_complete.json",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read test file using the embedded FS
			data, err := testDataFS.ReadFile("testdata/" + tt.filename)
			assert.NoError(t, err, "Failed to read test file")

			// Call the function being tested
			result, err := Parse(data)

			// Check if error behavior matches expected
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Additional validation for successful cases
				assert.NotEmpty(t, result.Name)
				assert.NotEmpty(t, result.Infrastructure)
				assert.NotEmpty(t, result.Nodes)
				assert.NotEmpty(t, result.Devcontainers)
			}
		})
	}
}
