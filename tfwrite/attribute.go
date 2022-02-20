package tfwrite

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/minamijoyo/hcledit/editor"
)

// Attribute is an attribute of resource.
type Attribute struct {
	raw *hclwrite.Attribute
}

// NewAttribute creates a new instance of Attribute.
func NewAttribute(attr *hclwrite.Attribute) *Attribute {
	return &Attribute{raw: attr}
}

// ValueAsString returns a value of Attribute as string.
func (a *Attribute) ValueAsString() (string, error) {
	return editor.GetAttributeValueAsString(a.raw)
}

// ValueAsTokens returns a value of Attribute as raw tokens.
func (a *Attribute) ValueAsTokens() hclwrite.Tokens {
	return a.raw.Expr().BuildTokens(nil)
}
