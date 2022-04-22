package migration

// Generate returns bytes of a migration file which reverts a given planned changes.
func Generate(planJSON []byte) ([]byte, error) {
	plan, err := NewPlan(planJSON)
	if err != nil {
		return nil, err
	}

	analyzer := NewDefaultPlanAnalyzer()
	migration := analyzer.Analyze(plan)

	return migration.Render()
}
