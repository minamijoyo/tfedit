package tfwrite

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
			r := findFirstTestResource(t, f)
			got := r.SchemaType()
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
			r := findFirstTestResource(t, f)
			got := r.Name()
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
