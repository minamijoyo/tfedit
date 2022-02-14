package awsv4upgrade

import (
	"testing"

	"github.com/minamijoyo/hcledit/editor"
)

func TestAWSS3BucketLoggingFilter(t *testing.T) {
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
  logging {
    target_bucket = "tfedit-test-log"
    target_prefix = "log/"
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_logging" "example" {
  bucket = aws_s3_bucket.example.id

  target_bucket = "tfedit-test-log"
  target_prefix = "log/"
}
`,
		},
		{
			name: "multiple blocks",
			src: `
resource "aws_s3_bucket" "example1" {
  bucket = "tfedit-test1"
  logging {
    target_bucket = "tfedit-test-log"
    target_prefix = "log/"
  }
}

resource "aws_s3_bucket" "example2" {
  bucket = "tfedit-test2"
  logging {
    target_bucket = "tfedit-test-log"
    target_prefix = "log/"
  }
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

resource "aws_s3_bucket_logging" "example1" {
  bucket = aws_s3_bucket.example1.id

  target_bucket = "tfedit-test-log"
  target_prefix = "log/"
}

resource "aws_s3_bucket_logging" "example2" {
  bucket = aws_s3_bucket.example2.id

  target_bucket = "tfedit-test-log"
  target_prefix = "log/"
}
`,
		},
		{
			name: "argument not found",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  foo {}
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  foo {}
}
`,
		},
		{
			name: "resource type not found",
			src: `
resource "aws_s3_bucket_foo" "example" {
  bucket = "tfedit-test"
  logging {
    target_bucket = "tfedit-test-log"
    target_prefix = "log/"
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket_foo" "example" {
  bucket = "tfedit-test"
  logging {
    target_bucket = "tfedit-test-log"
    target_prefix = "log/"
  }
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &AWSS3BucketLoggingFilter{}
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
