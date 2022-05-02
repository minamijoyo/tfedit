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
		if !rc.Change.Actions.NoOp() {
			c := NewConflict(rc)
			conflicts = append(conflicts, c)
		}
	}

	return &Subject{
		conflicts: conflicts,
	}
}

// UnresolvedConflicts returns a list of unresolved conflicts.
func (s *Subject) UnresolvedConflicts() []*Conflict {
	ret := []*Conflict{}
	for _, c := range s.conflicts {
		if !c.IsResolved() {
			ret = append(ret, c)
		}
	}

	return ret
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
