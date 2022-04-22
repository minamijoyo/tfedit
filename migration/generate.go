package migration

// Generate returns bytes of a migration file which reverts a given planned changes.
// The dir is set to a dir attribute in a migration file.
func Generate(planJSON []byte, dir string) ([]byte, error) {
	plan, err := NewPlan(planJSON)
	if err != nil {
		return nil, err
	}

	analyzer := NewDefaultPlanAnalyzer()
	migration := analyzer.Analyze(plan, dir)

	return migration.Render()
}
