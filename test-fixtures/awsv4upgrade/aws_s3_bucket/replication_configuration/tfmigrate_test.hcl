migration "state" "test" {
  actions = [
    "import aws_s3_bucket_replication_configuration.example tfedit-test",

    # versioning must be enabled to allow S3 bucket replication
    "import aws_s3_bucket_versioning.example tfedit-test",
  ]
}
