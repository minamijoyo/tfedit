package migration

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// StateAction is an interface of action for state migration.
type StateAction interface {
	// MigrationAction returns a string of action for state migration.
	// It escapes special characters in HCL for use as an action in a tfmigrate's
	// migration file.
	MigrationAction() string
}

// actionEscape is a helper function which escapes special characters in HCL for
// use as an action in a tfmigrate's migration file.
func actionEscape(raw string) string {
	// Since the hclwrite.escapeQuotedStringLit() is unexported,
	// implement the HCL escaping relying on the fact that hclwrite tokens are
	// implicitly escaped when generated.
	tokens := hclwrite.TokensForValue(cty.StringVal(raw))

	// The TokensForValue() wraps tokens with double quotes.
	// Remove `"` TokenOQuote at head and `"` TokenCQuote at tail
	unquoted := tokens[1 : len(tokens)-1]
	escaped := string(unquoted.Bytes())

	// If escaping was required, enclose it in single quotes
	// so that a shell does not interpret double quotes.
	// It is needed to use the result as an action in a tfmigrate's migration file.
	if raw != escaped {
		return "'" + escaped + "'"
	}
	return raw
}

// StateImportAction implements the StateAction interface.
type StateImportAction struct {
	address string
	id      string
}

var _ StateAction = (*StateImportAction)(nil)

// NewStateImportAction returns a new instance of StateImportAction.
func NewStateImportAction(address string, id string) StateAction {
	return &StateImportAction{
		address: address,
		id:      id,
	}
}

// MigrationAction returns a string of action for state migration.
// It escapes special characters in HCL for use as an action in a tfmigrate's
// migration file.
func (a *StateImportAction) MigrationAction() string {
	return fmt.Sprintf("import %s %s", actionEscape(a.address), actionEscape(a.id))
}
