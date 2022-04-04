package tfwrite

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func TestSplitTokensAsList(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want []hclwrite.Tokens
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  attr = ["foo", "bar"]
}
`,
			want: []hclwrite.Tokens{
				hclwrite.TokensForValue(cty.StringVal("foo")),
				hclwrite.TokensForValue(cty.StringVal("bar")),
			},
			ok: true,
		},
		{
			desc: "empty string (invalid list)",
			src: `
foo {
  attr = ""
}
`,
			want: nil,
			ok:   true,
		},
		{
			desc: "empty list",
			src: `
foo {
  attr = []
}
`,
			want: []hclwrite.Tokens{},
			ok:   true,
		},
		{
			desc: "variable",
			src: `
foo {
  attr = [var.foo, var.bar]
}
`,
			want: []hclwrite.Tokens{
				hclwrite.TokensForTraversal(hcl.Traversal{
					hcl.TraverseRoot{Name: "var"},
					hcl.TraverseAttr{Name: "foo"},
				}),
				hclwrite.TokensForTraversal(hcl.Traversal{
					hcl.TraverseRoot{Name: "var"},
					hcl.TraverseAttr{Name: "bar"},
				}),
			},
			ok: true,
		},
		{
			desc: "multi lines",
			src: `
foo {
  attr = [
    "foo",
    "bar"
  ]
}
`,
			want: []hclwrite.Tokens{
				hclwrite.TokensForValue(cty.StringVal("foo")),
				hclwrite.TokensForValue(cty.StringVal("bar")),
			},
			ok: true,
		},
		{
			desc: "multi lines with comma",
			src: `
foo {
  attr = [
    "foo",
    "bar",
  ]
}
`,
			want: []hclwrite.Tokens{
				hclwrite.TokensForValue(cty.StringVal("foo")),
				hclwrite.TokensForValue(cty.StringVal("bar")),
			},
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			attr := b.GetAttribute("attr")
			tokens := attr.ValueAsTokens()
			got := SplitTokensAsList(tokens)

			if diff := cmp.Diff(got, tc.want, cmpopts.IgnoreFields(hclwrite.Token{}, "SpacesBefore")); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%v", spew.Sdump(got), spew.Sdump(tc.want), diff)
			}
		})
	}
}
