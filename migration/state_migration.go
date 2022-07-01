package migration

import (
	"bytes"
	"fmt"
	"text/template"
)

// StateMigration is a type which corresponds to tfmigrate.StateMigratorConfig
// and config.MigrationBlock in minamijoyo/tfmigrate.
// The current implementation doesn't encode migration actions to a file
// directly with gohcl, so we define only what we need here.
type StateMigration struct {
	// A name label of migration block
	Name string
	// A working directory for executing terraform command.
	Dir string
	// A list of state action.
	Actions []StateAction
}

var migrationTemplate = `migration "state" "{{ .Name }}" {
{{- if ne .Dir "" }}
  dir = "{{ .Dir }}"
{{- end }}
  actions = [
  {{- range .Actions }}
    "{{ .MigrationAction }}",
  {{- end }}
  ]
}
`

var compiledMigrationTemplate = template.Must(template.New("migration").Parse(migrationTemplate))

// NewStateMigration returns a new instance of StateMigration.
func NewStateMigration(name string, dir string) *StateMigration {
	return &StateMigration{
		Name: name,
		Dir:  dir,
	}
}

// AppendActions appends a list of actions to migration.
func (m *StateMigration) AppendActions(actions ...StateAction) {
	m.Actions = append(m.Actions, actions...)
}

// Render converts a state migration config to bytes.
// Return an empty slice when no action without error.
// Encoding StateMigratorConfig directly with gohcl has some problems.
// An array contains multiple elements is output as one line. It's not readable
// for multiple actions. In additon, the default value is set explicitly, it's
// not only redundant but also increases cognitive load for user who isn't
// familiar with tfmigrate.
// So we use text/template to render a migration file.
func (m *StateMigration) Render() ([]byte, error) {
	// Return an empty slice when no action without error.
	if len(m.Actions) == 0 {
		return []byte{}, nil
	}

	var output bytes.Buffer
	if err := compiledMigrationTemplate.Execute(&output, m); err != nil {
		return nil, fmt.Errorf("failed to render migration file: %s", err)
	}

	return output.Bytes(), nil
}
