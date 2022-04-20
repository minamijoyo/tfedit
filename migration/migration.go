package migration

import (
	"github.com/minamijoyo/tfedit/migration/schema"
	_ "github.com/minamijoyo/tfedit/migration/schema/aws" // Register schema for aws
)

func Generate(planJSON []byte) ([]byte, error) {
	plan, err := NewPlan(planJSON)
	if err != nil {
		return nil, err
	}

	var migration StateMigration
	for _, rc := range plan.ResourceChanges() {
		if rc.Change.Actions.Create() {
			address := rc.Address
			after := rc.Change.After.(map[string]interface{})
			importID := schema.ImportID(rc.Type, after)
			action := NewStateImportAction(address, importID)
			migration.AppendAction(action)
		}
	}

	return migration.Render()
}
