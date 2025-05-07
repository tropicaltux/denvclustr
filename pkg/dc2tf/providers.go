package dc2tf

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/tropicaltux/denvclustr/pkg/schema"
	"github.com/zclconf/go-cty/cty"
)

func (c *converter) addProviders(body *hclwrite.Body) error {
	for _, infrastructure := range c.root.Infrastructure {
		if infrastructure.Provider != schema.ProviderAws {
			return fmt.Errorf("unsupported provider %q", infrastructure.Provider)
		}
		providerBody := body.AppendNewBlock("provider", []string{"aws"}).Body()
		providerBody.SetAttributeValue("region", cty.StringVal(string(infrastructure.Region)))
		providerBody.SetAttributeValue("alias", cty.StringVal(string(infrastructure.Id)))
	}
	return nil
}
