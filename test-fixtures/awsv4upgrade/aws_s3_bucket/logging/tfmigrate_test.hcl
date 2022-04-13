migration "state" "test" {
  actions = [
    "import aws_s3_bucket_logging.example tfedit-test",

    # You must give the log-delivery group WRITE and READ_ACP permissions to the target bucket
    "import aws_s3_bucket_acl.log tfedit-log,log-delivery-write",
  ]
}
