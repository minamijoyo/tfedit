resource "aws_s3_bucket" "example" {
  bucket        = "tfedit-test"
  request_payer = "Requester"
}
