package migration

import (
	"github.com/minamijoyo/tfedit/migration/schema"
	_ "github.com/minamijoyo/tfedit/migration/schema/aws" // Register schema for aws
)

type Resolver interface {
	Resolve(plan *Plan) []Action
}

type StateImportResolver struct {
}

var _ Resolver = (*StateImportResolver)(nil)

func (r *StateImportResolver) Resolve(plan *Plan) []Action {
	actions := []Action{}
	for _, rc := range plan.ResourceChanges() {
		if rc.Change.Actions.Create() {
			address := rc.Address
			after := rc.Change.After.(map[string]interface{})
			importID := schema.ImportID(rc.Type, after)
			action := NewStateImportAction(address, importID)
			actions = append(actions, action)
		}
	}

	return actions
}
