package dc2tf

import (
	"embed"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/stretchr/testify/require"
	"github.com/tmccombs/hcl2json/convert"

	"github.com/tropicaltux/denvclustr/pkg/schema"
)

//go:embed testdata/*.tf
var testdataFS embed.FS

// compareHCL compares two parsed HCL files
func compareHCL(t testing.TB, expect, actual *hcl.File) error {
	expectJsonBytes, err := convert.File(expect, convert.Options{Simplify: true})
	require.NoError(t, err)
	var expectJson map[string]interface{}
	err = json.Unmarshal(expectJsonBytes, &expectJson)
	require.NoError(t, err)

	actualJsonBytes, err := convert.File(actual, convert.Options{Simplify: true})
	require.NoError(t, err)
	var actualJson map[string]interface{}
	err = json.Unmarshal(actualJsonBytes, &actualJson)
	require.NoError(t, err)

	// Compare structures
	require.True(t, reflect.DeepEqual(expectJson, actualJson))

	return nil
}

func TestConvert(t *testing.T) {
	minimal, err := testdataFS.ReadFile("testdata/valid_minimal.tf")
	require.NoError(t, err)
	withDNS, err := testdataFS.ReadFile("testdata/with_dns.tf")
	require.NoError(t, err)
	multipleNodes, err := testdataFS.ReadFile("testdata/multiple_nodes.tf")
	require.NoError(t, err)
	multipleInfrastructure, err := testdataFS.ReadFile("testdata/multiple_infrastructures.tf")
	require.NoError(t, err)
	complexConfig, err := testdataFS.ReadFile("testdata/complex_config.tf")
	require.NoError(t, err)
	emptyDevcontainers, err := testdataFS.ReadFile("testdata/empty_devcontainers.tf")
	require.NoError(t, err)

	// Parse expected HCL files
	parser := hclparse.NewParser()
	expectedMinimal, diags1 := parser.ParseHCL(minimal, "expected_minimal.tf")
	require.False(t, diags1.HasErrors(), "failed parsing expected minimal: %v", diags1)
	expectedWithDNS, diags2 := parser.ParseHCL(withDNS, "expected_with_dns.tf")
	require.False(t, diags2.HasErrors(), "failed parsing expected with DNS: %v", diags2)
	expectedMultipleNodes, diags3 := parser.ParseHCL(multipleNodes, "expected_multiple_nodes.tf")
	require.False(t, diags3.HasErrors(), "failed parsing expected multiple nodes: %v", diags3)
	expectedMultipleInfrastructure, diags4 := parser.ParseHCL(multipleInfrastructure, "expected_multiple_infrastructures.tf")
	require.False(t, diags4.HasErrors(), "failed parsing expected multiple infrastructure: %v", diags4)
	expectedComplexConfig, diags5 := parser.ParseHCL(complexConfig, "expected_complex_config.tf")
	require.False(t, diags5.HasErrors(), "failed parsing expected complex config: %v", diags5)
	expectedEmptyDevcontainers, diags6 := parser.ParseHCL(emptyDevcontainers, "expected_empty_devcontainers.tf")
	require.False(t, diags6.HasErrors(), "failed parsing expected empty devcontainers: %v", diags6)

	cases := []struct {
		name     string
		spec     *schema.DenvclustrRoot
		expected *hcl.File
	}{
		{
			"minimal config",
			&schema.DenvclustrRoot{
				Name: schema.TrimmedString("test-cluster"),
				Infrastructure: []*schema.Infrastructure{{
					Id:       schema.TrimmedString("infrastructure1"),
					Provider: schema.ProviderAws,
					Kind:     schema.KindVm,
					Region:   schema.TrimmedString("us-west-2"),
				}},
				Nodes: []*schema.Node{{
					Id:               schema.TrimmedString("node1"),
					InfrastructureId: schema.TrimmedString("infrastructure1"),
					Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.micro")},
					RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
				}},
				Devcontainers: []*schema.Devcontainer{{
					Id:           schema.TrimmedString("dev1"),
					NodeId:       schema.TrimmedString("node1"),
					Source:       &schema.DevcontainerSource{URL: schema.TrimmedString("https://github.com/example/repo")},
					RemoteAccess: &schema.DevcontainerRemoteAccess{},
				}},
			}, expectedMinimal},
		{"config with DNS", &schema.DenvclustrRoot{
			Name: schema.TrimmedString("test-cluster"),
			Infrastructure: []*schema.Infrastructure{{
				Id:       schema.TrimmedString("infrastructure1"),
				Provider: schema.ProviderAws,
				Kind:     schema.KindVm,
				Region:   schema.TrimmedString("us-west-2"),
			}},
			Nodes: []*schema.Node{{
				Id:               schema.TrimmedString("node1"),
				InfrastructureId: schema.TrimmedString("infrastructure1"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.micro")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
				DNS:              &schema.NodeDNS{HighLevelDomain: schema.TrimmedString("example.com")},
			}},
			Devcontainers: []*schema.Devcontainer{{
				Id:           schema.TrimmedString("dev1"),
				NodeId:       schema.TrimmedString("node1"),
				Source:       &schema.DevcontainerSource{URL: schema.TrimmedString("https://github.com/example/repo")},
				RemoteAccess: &schema.DevcontainerRemoteAccess{},
			}},
		}, expectedWithDNS},
		{"multiple nodes", &schema.DenvclustrRoot{
			Name: schema.TrimmedString("test-cluster"),
			Infrastructure: []*schema.Infrastructure{{
				Id:       schema.TrimmedString("infrastructure1"),
				Provider: schema.ProviderAws,
				Kind:     schema.KindVm,
				Region:   schema.TrimmedString("us-west-2"),
			}},
			Nodes: []*schema.Node{{
				Id:               schema.TrimmedString("node1"),
				InfrastructureId: schema.TrimmedString("infrastructure1"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.micro")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
			}, {
				Id:               schema.TrimmedString("node2"),
				InfrastructureId: schema.TrimmedString("infrastructure1"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.large")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
			}},
			Devcontainers: []*schema.Devcontainer{{
				Id:           schema.TrimmedString("dev1"),
				NodeId:       schema.TrimmedString("node1"),
				Source:       &schema.DevcontainerSource{URL: schema.TrimmedString("https://github.com/example/repo")},
				RemoteAccess: &schema.DevcontainerRemoteAccess{},
			}, {
				Id:     schema.TrimmedString("dev2"),
				NodeId: schema.TrimmedString("node2"),
				Source: &schema.DevcontainerSource{
					URL:    schema.TrimmedString("https://github.com/example/app"),
					Branch: schema.TrimmedString("main"),
				},
				RemoteAccess: &schema.DevcontainerRemoteAccess{},
			}},
		}, expectedMultipleNodes},
		{"multiple infrastructure", &schema.DenvclustrRoot{
			Name: schema.TrimmedString("test-cluster"),
			Infrastructure: []*schema.Infrastructure{{
				Id:       schema.TrimmedString("infrastructure1"),
				Provider: schema.ProviderAws,
				Kind:     schema.KindVm,
				Region:   schema.TrimmedString("us-west-2"),
			}, {
				Id:       schema.TrimmedString("infrastructure2"),
				Provider: schema.ProviderAws,
				Kind:     schema.KindVm,
				Region:   schema.TrimmedString("eu-west-1"),
			}},
			Nodes: []*schema.Node{{
				Id:               schema.TrimmedString("node1"),
				InfrastructureId: schema.TrimmedString("infrastructure1"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.micro")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
			}, {
				Id:               schema.TrimmedString("node2"),
				InfrastructureId: schema.TrimmedString("infrastructure2"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.small")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
			}},
			Devcontainers: []*schema.Devcontainer{{
				Id:           schema.TrimmedString("dev1"),
				NodeId:       schema.TrimmedString("node1"),
				Source:       &schema.DevcontainerSource{URL: schema.TrimmedString("https://github.com/example/repo")},
				RemoteAccess: &schema.DevcontainerRemoteAccess{},
			}, {
				Id:           schema.TrimmedString("dev2"),
				NodeId:       schema.TrimmedString("node2"),
				Source:       &schema.DevcontainerSource{URL: schema.TrimmedString("https://github.com/example/app")},
				RemoteAccess: &schema.DevcontainerRemoteAccess{},
			}},
		}, expectedMultipleInfrastructure},
		{"complex configuration", &schema.DenvclustrRoot{
			Name: schema.TrimmedString("test-cluster"),
			Infrastructure: []*schema.Infrastructure{{
				Id:       schema.TrimmedString("infrastructure1"),
				Provider: schema.ProviderAws,
				Kind:     schema.KindVm,
				Region:   schema.TrimmedString("us-west-2"),
			}},
			Nodes: []*schema.Node{{
				Id:               schema.TrimmedString("node1"),
				InfrastructureId: schema.TrimmedString("infrastructure1"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.micro")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
				DNS:              &schema.NodeDNS{HighLevelDomain: schema.TrimmedString("example.com")},
			}},
			Devcontainers: []*schema.Devcontainer{{
				Id:     schema.TrimmedString("dev1"),
				NodeId: schema.TrimmedString("node1"),
				Source: &schema.DevcontainerSource{
					URL:              schema.TrimmedString("https://github.com/example/repo"),
					Branch:           schema.TrimmedString("feature-branch"),
					DevcontainerPath: schema.TrimmedString(".devcontainer"),
					SshKey: &schema.DevcontainerSourceSSHKey{
						Reference: schema.TrimmedString("github-key"),
						Source:    "secrets_manager",
					},
				},
				RemoteAccess: &schema.DevcontainerRemoteAccess{
					OpenVsCodeServer: &schema.DevcontainerOpenVSCodeServer{
						Port: func() *int { port := 3000; return &port }(),
					},
					Ssh: &schema.DevcontainerSSH{
						Port:         func() *int { port := 2222; return &port }(),
						PublicSshKey: schema.TrimmedString("~/.ssh/custom_key.pub"),
					},
				},
			}},
		}, expectedComplexConfig},
		{"empty devcontainers", &schema.DenvclustrRoot{
			Name: schema.TrimmedString("test-cluster"),
			Infrastructure: []*schema.Infrastructure{{
				Id:       schema.TrimmedString("infrastructure1"),
				Provider: schema.ProviderAws,
				Kind:     schema.KindVm,
				Region:   schema.TrimmedString("us-west-2"),
			}},
			Nodes: []*schema.Node{{
				Id:               schema.TrimmedString("node1"),
				InfrastructureId: schema.TrimmedString("infrastructure1"),
				Properties:       schema.NodeProperties{InstanceType: schema.TrimmedString("t3.micro")},
				RemoteAccess:     schema.NodeRemoteAccess{PublicSSHKey: schema.TrimmedString("~/.ssh/id_rsa.pub")},
			}},
			Devcontainers: []*schema.Devcontainer{},
		}, expectedEmptyDevcontainers},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := Convert(c.spec)
			require.NoError(t, err)

			// Parse the result into HCL format
			parser := hclparse.NewParser()
			actual, diags := parser.ParseHCL(out.Bytes(), "actual.tf")
			require.False(t, diags.HasErrors(), "failed parsing actual: %v", diags)

			// Compare with expected result using compareHCL
			compareHCL(t, c.expected, actual)
		})
	}

	// Error test cases
	t.Run("unsupported provider", func(t *testing.T) {
		_, err := Convert(&schema.DenvclustrRoot{
			Name:           schema.TrimmedString("test-cluster"),
			Infrastructure: []*schema.Infrastructure{{Provider: "azure"}},
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "unsupported provider")
	})

	t.Run("nil input", func(t *testing.T) {
		_, err := Convert(nil)
		require.Error(t, err)
	})
}
