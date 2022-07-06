package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
			name: "with empty prefix",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "log-expiration"
    enabled = true
    prefix  = ""

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 180
      storage_class = "GLACIER"
    }
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
    id = "log-expiration"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 180
      storage_class = "GLACIER"
    }
    status = "Enabled"

    filter {
      prefix = ""
    }
  }
}
`,
		},
		{
			name: "with non-empty prefix",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "log-expiration"
    enabled = true
    prefix  = "log/"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 180
      storage_class = "GLACIER"
    }
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
    id = "log-expiration"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 180
      storage_class = "GLACIER"
    }
    status = "Enabled"

    filter {
      prefix = "log/"
    }
  }
}
`,
		},
		{
			name: "with tags",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "log"
    enabled = true
    prefix = "log/"

    tags = {
      rule      = "log"
      autoclean = "true"
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "GLACIER"
    }

    expiration {
      days = 90
    }
  }

  lifecycle_rule {
    id      = "tmp"
    prefix  = "tmp/"
    enabled = true

    expiration {
      days = 90
    }
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
    id = "log"


    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "GLACIER"
    }

    expiration {
      days = 90
    }
    status = "Enabled"

    filter {

      and {
        prefix = "log/"
        tags = {
          rule      = "log"
          autoclean = "true"
        }
      }
    }
  }

  rule {
    id = "tmp"

    expiration {
      days = 90
    }
    status = "Enabled"

    filter {
      prefix = "tmp/"
    }
  }
}
`,
		},
		{
			name: "with empty prefix and tags",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "log"
    enabled = true
    prefix  = ""
    tags    = {}

    noncurrent_version_transition {
      days          = 30
      storage_class = "GLACIER"
    }

    noncurrent_version_expiration {
      days = 90
    }
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
    id = "log"

    noncurrent_version_transition {
      storage_class   = "GLACIER"
      noncurrent_days = 30
    }

    noncurrent_version_expiration {
      noncurrent_days = 90
    }
    status = "Enabled"

    filter {
    }
  }
}
`,
		},
		{
			name: "with tags but no prefix",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "log"
    enabled = true

    tags = {
      rule      = "log"
      autoclean = "true"
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "GLACIER"
    }

    expiration {
      days = 90
    }
  }

  lifecycle_rule {
    id      = "tmp"
    prefix  = "tmp/"
    enabled = true

    expiration {
      days = 90
    }
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
    id = "log"


    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "GLACIER"
    }

    expiration {
      days = 90
    }
    status = "Enabled"

    filter {

      and {
        prefix = ""
        tags = {
          rule      = "log"
          autoclean = "true"
        }
      }
    }
  }

  rule {
    id = "tmp"

    expiration {
      days = 90
    }
    status = "Enabled"

    filter {
      prefix = "tmp/"
    }
  }
}
`,
		},
		{
			name: "with date in transition and expiration",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "log-expiration"
    enabled = true
    prefix  = "log/"

    transition {
      date          = "2022-03-01"
      storage_class = "STANDARD_IA"
    }

    expiration {
      date = "2022-12-31"
    }
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
    id = "log-expiration"

    transition {
      date          = "2022-03-01T00:00:00Z"
      storage_class = "STANDARD_IA"
    }

    expiration {
      date = "2022-12-31T00:00:00Z"
    }
    status = "Enabled"

    filter {
      prefix = "log/"
    }
  }
}
`,
		},
		{
			name: "with abort_incomplete_multipart_upload_days = 0",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  lifecycle_rule {
    id      = "rule-0"
    enabled = true
    abort_incomplete_multipart_upload_days = 0
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
    id     = "rule-0"
    status = "Enabled"

    filter {
      prefix = ""
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

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
