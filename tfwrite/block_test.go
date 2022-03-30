package tfwrite

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBlockType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "no label",
			src:  `terraform  {}`,
			want: "terraform",
			ok:   true,
		},
		{
			desc: "with a single label",
			src:  `provider "aws" {}`,
			want: "provider",
			ok:   true,
		},
		{
			desc: "with multiple labels",
			src:  `resource "aws_s3_bucket" "example" {}`,
			want: "resource",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			got := b.Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestBlockSetType(t *testing.T) {
	cases := []struct {
		desc     string
		src      string
		typeName string
		want     string
		ok       bool
	}{
		{
			desc: "simple",
			src: `
foo {}
`,
			typeName: "bar",
			want: `
bar {}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := findFirstTestBlock(t, f)
			b.SetType(tc.typeName)
			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
