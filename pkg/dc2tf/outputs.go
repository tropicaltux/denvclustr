package dc2tf

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func (c *converter) addOutputs(body *hclwrite.Body) error {
	for _, node := range c.root.Nodes {
		name := fmt.Sprintf("%s_output", node.Id)
		outputBody := body.AppendNewBlock("output", []string{name}).Body()

		// Create an object with a single key "module" and no quotes around the reference
		objTokens := hclwrite.Tokens{}

		// Opening brace
		objTokens = append(objTokens, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrace,
			Bytes: []byte("{"),
		})

		// Key "module"
		objTokens = append(objTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenIdent,
			Bytes:        []byte("module"),
			SpacesBefore: 1,
		})

		// Equal sign
		objTokens = append(objTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenEqual,
			Bytes:        []byte("="),
			SpacesBefore: 1,
		})

		// Reference to module without quotes
		objTokens = append(objTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenIdent,
			Bytes:        []byte("module"),
			SpacesBefore: 1,
		})
		objTokens = append(objTokens, &hclwrite.Token{
			Type:  hclsyntax.TokenDot,
			Bytes: []byte("."),
		})
		objTokens = append(objTokens, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(string(node.Id)),
		})

		// Closing brace
		objTokens = append(objTokens, &hclwrite.Token{
			Type:         hclsyntax.TokenCBrace,
			Bytes:        []byte("}"),
			SpacesBefore: 1,
		})

		outputBody.SetAttributeRaw("value", objTokens)
		outputBody.SetAttributeValue("sensitive", cty.BoolVal(true))
	}
	return nil
}
