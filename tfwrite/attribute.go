package tfwrite

import (
	"strings"

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

// RenameReference renames all variable references contained in the value.
// The `from` and `to` arguments are specified as dot-delimited resource addresses.
// The `from` argument can be a partial prefix match, but must match the length
// of the `to` argument.
func (a *Attribute) RenameReference(from string, to string) {
	search := strings.Split(from, ".")
	replacement := strings.Split(to, ".")
	a.raw.Expr().RenameVariablePrefix(search, replacement)
}
