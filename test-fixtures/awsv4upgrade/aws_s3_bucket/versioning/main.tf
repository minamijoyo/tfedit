resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  versioning {
    enabled = true
  }
}
