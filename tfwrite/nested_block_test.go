package tfwrite

import (
	"testing"
)

func TestNestedBlockType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		name string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
foo {
  nested {
	  bar = "baz"
  }
}
`,
			name: "nested",
			want: "nested",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			nestedBlocks := b.FindNestedBlocksByType(tc.name)

			got := nestedBlocks[0].Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
