resource "aws_s3_bucket" "log" {
  bucket = "tfedit-log"

  # You must give the log-delivery group WRITE and READ_ACP permissions to the target bucket
  acl = "log-delivery-write"
}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  logging {
    target_bucket = aws_s3_bucket.log.id
    target_prefix = "log/"
  }
}
