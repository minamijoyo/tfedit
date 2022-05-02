package migration

import (
	"github.com/minamijoyo/tfedit/migration/schema"
	"github.com/minamijoyo/tfedit/migration/schema/aws"
)

// PlanAnalyzer is an interface that abstracts the analysis rules of plan.
type PlanAnalyzer interface {
	// Analyze analyzes a given plan and generates a state migration so that
	// the plan results in no changes.
	// The dir is set to a dir attribute in a migration file.
	Analyze(plan *Plan, dir string) (*StateMigration, error)
}

// defaultPlanAnalyzer is a default implementation for PlanAnalyzer.
// This is a predefined rule-based analyzer.
type defaultPlanAnalyzer struct {
	// A dictionary for provider schema.
	dictionary *schema.Dictionary
	// A list of rules used for analysis.
	resolvers []Resolver
}

var _ PlanAnalyzer = (*defaultPlanAnalyzer)(nil)

// NewDefaultPlanAnalyzer returns a new instance of defaultPlanAnalyzer.
// The current implementation only supports import, but allows us to compose
// multiple resolvers for future extension.
func NewDefaultPlanAnalyzer(d *schema.Dictionary) PlanAnalyzer {
	return &defaultPlanAnalyzer{
		dictionary: d,
		resolvers: []Resolver{
			NewStateImportResolver(d),
		},
	}
}

// Analyze analyzes a given plan and generates a state migration so that
// the plan results in no changes.
// The dir is set to a dir attribute in a migration file.
func (a *defaultPlanAnalyzer) Analyze(plan *Plan, dir string) (*StateMigration, error) {
	subject := NewSubject(plan)

	migration := NewStateMigration("fromplan", dir)
	current := subject
	for _, r := range a.resolvers {
		next, actions, err := r.Resolve(current)
		if err != nil {
			return nil, err
		}
		migration.AppendActions(actions...)
		current = next
	}

	return migration, nil
}

// GenerateFromPlan returns bytes of a migration file which reverts a given
// planned changes.
// The dir is set to a dir attribute in a migration file.
func GenerateFromPlan(planJSON []byte, dir string) ([]byte, error) {
	plan, err := NewPlan(planJSON)
	if err != nil {
		return nil, err
	}

	dictionary := NewDefaultDictionary()
	analyzer := NewDefaultPlanAnalyzer(dictionary)
	migration, err := analyzer.Analyze(plan, dir)
	if err != nil {
		return nil, err
	}

	return migration.Render()
}

// NewDefaultDictionary returns a default built-in Dictionary.
func NewDefaultDictionary() *schema.Dictionary {
	d := schema.NewDictionary()
	aws.RegisterSchema(d)
	return d
}
