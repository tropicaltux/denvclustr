package dc2tf

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func (c *converter) addOutputs(body *hclwrite.Body) error {
	for _, node := range c.root.Nodes {
		name := fmt.Sprintf("%s_output", node.Id)
		outputBody := body.AppendNewBlock("output", []string{name}).Body()
		outputBody.SetAttributeValue("value", cty.ObjectVal(map[string]cty.Value{
			"module": cty.StringVal(fmt.Sprintf("module.%s", node.Id)),
		}))
		outputBody.SetAttributeValue("sensitive", cty.BoolVal(true))
	}
	return nil
}
