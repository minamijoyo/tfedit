migration "state" "fromplan" {
  actions = [
    "import 'aws_s3_bucket_acl.example[\"bar\"]' tfedit-test-bar,private",
    "import 'aws_s3_bucket_acl.example[\"foo\"]' tfedit-test-foo,private",
  ]
}
