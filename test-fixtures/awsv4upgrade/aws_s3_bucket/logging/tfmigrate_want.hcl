migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_logging.example tfedit-test",
    "import aws_s3_bucket_acl.log tfedit-log,log-delivery-write",
  ]
}
