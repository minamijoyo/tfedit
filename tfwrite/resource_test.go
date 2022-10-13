package tfwrite

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestResourceType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
resource "foo_test" "example" {
  nested {
	  bar = "baz"
  }
}
`,
			want: "resource",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := NewResource(findFirstTestBlock(t, f).Raw())

			got := b.Type()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestResourceSchemaType(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
resource "aws_s3_bucket" "example" {}
`,
			want: "aws_s3_bucket",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := NewResource(findFirstTestBlock(t, f).Raw())
			got := b.SchemaType()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestResourceName(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
resource "aws_s3_bucket" "example" {}
`,
			want: "example",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := NewResource(findFirstTestBlock(t, f).Raw())
			got := b.Name()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestResourceReferableName(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "count",
			src: `
resource "aws_s3_bucket" "example" {
  count = 2
}
`,
			want: "example[count.index]",
			ok:   true,
		},
		{
			desc: "for_each",
			src: `
resource "aws_s3_bucket" "example" {
  for_each = toset(["foo", "bar"])
}
`,
			want: "example[each.key]",
			ok:   true,
		},
		{
			desc: "default",
			src: `
resource "aws_s3_bucket" "example" {
}
`,
			want: "example",
			ok:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			b := NewResource(findFirstTestBlock(t, f).Raw())
			got := b.ReferableName()
			if got != tc.want {
				t.Errorf("got = %s, but want = %s", got, tc.want)
			}
		})
	}
}

func TestResourceSetAttributeByReference(t *testing.T) {
	cases := []struct {
		desc string
		src  string
		want string
		ok   bool
	}{
		{
			desc: "simple",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
}
`,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
}
`,
			ok: true,
		},
		{
			desc: "count",
			src: `
resource "aws_s3_bucket" "example" {
  count  = 2
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  count = 2
}
`,
			want: `
resource "aws_s3_bucket" "example" {
  count  = 2
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  count  = 2
  bucket = aws_s3_bucket.example[count.index].id
}
`,
			ok: true,
		},
		{
			desc: "for_each",
			src: `
resource "aws_s3_bucket" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  for_each = toset(["foo", "bar"])
}
`,
			want: `
resource "aws_s3_bucket" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = aws_s3_bucket.example[each.key].id
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := parseTestFile(t, tc.src)
			r := f.FindResourcesByType("aws_s3_bucket_acl")[0]
			refResource := f.FindResourcesByType("aws_s3_bucket")[0]
			r.SetAttributeByReference("bucket", refResource, "id")
			got := printTestFile(t, f)
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
