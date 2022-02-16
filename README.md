# tfedit
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/minamijoyo/tfedit.svg)](https://github.com/minamijoyo/tfedit/releases/latest)
[![GoDoc](https://godoc.org/github.com/minamijoyo/tfedit/tfedit?status.svg)](https://godoc.org/github.com/minamijoyo/tfedit)

## Features

- Easy refactoring Terraform configurations in a scalable way with built-in filters.
- CLI-friendly: Read HCL from stdin, apply filters and write results to stdout, easily pipe and combine other commands. You can also update a given file in-place.
- Keep comments: You can update lots of existing Terraform configurations without losing any valuable comments.
- Available operations:
  - filter awsv4upgrade

For upgrading AWS provider v4, some rules have not been implemented yet. The current implementation status is as follows:

[S3 Bucket Refactor](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor)

- [ ] acceleration_status Argument
- [x] acl Argument
- [ ] cors_rule Argument
- [ ] grant Argument
- [ ] lifecycle_rule Argument
- [x] logging Argument
- [ ] object_lock_configuration rule Argument
- [ ] policy Argument
- [ ] replication_configuration Argument
- [ ] request_payer Argument
- [ ] server_side_encryption_configuration Argument
- [ ] versioning Argument
- [ ] website, website_domain, and website_endpoint Arguments

The tfedit focuses on refactoring Terraform configurations in a scalable way. The initial goal of this project is a sharing migration tool for bulk refactoring `aws_s3_bucket` resource caused by breaking changes in AWS provider v4, but the project scope is not limited to specific use-cases. The scope of this project is sharing reusable operations often used in Terraform refactoring and well-known refactoring use-case as a built-in filter.

As you know, Terraform refactoring often requires not only configuration changes, but also terraform state migrations. However, it's error-prone and not suitable for a GitOps workflow. To migrate Terraform state in a declarative way, use [tfmigrate](https://github.com/minamijoyo/tfmigrate).

If you are not ready for the upgrade, you can pin version constraints in your Terraform configurations with [tfupdate](https://github.com/minamijoyo/tfupdate).

## Install

### Source

If you have Go 1.17+ development environment:

```
$ go install github.com/minamijoyo/tfedit@latest
$ tfedit version
```

## Usage

```
$ tfedit --help
A refactoring tool for Terraform

Usage:
  tfedit [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  filter      Apply a built-in filter
  help        Help about any command
  version     Print version

Flags:
  -f, --file string   A path of input file (default "-")
  -h, --help          help for tfedit
  -u, --update        Update files in-place

Use "tfedit [command] --help" for more information about a command.
```

```
$ tfedit filter --help
Apply a built-in filter

Arguments:
  FILTER_TYPE    A type of filter.
                 Valid values are:
                 - awsv4upgrade
                   Upgrade configurations to AWS provider v4.
                   Only aws_s3_bucket refactor is supported.

Usage:
  tfedit filter <FILTER_TYPE> [flags]

Flags:
  -h, --help   help for filter

Global Flags:
  -f, --file string   A path of input file (default "-")
  -u, --update        Update files in-place
```

By default, the input is read from stdin, and the output is written to stdout.
You can also read a file with `-f` flag, and update the file in-place with `-u` flag.

## Example

Given the following file:

```aws_s3_bucket.tf
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

```

```aws_s3_bucket.tf
$ tfedit filter awsv4upgrade -f ./test-fixtures/awsv4upgrade/aws_s3_bucket.tf
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

}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}

resource "aws_s3_bucket_logging" "example" {
  bucket = aws_s3_bucket.example.id

  target_bucket = "minamijoyo-tf-aws-v4-test1-log"
  target_prefix = "log/"
}
```

## License

MIT
