package migration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mitchellh/mapstructure"
)

var migrationTemplate string = `migration "state" "awsv4upgrade" {
  actions = [
  {{- range . }}
    "{{ . }}",
  {{- end }}
  ]
}
`

type AWSS3BucketACLFilterResource struct {
	Bucket string
	ACL    string
}

func (r *AWSS3BucketACLFilterResource) ID() string {
	return r.Bucket + "," + r.ACL
}

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
			var r AWSS3BucketACLFilterResource
			mapstructure.Decode(after, &r)
			migrateAction := fmt.Sprintf("import %s %s", address, r.ID())
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
