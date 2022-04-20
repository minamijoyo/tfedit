package migration

func Generate(planJSON []byte) ([]byte, error) {
	plan, err := NewPlan(planJSON)
	if err != nil {
		return nil, err
	}

	analyzer := NewDefaultAnalyzer()
	migration := analyzer.Analyze(plan)

	return migration.Render()
}
