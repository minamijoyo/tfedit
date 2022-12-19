package tfeditor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
)

func TestVerticalFormatterBlockFilter(t *testing.T) {
	cases := []struct {
		name       string
		src        string
		blockType  string
		schemaType string
		filter     BlockFilter
		ok         bool
		want       string
	}{
		{
			name: "simple",
			src: `
block "foo_bar" "example" {
  baz = "test1"

}
`,
			blockType:  "block",
			schemaType: "foo_bar",
			ok:         true,
			want: `
block "foo_bar" "example" {
  baz = "test1"
}
`,
		},
		{
			name: "block type not found",
			src: `
block "foo_bar" "example" {
  baz = "test1"

}
`,
			blockType:  "aaa",
			schemaType: "foo_bar",
			ok:         true,
			want: `
block "foo_bar" "example" {
  baz = "test1"

}
`,
		},
		{
			name: "schema type not found",
			src: `
block "foo_bar" "example" {
  baz = "test1"

}
`,
			blockType:  "block",
			schemaType: "aaa",
			ok:         true,
			want: `
block "foo_bar" "example" {
  baz = "test1"

}
`,
		},
		{
			name: "multiple blocks",
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
			blockType:  "block",
			schemaType: "foo_bar1",
			ok:         true,
			want: `
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
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bf := NewVerticalFormatterBlockFilter(tc.blockType, tc.schemaType)
			filter := NewAllBlocksFilter(bf)
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
