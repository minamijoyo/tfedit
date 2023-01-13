package tfwrite

import (
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

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

// References returns all variable references contained in the value.
// It returns a unique and sorted list.
func (a *Attribute) References() []string {
	refs := map[string]struct{}{}
	for _, key := range a.rawReferences() {
		// To remove duplicates, append keys to a map.
		refs[key] = struct{}{}
	}

	keys := maps.Keys(refs)
	slices.Sort(keys)
	return keys
}

// rawReferences returns a list of raw variable references.
// The result may be unsorted and duplicated,
// but it is efficient to merge the result later.
func (a *Attribute) rawReferences() []string {
	refs := []string{}
	for _, v := range a.raw.Expr().Variables() {
		// Covert *hclwrite.Traversal to string.
		// This is a bit dirty, but it’s the only way to do at the time of writing.
		key := strings.TrimSpace(string(v.BuildTokens(nil).Bytes()))

		refs = append(refs, key)
	}

	return refs
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
