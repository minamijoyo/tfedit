migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_replication_configuration.example tfedit-test",
    "import aws_s3_bucket_versioning.example tfedit-test",
  ]
}
