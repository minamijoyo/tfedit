package migration

import (
	"bytes"
	"fmt"
	"text/template"
)

// File is a type which is equivalent to tfmigrate.StateMigratorConfig of
// minamijoyo/tfmigrate.
// The current implementation doesn't encode migration actions to a file
// directly with gohcl, so we avoid to depend on tfmigrate's type and we define
// only what we need here.
type File struct {
	// Dir is a working directory for executing terraform command.
	Dir string
	// Actions is a list of state action.
	Actions []Action
}

var migrationTemplate = `migration "state" "awsv4upgrade" {
  actions = [
  {{- range . }}
    "{{ .MigrationAction }}",
  {{- end }}
  ]
}
`

var compiledMigrationTemplate = template.Must(template.New("migration").Parse(migrationTemplate))

// AppendAction appends an action to migration.
func (f *File) AppendAction(action Action) {
	f.Actions = append(f.Actions, action)
}

// Render converts a state migration config to bytes
// Encoding StateMigratorConfig directly with gohcl has some problems.
// An array contains multiple elements is output as one line. It's not readable
// for multiple actions. In additon, the default value is set explicitly, it's
// not only redundant but also increases cognitive load for user who isn't
// familiar with tfmigrate.
// So we use text/template to render a migration file.
func (f *File) Render() ([]byte, error) {
	var output bytes.Buffer
	if err := compiledMigrationTemplate.Execute(&output, f.Actions); err != nil {
		return nil, fmt.Errorf("failed to render migration file: %s", err)
	}

	return output.Bytes(), nil
}
