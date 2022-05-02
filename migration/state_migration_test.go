package migration

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStateMigrationRender(t *testing.T) {
	cases := []struct {
		desc    string
		name    string
		dir     string
		actions []StateAction
		m       *StateMigration
		ok      bool
		want    string
	}{
		{
			desc: "simple",
			name: "mytest",
			dir:  "",
			actions: []StateAction{
				&StateImportAction{
					address: "foo_bar.example1",
					id:      "test1",
				},
				&StateImportAction{
					address: "foo_bar.example2",
					id:      "test2",
				},
			},
			ok: true,
			want: `migration "state" "mytest" {
  actions = [
    "import foo_bar.example1 test1",
    "import foo_bar.example2 test2",
  ]
}
`,
		},
		{
			desc: "simple with dir",
			name: "mytest",
			dir:  "tmp/dir1",
			actions: []StateAction{
				&StateImportAction{
					address: "foo_bar.example1",
					id:      "test1",
				},
				&StateImportAction{
					address: "foo_bar.example2",
					id:      "test2",
				},
			},
			ok: true,
			want: `migration "state" "mytest" {
  dir = "tmp/dir1"
  actions = [
    "import foo_bar.example1 test1",
    "import foo_bar.example2 test2",
  ]
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			m := NewStateMigration(tc.name, tc.dir)
			m.AppendActions(tc.actions...)
			output, err := m.Render()
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %s", got)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
