package migration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/minamijoyo/tfedit/migration/schema"
	_ "github.com/minamijoyo/tfedit/migration/schema/aws"
)

var migrationTemplate string = `migration "state" "awsv4upgrade" {
  actions = [
  {{- range . }}
    "{{ . }}",
  {{- end }}
  ]
}
`

func Generate(planJSON []byte) ([]byte, error) {
	var plan tfjson.Plan
	if err := json.Unmarshal(planJSON, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan file: %s", err)
	}

	var migrateActions []string
	for _, rc := range plan.ResourceChanges {
		if rc.Change.Actions.Create() {
			address := rc.Address
			after := rc.Change.After.(map[string]interface{})
			importID := schema.ImportID(rc.Type, after)
			migrateAction := fmt.Sprintf("import %s %s", address, importID)
			migrateActions = append(migrateActions, migrateAction)
		}
	}

	tpl := template.Must(template.New("migration").Parse(migrationTemplate))
	var output bytes.Buffer
	if err := tpl.Execute(&output, migrateActions); err != nil {
		return nil, fmt.Errorf("failed to render migration file: %s", err)
	}

	return output.Bytes(), nil
}
