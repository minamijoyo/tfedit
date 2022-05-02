package migration

import (
	"github.com/minamijoyo/tfedit/migration/schema"
)

// StateImportResolver is an implementation of Resolver for import.
type StateImportResolver struct {
	// A dictionary for provider schema.
	dictionary *schema.Dictionary
}

var _ Resolver = (*StateImportResolver)(nil)

// NewStateImportResolver returns a new instance of StateImportResolver.
func NewStateImportResolver(d *schema.Dictionary) Resolver {
	return &StateImportResolver{
		dictionary: d,
	}
}

// Resolve tries to resolve some conflicts in a given subject and returns the
// updated subject and state migration actions.
// It translates a planned create action into an import state migration.
func (r *StateImportResolver) Resolve(s *Subject) (*Subject, []StateAction, error) {
	actions := []StateAction{}
	for _, c := range s.UnresolvedConflicts() {
		switch c.PlannedActionType() {
		case "create":
			resource, err := c.ResourceAfter()
			if err != nil {
				return nil, nil, err
			}

			importID, err := r.dictionary.ImportID(c.ResourceType(), resource)
			if err != nil {
				return nil, nil, err
			}

			action := NewStateImportAction(c.Address(), importID)
			actions = append(actions, action)
			c.MarkAsResolved()
		}
	}

	return s, actions, nil
}
