package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

func TestAWSS3BucketVersioningFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
	}{
		{
			name: "enabled = true",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  versioning {
    enabled = true
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    status = "Enabled"
  }
}
`,
		},
		{
			name: "enabled = false",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  versioning {
    enabled = false
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    status = "Suspended"
  }
}
`,
		},
		{
			name: "enabled = var.enabled",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  versioning {
    enabled = var.enabled
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    status = var.enabled
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
		{
			name: "mfa_delete = true",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  versioning {
    enabled    = true
    mfa_delete = true
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    mfa_delete = "Enabled"
    status     = "Enabled"
  }
}
`,
		},
		{
			name: "mfa_delete = false",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  versioning {
    enabled    = true
    mfa_delete = false
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    mfa_delete = "Disabled"
    status     = "Enabled"
  }
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := &AWSS3BucketFilter{filters: []tfeditor.ResourceFilter{&AWSS3BucketVersioningFilter{}}}
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
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%v", got, tc.want, diff)
			}
		})
	}
}
