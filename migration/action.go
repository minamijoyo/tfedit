package migration

import "fmt"

// Action is an interface of terraform command for state migration.
type Action interface {
	// MigrationAction returns a string of terraform command for state migration.
	MigrationAction() string
}

// StateImportAction implements the Action interface.
type StateImportAction struct {
	address string
	id      string
}

var _ Action = (*StateImportAction)(nil)

// NewStateImportAction returns a new instance of StateImportAction.
func NewStateImportAction(address string, id string) Action {
	return &StateImportAction{
		address: address,
		id:      id,
	}
}

// MigrationAction returns a string of terraform command for state migration.
func (a *StateImportAction) MigrationAction() string {
	return fmt.Sprintf("import %s %s", a.address, a.id)
}
