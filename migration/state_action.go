package migration

import "fmt"

// StateAction is an interface of terraform command for state migration.
type StateAction interface {
	// MigrationAction returns a string of terraform command for state migration.
	MigrationAction() string
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

// MigrationAction returns a string of terraform command for state migration.
func (a *StateImportAction) MigrationAction() string {
	return fmt.Sprintf("import %s %s", a.address, a.id)
}
