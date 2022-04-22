migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_acl.log tfedit-log,log-delivery-write",
    "import aws_s3_bucket_accelerate_configuration.example tfedit-test",
    "import aws_s3_bucket_acl.example tfedit-test,private",
    "import aws_s3_bucket_cors_configuration.example tfedit-test",
    "import aws_s3_bucket_lifecycle_configuration.example tfedit-test",
    "import aws_s3_bucket_logging.example tfedit-test",
    "import aws_s3_bucket_object_lock_configuration.example tfedit-test",
    "import aws_s3_bucket_policy.example tfedit-test",
    "import aws_s3_bucket_replication_configuration.example tfedit-test",
    "import aws_s3_bucket_request_payment_configuration.example tfedit-test",
    "import aws_s3_bucket_server_side_encryption_configuration.example tfedit-test",
    "import aws_s3_bucket_versioning.example tfedit-test",
    "import aws_s3_bucket_website_configuration.example tfedit-test",
  ]
}
