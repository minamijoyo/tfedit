package tfwrite

import (
	"testing"
)

func TestTerraformType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
terraform {}
`,
			want: "terraform",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := NewTerraform(findFirstTestBlock(t, f).Raw())

			got := b.Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestTerraformSchemaType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
terraform {}
`,
			want: "",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := NewTerraform(findFirstTestBlock(t, f).Raw())
			got := b.SchemaType()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}
