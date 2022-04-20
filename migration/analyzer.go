package migration

type Analyzer interface {
	Analyze(plan *Plan) *StateMigration
}

type defaultAnalyzer struct {
	resolvers []Resolver
}

var _ Analyzer = (*defaultAnalyzer)(nil)

func NewDefaultAnalyzer() Analyzer {
	return &defaultAnalyzer{
		resolvers: []Resolver{
			&StateImportResolver{},
		},
	}
}

func (a *defaultAnalyzer) Analyze(plan *Plan) *StateMigration {
	var migration StateMigration

	for _, r := range a.resolvers {
		actions := r.Resolve(plan)
		migration.AppendActions(actions...)
	}

	return &migration
}
