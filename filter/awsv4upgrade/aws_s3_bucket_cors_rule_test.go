package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

func TestAWSS3BucketCorsRuleFilter(t *testing.T) {
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

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_cors_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
`,
		},
		{
			name: "multiple rules",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"


}

resource "aws_s3_bucket_cors_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
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
			filter := tfeditor.NewAllBlocksFilter(
				&AWSS3BucketFilter{
					filters: []tfeditor.BlockFilter{
						tfeditor.ResourceFilterFunc(AWSS3BucketCorsRuleResourceFilter),
					},
				},
			)
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
