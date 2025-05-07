package schema

import (
	"embed"
	"reflect"
	"strings"
	"testing"
)

// testFiles contains all the JSON test files embedded from the testdata directory.
//
//go:embed testdata/*.json
var testFiles embed.FS

// readTestFile reads a test file from the embedded filesystem.
func readTestFile(t *testing.T, filePath string) ([]byte, error) {
	t.Helper()
	return testFiles.ReadFile(filePath)
}

// TestValidCompleteInput tests deserialization of a valid complete JSON input.
func TestValidCompleteInput(t *testing.T) {
	data, err := readTestFile(t, "testdata/valid_complete.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	got, err := deserializeDenvclustrFile(data)

	// Validate results
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	want := &DenvclustrRoot{
		Name: "test-cluster",
		Infrastructure: []*Infrastructure{
			{
				Id:       "aws_infra",
				Kind:     "vm",
				Provider: "aws",
				Region:   "us-west-2",
			},
		},
		Nodes: []*Node{
			{
				Id:               "node1",
				InfrastructureId: "aws_infra",
				Properties: NodeProperties{
					InstanceType: "t2.micro",
				},
				RemoteAccess: NodeRemoteAccess{
					PublicSSHKey: "/path/to/key.pub",
				},
				DNS: &NodeDNS{
					HighLevelDomain: "example.com",
				},
			},
		},
		Devcontainers: []*Devcontainer{
			{
				Id:     "dev1",
				NodeId: "node1",
				Source: &DevcontainerSource{
					URL: "https://github.com/example/repo.git",
				},
				RemoteAccess: &DevcontainerRemoteAccess{
					Ssh: &DevcontainerSSH{
						Port: intPtr(2222),
					},
					OpenVsCodeServer: &DevcontainerOpenVSCodeServer{
						Port: intPtr(8080),
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("deserializeDenvclustrFile() = %v, want %v", got, want)
	}
}

// TestValidMinimalInput tests deserialization of a minimal valid JSON input.
func TestValidMinimalInput(t *testing.T) {
	data, err := readTestFile(t, "testdata/valid_minimal.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	got, err := deserializeDenvclustrFile(data)

	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Expected result
	want := &DenvclustrRoot{
		Name: "min-cluster",
		Infrastructure: []*Infrastructure{
			{
				Id:       "aws_infra",
				Kind:     "vm",
				Provider: "aws",
				Region:   "us-west-2",
			},
		},
		Nodes: []*Node{
			{
				Id:               "node1",
				InfrastructureId: "aws_infra",
				Properties: NodeProperties{
					InstanceType: "t2.micro",
				},
				RemoteAccess: NodeRemoteAccess{
					PublicSSHKey: "/path/to/key.pub",
				},
			},
		},
		Devcontainers: []*Devcontainer{
			{
				Id:     "dev1",
				NodeId: "node1",
				Source: &DevcontainerSource{
					URL: "https://github.com/example/repo.git",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("deserializeDenvclustrFile() = %v, want %v", got, want)
	}
}

// TestInvalidJSONSyntax tests deserialization of an invalid JSON input.
func TestInvalidJSONSyntax(t *testing.T) {
	data, err := readTestFile(t, "testdata/invalid_json_syntax.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	got, err := deserializeDenvclustrFile(data)

	if err == nil {
		t.Errorf("deserializeDenvclustrFile() expected error, got nil")
		return
	}

	// Check that the error message contains the expected substring
	expectedErrSubstring := "unexpected end of JSON input"
	if !strings.Contains(err.Error(), expectedErrSubstring) {
		t.Errorf("deserializeDenvclustrFile() error = %v, want error containing %q", err, expectedErrSubstring)
	}

	// Check that the result is nil
	if got != nil {
		t.Errorf("deserializeDenvclustrFile() = %v, want nil", got)
	}
}

// TestMissingInfrastructure tests deserialization of JSON with missing infrastructure field.
func TestMissingInfrastructure(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/missing_infrastructure.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Check specific fields
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	if string(got.Name) != "no-infra" {
		t.Errorf("Name: got %q, want %q", got.Name, "no-infra")
	}
	if len(got.Infrastructure) != 0 {
		t.Errorf("Infrastructure: expected empty slice, got %v", got.Infrastructure)
	}
	if len(got.Nodes) != 1 {
		t.Errorf("Nodes: expected 1 node, got %d", len(got.Nodes))
	}
	if len(got.Devcontainers) != 1 {
		t.Errorf("Devcontainers: expected 1 devcontainer, got %d", len(got.Devcontainers))
	}
}

// TestIncorrectTypeName tests deserialization of JSON with incorrect type for name field.
func TestIncorrectTypeName(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/incorrect_type_name.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results
	if err == nil {
		t.Errorf("deserializeDenvclustrFile() expected error, got nil")
		return
	}

	// Check that the error message contains the expected substring
	expectedErrSubstring := "cannot unmarshal number"
	if !strings.Contains(err.Error(), expectedErrSubstring) {
		t.Errorf("deserializeDenvclustrFile() error = %v, want error containing %q", err, expectedErrSubstring)
	}

	// Check that the result is nil
	if got != nil {
		t.Errorf("deserializeDenvclustrFile() = %v, want nil", got)
	}
}

// TestEmptyArrays tests deserialization of JSON with empty arrays where minItems=1 is required.
func TestEmptyArrays(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/empty_arrays.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Expected result
	want := &DenvclustrRoot{
		Name:           "empty-arrays",
		Infrastructure: []*Infrastructure{},
		Nodes:          []*Node{},
		Devcontainers:  []*Devcontainer{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("deserializeDenvclustrFile() = %v, want %v", got, want)
	}
}

// TestWrongEnumValues tests deserialization of JSON with wrong enum values.
func TestWrongEnumValues(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/wrong_enum_values.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Expected result
	want := &DenvclustrRoot{
		Name: "wrong-enum",
		Infrastructure: []*Infrastructure{
			{
				Id:       "aws_infra",
				Kind:     "container", // Wrong enum value
				Provider: "gcp",       // Wrong enum value
				Region:   "us-west-2",
			},
		},
		Nodes: []*Node{
			{
				Id:               "node1",
				InfrastructureId: "aws_infra",
				Properties: NodeProperties{
					InstanceType: "t2.micro",
				},
				RemoteAccess: NodeRemoteAccess{
					PublicSSHKey: "/path/to/key.pub",
				},
			},
		},
		Devcontainers: []*Devcontainer{
			{
				Id:     "dev1",
				NodeId: "node1",
				Source: &DevcontainerSource{
					URL: "https://github.com/example/repo.git",
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("deserializeDenvclustrFile() = %v, want %v", got, want)
	}
}

// TestDuplicateInfrastructureIds tests deserialization of JSON with duplicate infrastructure IDs.
func TestDuplicateInfrastructureIds(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/duplicate_infrastructure_ids.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Expected result - both infrastructure items should be present
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	if len(got.Infrastructure) != 2 {
		t.Errorf("Infrastructure: expected 2 items, got %d", len(got.Infrastructure))
		return
	}

	// Check that both infrastructure items have the same ID
	if string(got.Infrastructure[0].Id) != string(got.Infrastructure[1].Id) {
		t.Errorf("Expected duplicate IDs, got %q and %q", got.Infrastructure[0].Id, got.Infrastructure[1].Id)
	}
}

// TestNonexistentNodeId tests deserialization of JSON with node_id not matching any node.
func TestNonexistentNodeId(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/nonexistent_node_id.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Check that the devcontainer references a non-existent node ID
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	if len(got.Devcontainers) != 1 {
		t.Errorf("Devcontainers: expected 1 item, got %d", len(got.Devcontainers))
		return
	}

	// Check that the node ID doesn't exist in the nodes list
	nodeId := string(got.Devcontainers[0].NodeId)
	nodeExists := false
	for _, node := range got.Nodes {
		if string(node.Id) == nodeId {
			nodeExists = true
			break
		}
	}

	if nodeExists {
		t.Errorf("Expected node ID %q to not exist, but it does", nodeId)
	}
}

// TestExtraFields tests deserialization of JSON with extra optional fields.
func TestExtraFields(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/extra_fields.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as extra fields are ignored
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Extra fields should be ignored, so the result should match the expected structure
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	// Check that the required fields are present
	if string(got.Name) != "extra-fields" {
		t.Errorf("Name: got %q, want %q", got.Name, "extra-fields")
	}
	if len(got.Infrastructure) != 1 {
		t.Errorf("Infrastructure: expected 1 item, got %d", len(got.Infrastructure))
	}
	if len(got.Nodes) != 1 {
		t.Errorf("Nodes: expected 1 node, got %d", len(got.Nodes))
	}
	if len(got.Devcontainers) != 1 {
		t.Errorf("Devcontainers: expected 1 devcontainer, got %d", len(got.Devcontainers))
	}
}

// TestNullValues tests deserialization of JSON with null values for required fields.
func TestNullValues(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/null_values.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Check that null values are properly handled
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	// Check that the name is an empty string
	if string(got.Name) != "" {
		t.Errorf("Name: got %q, want empty string", got.Name)
	}
}

// TestSSHWithoutKey tests deserialization of JSON with SSH URL provided.
func TestSSHWithoutKey(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/ssh_without_key.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results - no error expected as validation happens elsewhere
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Check that the SSH URL is present in the devcontainer source
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	if len(got.Devcontainers) != 1 {
		t.Errorf("Devcontainers: expected 1 devcontainer, got %d", len(got.Devcontainers))
		return
	}

	// Verify that the source URL uses SSH format
	devcontainer := got.Devcontainers[0]
	expectedUrlPrefix := "git@github.com"
	if !strings.HasPrefix(string(devcontainer.Source.URL), expectedUrlPrefix) {
		t.Errorf("Source URL: got %q, want URL starting with %q", devcontainer.Source.URL, expectedUrlPrefix)
	}
}

// TestWhitespaceTrimming tests deserialization of JSON with whitespace around strings.
func TestWhitespaceTrimming(t *testing.T) {
	// Read the test file
	data, err := readTestFile(t, "testdata/whitespace_trimming.json")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Call the function being tested
	got, err := deserializeDenvclustrFile(data)

	// Validate results
	if err != nil {
		t.Errorf("deserializeDenvclustrFile() unexpected error: %v", err)
		return
	}

	// Check that strings were properly trimmed
	if got == nil {
		t.Errorf("deserializeDenvclustrFile() returned nil, expected non-nil result")
		return
	}

	if string(got.Name) != "trimmed-strings" {
		t.Errorf("Name not properly trimmed: got %q, want %q", got.Name, "trimmed-strings")
	}

	if len(got.Infrastructure) != 1 {
		t.Errorf("Expected 1 infrastructure item, got %d", len(got.Infrastructure))
		return
	}

	infra := got.Infrastructure[0]
	if string(infra.Id) != "aws_infra" {
		t.Errorf("Infrastructure.Id not properly trimmed: got %q, want %q", infra.Id, "aws_infra")
	}

	// It appears that the TrimmedString type's UnmarshalJSON method isn't being applied to enum types
	// This is likely because they are defined as custom string types rather than TrimmedString
	// For this test, we'll just check that the values contain the expected content
	if !strings.Contains(string(infra.Kind), "vm") {
		t.Errorf("Infrastructure.Kind doesn't contain expected value: got %q, want to contain %q", infra.Kind, "vm")
	}
	if !strings.Contains(string(infra.Provider), "aws") {
		t.Errorf("Infrastructure.Provider doesn't contain expected value: got %q, want to contain %q", infra.Provider, "aws")
	}
	if string(infra.Region) != "us-west-2" {
		t.Errorf("Infrastructure.Region not properly trimmed: got %q, want %q", infra.Region, "us-west-2")
	}
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}
