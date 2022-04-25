package migration

import (
	"github.com/minamijoyo/tfedit/migration/schema"
	"github.com/minamijoyo/tfedit/migration/schema/aws"
)

// Generate returns bytes of a migration file which reverts a given planned changes.
// The dir is set to a dir attribute in a migration file.
func Generate(planJSON []byte, dir string) ([]byte, error) {
	plan, err := NewPlan(planJSON)
	if err != nil {
		return nil, err
	}

	dictionary := NewDefaultDictionary()
	analyzer := NewDefaultPlanAnalyzer(dictionary)
	migration := analyzer.Analyze(plan, dir)

	return migration.Render()
}

// NewDefaultDictionary returns a default built-in Dictionary.
func NewDefaultDictionary() *schema.Dictionary {
	d := schema.NewDictionary()
	aws.RegisterSchema(d)
	return d
}
