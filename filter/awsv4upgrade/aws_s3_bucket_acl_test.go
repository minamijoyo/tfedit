package awsv4upgrade

import (
	"testing"

	"github.com/minamijoyo/hcledit/editor"
)

func TestAWSS3BucketACLFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
	}{
		{
			name: "single block",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  acl    = "private"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}
`,
		},
		{
			name: "multiple blocks",
			src: `
resource "aws_s3_bucket" "example1" {
  bucket = "tfedit-test1"
  acl    = "private"
}

resource "aws_s3_bucket" "example2" {
  bucket = "tfedit-test2"
  acl    = "private"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example1" {
  bucket = "tfedit-test1"
}

resource "aws_s3_bucket" "example2" {
  bucket = "tfedit-test2"
}

resource "aws_s3_bucket_acl" "example1" {
  bucket = aws_s3_bucket.example1.id
  acl    = "private"
}

resource "aws_s3_bucket_acl" "example2" {
  bucket = aws_s3_bucket.example2.id
  acl    = "private"
}
`,
		},
		{
			name: "argument not found",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  foo    = "bar"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  foo    = "bar"
}
`,
		},
		{
			name: "resource type not found",
			src: `
resource "aws_s3_bucket_foo" "example" {
  bucket = "tfedit-test"
  acl    = "private"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket_foo" "example" {
  bucket = "tfedit-test"
  acl    = "private"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &AWSS3BucketACLFilter{}
			o := editor.NewEditOperator(filter)
			output, err := o.Apply([]byte(tc.src), "test")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
