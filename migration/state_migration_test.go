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
			desc:    "empty",
			name:    "mytest",
			dir:     "",
			actions: []StateAction{},
			ok:      true,
			want:    "",
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
		{
			desc: "count",
			name: "mytest",
			dir:  "",
			actions: []StateAction{
				&StateImportAction{
					address: "foo_bar.example[0]",
					id:      "test-0",
				},
				&StateImportAction{
					address: "foo_bar.example[1]",
					id:      "test-1",
				},
			},
			ok: true,
			want: `migration "state" "mytest" {
  actions = [
    "import foo_bar.example[0] test-0",
    "import foo_bar.example[1] test-1",
  ]
}
`,
		},
		{
			desc: "for_each",
			name: "mytest",
			dir:  "",
			actions: []StateAction{
				&StateImportAction{
					address: "foo_bar.example[\"foo\"]",
					id:      "test-foo",
				},
				&StateImportAction{
					address: "foo_bar.example[\"bar\"]",
					id:      "test-bar",
				},
			},
			ok: true,
			want: `migration "state" "mytest" {
  actions = [
    "import 'foo_bar.example[\"foo\"]' test-foo",
    "import 'foo_bar.example[\"bar\"]' test-bar",
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
