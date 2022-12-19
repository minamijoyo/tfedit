package tfeditor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfwrite"
	"github.com/zclconf/go-cty/cty"
)

func TestFileFilter(t *testing.T) {
	cases := []struct {
		name   string
		src    string
		filter BlockFilter
		ok     bool
		want   string
	}{
		{
			name: "simple",
			src: `
block "foo_bar1" "example11" {
  baz = "test1"
}

block "foo_bar2" "example21" {
  baz = "test1"
}

block "foo_bar1" "example13" {
  baz = "test1"
}
`,
			filter: BlockFilterFunc(func(inFile *tfwrite.File, block tfwrite.Block) (*tfwrite.File, error) {
				block.SetAttributeValue("baz", cty.StringVal("test2"))
				return inFile, nil
			}),
			ok: true,
			want: `
block "foo_bar1" "example11" {
  baz = "test2"
}

block "foo_bar2" "example21" {
  baz = "test2"
}

block "foo_bar1" "example13" {
  baz = "test2"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := NewFileFilter(tc.filter)
			o := editor.NewEditOperator(filter)
			output, err := o.Apply([]byte(tc.src), "test")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
