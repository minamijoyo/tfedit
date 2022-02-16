terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.0.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-1"
}

resource "aws_s3_bucket" "example" {
  bucket = "minamijoyo-tf-aws-v4-test1"
  acl    = "private"

  logging {
    target_bucket = "minamijoyo-tf-aws-v4-test1-log"
    target_prefix = "log/"
  }
}
