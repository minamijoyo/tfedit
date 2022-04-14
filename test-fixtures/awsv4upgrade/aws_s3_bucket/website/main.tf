resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
