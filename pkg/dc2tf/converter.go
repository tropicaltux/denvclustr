package dc2tf

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/tropicaltux/denvclustr/pkg/schema"
)

// Convert is the single public entry‑point that takes
// denvclustr configuration and returns a Terraform HCL file.
func Convert(spec *schema.DenvclustrRoot) (*hclwrite.File, error) {
	if spec == nil {
		return nil, fmt.Errorf("nil input: spec cannot be nil")
	}
	c := &converter{root: spec}
	return c.toTerraform()
}

type converter struct {
	root *schema.DenvclustrRoot
}

func (c *converter) toTerraform() (*hclwrite.File, error) {
	f := hclwrite.NewEmptyFile()
	root := f.Body()

	if err := c.addProviders(root); err != nil {
		return nil, err
	}
	if err := c.addModules(root); err != nil {
		return nil, err
	}
	if err := c.addOutputs(root); err != nil {
		return nil, err
	}
	return f, nil
}
