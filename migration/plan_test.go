package migration

import (
	"os"
	"testing"
)

func TestNewPlan(t *testing.T) {
	cases := []struct {
		desc     string
		planFile string
		ok       bool
	}{
		{
			desc:     "valid",
			planFile: "test-fixtures/import_simple.tfplan.json",
			ok:       true,
		},
		{
			desc:     "invalid",
			planFile: "test-fixtures/invalid.tfplan.json",
			ok:       false,
		},
		{
			desc:     "unknown format version",
			planFile: "test-fixtures/unknown_format_version.tfplan.json",
			ok:       false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			planJSON, err := os.ReadFile(tc.planFile)
			if err != nil {
				t.Fatalf("failed to read file: %s", err)
			}

			got, err := NewPlan(planJSON)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %#v", got)
			}
		})
	}
}
