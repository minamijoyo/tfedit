package awsv4upgrade

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minamijoyo/hcledit/editor"
	"github.com/minamijoyo/tfedit/tfeditor"
)

func TestAWSS3BucketFilter(t *testing.T) {
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
			name: "multiple arguments",
			src: `
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  acl    = "private"

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

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_s3_bucket_logging" "example" {
  bucket = aws_s3_bucket.example.id

  target_bucket = "tfedit-test-log"
  target_prefix = "log/"
}
`,
		},
		{
			name: "multiple resources",
			src: `
resource "aws_s3_bucket" "example1" {
  bucket = "tfedit-test1"
  acl    = "private"

  logging {
    target_bucket = "tfedit-test-log"
    target_prefix = "log/"
  }
}

resource "aws_s3_bucket" "example2" {
  bucket = "tfedit-test2"
  acl    = "private"

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

resource "aws_s3_bucket_acl" "example1" {
  bucket = aws_s3_bucket.example1.id
  acl    = "private"
}

resource "aws_s3_bucket_logging" "example1" {
  bucket = aws_s3_bucket.example1.id

  target_bucket = "tfedit-test-log"
  target_prefix = "log/"
}

resource "aws_s3_bucket_acl" "example2" {
  bucket = aws_s3_bucket.example2.id
  acl    = "private"
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
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket_foo" "example" {
  bucket = "tfedit-test"
}
`,
		},
		{
			name: "full arguments except for grant",
			src: `
resource "aws_s3_bucket" "log" {
  bucket = "tfedit-log"

  # You must give the log-delivery group WRITE and READ_ACP permissions to the target bucket
  acl = "log-delivery-write"
}

resource "aws_s3_bucket" "destination" {
  bucket = "tfedit-destination"
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  acceleration_status = "Enabled"
  acl    = "private"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }

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

  logging {
    target_bucket = aws_s3_bucket.log.id
    target_prefix = "log/"
  }

  object_lock_configuration {
    object_lock_enabled = "Enabled"

    rule {
      default_retention {
        mode = "COMPLIANCE"
        days = 3
      }
    }
  }

  policy = <<EOF
{
  "Id": "Policy1446577137248",
  "Statement": [
    {
      "Action": "s3:PutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:root"
      },
      "Resource": "arn:aws:s3:::example/*",
      "Sid": "Stmt1446575236270"
    }
  ],
  "Version": "2012-10-17"
}
EOF

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

  request_payer = "Requester"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = "aws/s3"
        sse_algorithm     = "aws:kms"
      }
    }
  }

  versioning {
    enabled = true
  }

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "log" {
  bucket = "tfedit-log"
}

resource "aws_s3_bucket" "destination" {
  bucket = "tfedit-destination"
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  object_lock_enabled = true
}

resource "aws_s3_bucket_acl" "log" {
  bucket = aws_s3_bucket.log.id
  # You must give the log-delivery group WRITE and READ_ACP permissions to the target bucket
  acl = "log-delivery-write"
}

resource "aws_s3_bucket_accelerate_configuration" "example" {
  bucket = aws_s3_bucket.example.id
  status = "Enabled"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
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

resource "aws_s3_bucket_logging" "example" {
  bucket = aws_s3_bucket.example.id

  target_bucket = aws_s3_bucket.log.id
  target_prefix = "log/"
}

resource "aws_s3_bucket_object_lock_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  rule {
    default_retention {
      mode = "COMPLIANCE"
      days = 3
    }
  }
}

resource "aws_s3_bucket_policy" "example" {
  bucket = aws_s3_bucket.example.id
  policy = <<EOF
{
  "Id": "Policy1446577137248",
  "Statement": [
    {
      "Action": "s3:PutObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::123456789012:root"
      },
      "Resource": "arn:aws:s3:::example/*",
      "Sid": "Stmt1446575236270"
    }
  ],
  "Version": "2012-10-17"
}
EOF
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

resource "aws_s3_bucket_request_payment_configuration" "example" {
  bucket = aws_s3_bucket.example.id
  payer  = "Requester"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "example" {
  bucket = aws_s3_bucket.example.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = "aws/s3"
      sse_algorithm     = "aws:kms"
    }
  }
}

resource "aws_s3_bucket_versioning" "example" {
  bucket = aws_s3_bucket.example.id

  versioning_configuration {
    status = "Enabled"
  }
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
			name: "grant (conflict with acl)",
			src: `
data "aws_canonical_user_id" "current_user" {}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  grant {
    id          = data.aws_canonical_user_id.current_user.id
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
  }

  grant {
    type        = "Group"
    permissions = ["READ_ACP", "WRITE"]
    uri         = "http://acs.amazonaws.com/groups/s3/LogDelivery"
  }
}
`,
			ok: true,
			want: `
data "aws_canonical_user_id" "current_user" {}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id

  access_control_policy {

    grant {

      grantee {

        id   = data.aws_canonical_user_id.current_user.id
        type = "CanonicalUser"
      }
      permission = "FULL_CONTROL"
    }

    grant {

      grantee {

        type = "Group"
        uri  = "http://acs.amazonaws.com/groups/s3/LogDelivery"
      }
      permission = "READ_ACP"
    }

    grant {

      grantee {

        type = "Group"
        uri  = "http://acs.amazonaws.com/groups/s3/LogDelivery"
      }
      permission = "WRITE"
    }

    owner {
      id = "set_aws_canonical_user_id"
    }
  }
}
`,
		},
		{
			name: "count",
			src: `
resource "aws_s3_bucket" "example" {
  count  = 2
  bucket = "tfedit-test-${count.index}"
  acl    = "private"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  count  = 2
  bucket = "tfedit-test-${count.index}"
}

resource "aws_s3_bucket_acl" "example" {
  count  = 2
  bucket = aws_s3_bucket.example[count.index].id
  acl    = "private"
}
`,
		},
		{
			name: "for_each",
			src: `
resource "aws_s3_bucket" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = "tfedit-test-${each.key}"
  acl      = "private"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = "tfedit-test-${each.key}"
}

resource "aws_s3_bucket_acl" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = aws_s3_bucket.example[each.key].id
  acl      = "private"
}
`,
		},
		{
			name: "provider",
			src: `
resource "aws_s3_bucket" "example" {
  provider = aws.foo
  bucket   = "tfedit-test"
  acl      = "private"
}
`,
			ok: true,
			want: `
resource "aws_s3_bucket" "example" {
  provider = aws.foo
  bucket   = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  provider = aws.foo
  bucket   = aws_s3_bucket.example.id
  acl      = "private"
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			filter := tfeditor.NewAllBlocksFilter(NewAWSS3BucketFilter())
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
