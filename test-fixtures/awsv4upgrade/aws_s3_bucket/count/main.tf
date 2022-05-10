resource "aws_s3_bucket" "example" {
  count  = 2
  bucket = "tfedit-test-${count.index}"
  acl    = "private"
}
