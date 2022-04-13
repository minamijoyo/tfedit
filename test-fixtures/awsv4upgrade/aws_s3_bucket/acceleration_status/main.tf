resource "aws_s3_bucket" "example" {
  bucket              = "tfedit-test"
  acceleration_status = "Enabled"
}
