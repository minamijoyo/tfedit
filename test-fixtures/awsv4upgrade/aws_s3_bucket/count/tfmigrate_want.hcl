migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_acl.example[0] tfedit-test-0,private",
    "import aws_s3_bucket_acl.example[1] tfedit-test-1,private",
  ]
}
