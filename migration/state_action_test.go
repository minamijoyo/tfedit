package migration

import (
	"testing"
)

func TestActionEscape(t *testing.T) {
	cases := []struct {
		desc string
		raw  string
		want string
	}{
		{
			desc: "noop",
			raw:  "foo",
			want: "foo",
		},
		{
			desc: "double quote",
			raw:  `"foo"`,
			want: `'\"foo\"'`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := actionEscape(tc.raw)

			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

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
		{
			desc:    "count",
			address: "foo_bar.example[0]",
			id:      "test",
			want:    "import foo_bar.example[0] test",
		},
		{
			desc:    "for_each",
			address: `foo_bar.example["foo"]`,
			id:      "test",
			want:    `import 'foo_bar.example[\"foo\"]' test`,
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
