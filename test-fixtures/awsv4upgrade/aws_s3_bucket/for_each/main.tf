resource "aws_s3_bucket" "example" {
  for_each = toset(["foo", "bar"])
  bucket   = "tfedit-test-${each.key}"
  acl      = "private"
}
