package migration

import (
	"github.com/minamijoyo/tfedit/migration/schema"
	_ "github.com/minamijoyo/tfedit/migration/schema/aws" // Register schema for aws
)

// Resolver is an interface that abstracts a rule for solving a subject.
type Resolver interface {
	// Resolve tries to resolve some conflicts in a given subject and returns the
	// updated subject and state migration actions.
	Resolve(s *Subject) (*Subject, []StateAction)
}

// StateImportResolver is an implementation of Resolver for import.
type StateImportResolver struct {
}

var _ Resolver = (*StateImportResolver)(nil)

// Resolve tries to resolve some conflicts in a given subject and returns the
// updated subject and state migration actions.
// It translates a planned create action into an import state migration.
func (r *StateImportResolver) Resolve(s *Subject) (*Subject, []StateAction) {
	actions := []StateAction{}
	for _, c := range s.Conflicts() {
		if c.IsResolved() {
			continue
		}

		switch c.PlannedActionType() {
		case "create":
			importID := schema.ImportID(c.ResourceType(), c.ResourceAfter())
			action := NewStateImportAction(c.Address(), importID)
			actions = append(actions, action)
			c.MarkAsResolved()
		}
	}

	return s, actions
}
