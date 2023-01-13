package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
)

func TestAWSS3BucketWebsiteFilter(t *testing.T) {
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

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_website_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
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
			name: "rename references for website_domain and website_endpoint",
			src: `
resource "aws_route53_zone" "test" {
  name = "example.com"
}

resource "aws_route53_record" "alias" {
  zone_id = aws_route53_zone.test.zone_id
  name    = "www"
  type    = "A"

  alias {
    zone_id                = aws_s3_bucket.example.hosted_zone_id
    name                   = aws_s3_bucket.example.website_domain
    evaluate_target_health = true
  }
}

output "test_endpoint" {
  value = aws_s3_bucket.example.website_endpoint
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`,
			ok: true,
			want: `
resource "aws_route53_zone" "test" {
  name = "example.com"
}

resource "aws_route53_record" "alias" {
  zone_id = aws_route53_zone.test.zone_id
  name    = "www"
  type    = "A"

  alias {
    zone_id                = aws_s3_bucket.example.hosted_zone_id
    name                   = aws_s3_bucket_website_configuration.example.website_domain
    evaluate_target_health = true
  }
}

output "test_endpoint" {
  value = aws_s3_bucket_website_configuration.example.website_endpoint
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

}

resource "aws_s3_bucket_website_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := buildTestBlockFilter(AWSS3BucketWebsiteBlockFilter)
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
