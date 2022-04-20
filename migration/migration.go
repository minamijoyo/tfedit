package migration

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/minamijoyo/tfedit/migration/schema"
	_ "github.com/minamijoyo/tfedit/migration/schema/aws" // Register schema for aws
)

func Generate(planJSON []byte) ([]byte, error) {
	var plan tfjson.Plan
	if err := json.Unmarshal(planJSON, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan file: %s", err)
	}

	var file File
	for _, rc := range plan.ResourceChanges {
		if rc.Change.Actions.Create() {
			address := rc.Address
			after := rc.Change.After.(map[string]interface{})
			importID := schema.ImportID(rc.Type, after)
			migrateAction := fmt.Sprintf("import %s %s", address, importID)
			file.AppendAction(migrateAction)
		}
	}

	return file.Render()
}
