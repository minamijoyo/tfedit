package migration

import (
	"testing"
)

func TestStateImportActionMigrationAction(t *testing.T) {
	cases := []struct {
		desc    string
		address string
		id      string
		want    string
	}{
		{
			desc:    "simple",
			address: "foo_bar.example",
			id:      "test",
			want:    "import foo_bar.example test",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			a := NewStateImportAction(tc.address, tc.id)
			got := a.MigrationAction()

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
