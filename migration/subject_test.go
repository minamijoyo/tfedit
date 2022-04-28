package migration

import (
	"os"
	"testing"
)

func TestNewSubject(t *testing.T) {
	cases := []struct {
		desc     string
		planFile string
		ok       bool
		resolved bool
		want     int
	}{
		{
			desc:     "simple",
			planFile: "test-fixtures/import_simple.tfplan.json",
			ok:       true,
			resolved: false,
			want:     1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			planJSON, err := os.ReadFile(tc.planFile)
			if err != nil {
				t.Fatalf("failed to read file: %s", err)
			}

			plan, err := NewPlan(planJSON)
			if err != nil {
				t.Fatalf("failed to new plan: %s", err)
			}

			s := NewSubject(plan)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %#v", s)
			}

			if tc.ok {
				if s.IsResolved() != tc.resolved {
					t.Errorf("unexpected the resolved status of subject. got = %t, but want = %t", s.IsResolved(), tc.resolved)
				}
				got := len(s.UnresolvedConflicts())
				if got != tc.want {
					t.Errorf("got = %d, but want = %d", got, tc.want)
				}
			}
		})
	}
}

func TestSubjectUnresolvedConflicts(t *testing.T) {
	cases := []struct {
		desc string
		s    *Subject
		want int
	}{
		{
			desc: "unresolved",
			s: &Subject{
				conflicts: []*Conflict{
					{resolved: true},
					{resolved: false},
				},
			},
			want: 1,
		},
		{
			desc: "resolved",
			s: &Subject{
				conflicts: []*Conflict{
					{resolved: true},
					{resolved: true},
				},
			},
			want: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := len(tc.s.UnresolvedConflicts())
			if got != tc.want {
				t.Errorf("got = %d, but want = %d", got, tc.want)
			}
		})
	}
}

func TestSubjectIsResolved(t *testing.T) {
	cases := []struct {
		desc string
		s    *Subject
		want bool
	}{
		{
			desc: "unresolved",
			s: &Subject{
				conflicts: []*Conflict{
					{resolved: true},
					{resolved: false},
				},
			},
			want: false,
		},
		{
			desc: "resolved",
			s: &Subject{
				conflicts: []*Conflict{
					{resolved: true},
					{resolved: true},
				},
			},
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.s.IsResolved()
			if got != tc.want {
				t.Errorf("got = %t, but want = %t", got, tc.want)
			}
		})
	}
}
