terraform {
  # https://www.terraform.io/docs/backends/types/s3.html
  backend "s3" {
    region = "ap-northeast-1"
    bucket = "tfstate-test"
    key    = "test/terraform.tfstate"

    # mock s3 endpoint with localstack
    endpoint                    = "http://localstack:4566"
    access_key                  = "dummy"
    secret_key                  = "dummy"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    force_path_style            = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "3.74.3"
    }
  }
}

# https://www.terraform.io/docs/providers/aws/index.html
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/custom-service-endpoints#localstack
provider "aws" {
  region = "ap-northeast-1"

  access_key                  = "dummy"
  secret_key                  = "dummy"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_region_validation      = true
  skip_requesting_account_id  = true

  # mock endpoints with localstack
  endpoints {
    s3 = "http://localstack:4566"
  }

  s3_force_path_style = true
}
