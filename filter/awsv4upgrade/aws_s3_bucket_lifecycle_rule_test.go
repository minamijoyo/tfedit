package awsv4upgrade

import (
	"testing"

	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

func TestAWSS3BucketLifecycleRuleFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
	}{
		{
			name: "simple",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "Keep previous version 30 days, then in Glacier another 60"
    enabled = true

    noncurrent_version_transition {
      days          = 30
      storage_class = "GLACIER"
    }

    noncurrent_version_expiration {
      days = 90
    }
  }

  lifecycle_rule {
    id                                     = "Delete old incomplete multi-part uploads"
    enabled                                = true
    abort_incomplete_multipart_upload_days = 7
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"


}

resource "aws_s3_bucket_lifecycle_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  rule {
    id = "Keep previous version 30 days, then in Glacier another 60"

    noncurrent_version_transition {
      storage_class   = "GLACIER"
      noncurrent_days = 30
    }

    noncurrent_version_expiration {
      noncurrent_days = 90
    }
    status = "Enabled"

    filter {
      prefix = ""
    }
  }

  rule {
    id     = "Delete old incomplete multi-part uploads"
    status = "Enabled"

    filter {
      prefix = ""
    }

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &AWSS3BucketFilter{filters: []tfeditor.ResourceFilter{&AWSS3BucketLifecycleRuleFilter{}}}
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
