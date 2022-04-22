package migration

// PlanAnalyzer is an interface that abstracts the analysis rules of plan.
type PlanAnalyzer interface {
	// Analyze analyzes a given plan and generates a state migration so that
	// the plan results in no changes.
	Analyze(plan *Plan) *StateMigration
}

// defaultPlanAnalyzer is a default implementation for PlanAnalyzer.
// This is a predefined rule-based analyzer.
type defaultPlanAnalyzer struct {
	// A list of rules used for analysis.
	resolvers []Resolver
}

var _ PlanAnalyzer = (*defaultPlanAnalyzer)(nil)

// NewDefaultPlanAnalyzer returns a new instance of defaultPlanAnalyzer.
// The current implementation only supports import, but allows us to compose
// multiple resolvers for future extension.
func NewDefaultPlanAnalyzer() PlanAnalyzer {
	return &defaultPlanAnalyzer{
		resolvers: []Resolver{
			&StateImportResolver{},
		},
	}
}

// Analyze analyzes a given plan and generates a state migration so that
// the plan results in no changes.
func (a *defaultPlanAnalyzer) Analyze(plan *Plan) *StateMigration {
	subject := NewSubject(plan)

	var migration StateMigration
	current := subject
	for _, r := range a.resolvers {
		next, actions := r.Resolve(current)
		migration.AppendActions(actions...)
		current = next
	}

	return &migration
}
