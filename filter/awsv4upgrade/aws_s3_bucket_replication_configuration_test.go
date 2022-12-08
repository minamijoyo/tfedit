package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

func TestAWSS3BucketReplicationConfigurationFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		ok   bool
		want string
	}{
		{
			name: "simple",
			src: `
resource "aws_s3_bucket" "destination" {
  bucket = "tfedit-destination"
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  replication_configuration {
    role = "arn:aws:iam::123456789012:role/tfedit-role"
    rules {
      id     = "foobar"
      status = "Enabled"

      filter {}
      delete_marker_replication_status = "Enabled"

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"

        replication_time {
          status  = "Enabled"
          minutes = 15
        }

        metrics {
          status  = "Enabled"
          minutes = 15
        }
      }
    }
  }

  # versioning must be enabled to allow S3 bucket replication
	versioning {
    enabled = true
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "destination" {
  bucket = "tfedit-destination"
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"


  # versioning must be enabled to allow S3 bucket replication
  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket_replication_configuration" "example" {
  bucket = aws_s3_bucket.example.id
  role   = "arn:aws:iam::123456789012:role/tfedit-role"

  rule {
    id     = "foobar"
    status = "Enabled"

    filter {}

    destination {
      bucket        = aws_s3_bucket.destination.arn
      storage_class = "STANDARD"

      replication_time {
        status = "Enabled"

        time {
          minutes = 15
        }
      }

      metrics {
        status = "Enabled"

        event_threshold {
          minutes = 15
        }
      }
    }

    delete_marker_replication {
      status = "Enabled"
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
			filter := tfeditor.NewAllBlocksFilter(
				&AWSS3BucketFilter{
					filters: []tfeditor.BlockFilter{
						tfeditor.ResourceFilterFunc(AWSS3BucketReplicationConfigurationResourceFilter),
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
