package migration

// Resolver is an interface that abstracts a rule for solving a subject.
type Resolver interface {
	// Resolve tries to resolve some conflicts in a given subject and returns the
	// updated subject and state migration actions.
	Resolve(s *Subject) (*Subject, []StateAction, error)
}
