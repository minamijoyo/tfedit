package tfwrite

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseBlock(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "resource",
			src: `
resource "foo_test" "example" {}
`,
			want: "resource",
			ok:   true,
		},
		{
			desc: "data",
			src: `
data "foo_test" "example" {}
`,
			want: "data",
			ok:   true,
		},
		{
			desc: "provider",
			src: `
provider "foo" {}
`,
			want: "provider",
			ok:   true,
		},
		{
			desc: "unknown",
			src: `
b "foo" {}
`,
			want: "b",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := parseBlock(findFirstTestBlock(t, f).Raw())
			got := b.Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestFileBlocks(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want []Block
		ok   bool
	}{
		{
			desc: "simple",
			src: `
resource "foo_test" "example1" {}
resource "foo_test" "example2" {}
resource "foo_bar" "example1" {}
data "foo_test" "example1" {}
`,
			want: []Block{
				NewEmptyResource("foo_test", "example1"),
				NewEmptyResource("foo_test", "example2"),
				NewEmptyResource("foo_bar", "example1"),
				NewEmptyDataSource("foo_test", "example1"),
			},
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			got := f.Blocks()
			opts := cmpopts.IgnoreUnexported(Resource{}, DataSource{})
			if diff := cmp.Diff(got, tc.want, opts); diff != "" {
				t.Errorf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestFileFindBlocksByType(t *testing.T) {
	cases := []struct {
		desc       string
		src        string
		blockType  string
		schemaType string
		want       string
		ok         bool
	}{
		{
			desc: "resource",
			src: `
resource "foo_test" "example1" {}
resource "foo_test" "example2" {}
resource "foo_bar" "example1" {}
data "foo_test" "example1" {}
`,
			blockType:  "resource",
			schemaType: "",
			want: `
resource "foo_test" "example1" {}

resource "foo_test" "example2" {}

resource "foo_bar" "example1" {}
`,
			ok: true,
		},
		{
			desc: "resource with schemaType",
			src: `
resource "foo_test" "example1" {}
resource "foo_test" "example2" {}
resource "foo_bar" "example1" {}
data "foo_test" "example1" {}
`,
			blockType:  "resource",
			schemaType: "foo_test",
			want: `
resource "foo_test" "example1" {}

resource "foo_test" "example2" {}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			blocks := f.FindBlocksByType(tc.blockType, tc.schemaType)

			newFile := NewEmptyFile()
			for _, r := range blocks {
				newFile.AppendBlock(r)
			}

			got := printTestFile(t, newFile)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}

func TestFileAppendBlock(t *testing.T) {
	cases := []struct {
		desc         string
		src          string
		resourceType string
		resourceName string
		want         string
		ok           bool
	}{
		{
			desc: "simple",
			src: `
resource "foo_test" "example1" {}
`,
			resourceType: "foo_test",
			resourceName: "example2",
			want: `
resource "foo_test" "example1" {}

resource "foo_test" "example2" {
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := NewEmptyResource(tc.resourceType, tc.resourceName)
			f.AppendBlock(r)

			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
