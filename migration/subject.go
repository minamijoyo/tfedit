package migration

// Subject is a problem to be solved. It contains multiple conflicts.
type Subject struct {
	// A list of conflicts to be solved.
	conflicts []*Conflict
}

// NewSubject finds conflicts contained in a given plan and defines a problem.
func NewSubject(plan *Plan) *Subject {
	conflicts := []*Conflict{}
	for _, rc := range plan.ResourceChanges() {
		c := NewConflict(rc)
		conflicts = append(conflicts, c)
	}

	return &Subject{
		conflicts: conflicts,
	}
}

// Conflicts returns a list of conflicts. It may include already resolved.
func (s *Subject) Conflicts() []*Conflict {
	return s.conflicts
}

// IsResolved returns true if all conflicts have been resolved, otherwise false.
func (s *Subject) IsResolved() bool {
	for _, c := range s.conflicts {
		if !c.IsResolved() {
			return false
		}
	}

	return true
}
