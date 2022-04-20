package migration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/minamijoyo/tfedit/migration/schema"
	_ "github.com/minamijoyo/tfedit/migration/schema/aws" // Register schema for aws
)

// Migration is a type which is equivalent to tfmigrate.StateMigratorConfig of
// minamijoyo/tfmigrate.
// The current implementation doesn't encode migration actions to a file
// directly with gohcl, so we avoid to depend on tfmigrate's type and we define
// only what we need here.
type Migration struct {
	// Dir is a working directory for executing terraform command.
	Dir string
	// Actions is a list of state action.
	Actions []string
}

var migrationTemplate = `migration "state" "awsv4upgrade" {
  actions = [
  {{- range . }}
    "{{ . }}",
  {{- end }}
  ]
}
`

var compiledMigrationTemplate = template.Must(template.New("migration").Parse(migrationTemplate))

// AppendAction appends an action to migration.
func (m *Migration) AppendAction(action string) {
	m.Actions = append(m.Actions, action)
}

// Encode converts a state migration config to bytes
func (m *Migration) Encode() ([]byte, error) {
	// Encoding StateMigratorConfig directly with gohcl has some problems.
	// An array contains multiple elements is output as one line. It's not readable
	// for multiple actions. In additon, the default value is set explicitly, it's
	// not only redundant but also increases cognitive load for user who isn't
	// familiar with tfmigrate.
	// So we use text/template to render a migration file.
	var output bytes.Buffer
	if err := compiledMigrationTemplate.Execute(&output, m.Actions); err != nil {
		return nil, fmt.Errorf("failed to render migration file: %s", err)
	}

	return output.Bytes(), nil
}

func Generate(planJSON []byte) ([]byte, error) {
	var plan tfjson.Plan
	if err := json.Unmarshal(planJSON, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan file: %s", err)
	}

	var migration Migration
	for _, rc := range plan.ResourceChanges {
		if rc.Change.Actions.Create() {
			address := rc.Address
			after := rc.Change.After.(map[string]interface{})
			importID := schema.ImportID(rc.Type, after)
			migrateAction := fmt.Sprintf("import %s %s", address, importID)
			migration.AppendAction(migrateAction)
		}
	}

	return migration.Encode()
}
