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
