# In the sandbox environment with localstack,
# aws s3api list-buckets and aws s3api get-bucket-acl --bucket tfedit-test
# return different canonical user IDs.
# The former has the same ID as the tfstate-test bucket,
# so there may be a problem with the initialization of the sandbox environment.
data "aws_canonical_user_id" "current_user" {}

resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"

  grant {
    id          = data.aws_canonical_user_id.current_user.id
    type        = "CanonicalUser"
    permissions = ["FULL_CONTROL"]
  }

  grant {
    type        = "Group"
    permissions = ["READ_ACP", "WRITE"]
    uri         = "http://acs.amazonaws.com/groups/s3/LogDelivery"
  }
}
