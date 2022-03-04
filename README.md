# tfedit
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/minamijoyo/tfedit.svg)](https://github.com/minamijoyo/tfedit/releases/latest)
[![GoDoc](https://godoc.org/github.com/minamijoyo/tfedit/tfedit?status.svg)](https://godoc.org/github.com/minamijoyo/tfedit)

## Features

Easy refactoring Terraform configurations in a scalable way.

- CLI-friendly: Read HCL from stdin, apply filters and write results to stdout, easily pipe and combine other commands.
- Keep comments: Update lots of existing Terraform configurations without losing comments as much as possible.
- Built-in operations:
  - filter awsv4upgrade: Upgrade configurations to AWS provider v4. Only `aws_s3_bucket` refactor is supported.

Although the initial goal of this project is providing a way for bulk refactoring of the `aws_s3_bucket` resource required by breaking changes in AWS provider v4, but the project scope is not limited to specific use-cases. It's by no means intended to be an upgrade tool for all your providers. Instead of covering all you need, it provides reusable building blocks for Terraform refactoring and shows examples for how to compose them in real world use-cases.

As you know, Terraform refactoring often requires not only configuration changes, but also Terraform state migrations. However, it's error-prone and not suitable for CI/CD. For declarative Terraform state migration, use [tfmigrate](https://github.com/minamijoyo/tfmigrate).

If you are not ready for the upgrade, you can pin version constraints in your Terraform configurations with [tfupdate](https://github.com/minamijoyo/tfupdate).

## awsv4upgrade

For upgrading AWS provider v4, some rules have not been implemented yet. The current implementation status is as follows:

[S3 Bucket Refactor](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor)

- [ ] Arguments of aws_s3_bucket resource
  - [ ] acceleration_status
  - [x] acl
  - [x] cors_rule
  - [ ] grant
  - [x] lifecycle_rule
  - [x] logging
  - [ ] object_lock_configuration rule
  - [x] policy
  - [ ] replication_configuration
  - [ ] request_payer
  - [x] server_side_encryption_configuration
  - [x] versioning
  - [ ] website
- [ ] Meta arguments of resource
  - [ ] count
  - [ ] for_each
  - [ ] dynamic
- [ ] Rename references in an expression to new resource type
- [ ] Generate import commands for new split resources

Known limitations:
- Some arguments were changed not only their names but also valid values. In this case, if a value of the argument is a variable, not literal, it's impossible to automatically rewrite the value of the variable. It potentially could be passed from outside of module or even overwritten at runtime. If it's not literal, you need to change the value of the variable by yourself. The following arguments have this limitation:
  - lifecycle_rule:
    - enabled = true => status = "Enabled"
    - enabled = false => status = "Disabled"
    - transition:
      - date = "2022-12-31" => date = "2022-12-31T00:00:00Z"
    - expiration:
      - date = "2022-12-31" => date = "2022-12-31T00:00:00Z"
  - versioning:
    - enabled = true => status = "Enabled"
    - enabled = false => It also depends on the current status of your bucket. Set `status = "Suspended"` or use `for_each` to avoid creating `aws_s3_bucket_versioning` resource.
- Some arguments cannot be converted correctly without knowing the current state of AWS resources. The tfedit never calls the AWS API on your behalf. You have to check it by yourself. The following arguments have this limitation:
  - lifecycle_rule:
    - filter: When [`aws s3api get-bucket-lifecycle-configuration --bucket <bucketname>`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/get-bucket-lifecycle-configuration.html) returns `"Filter" : {}` without a prefix, you need to set rule.filter as `filter {}`.
  - versioning:
    - enabled: Starting from v3.70.0, `enabled = false` for a new bucket doesn't set "Suspended" explicitly. When [`aws s3api get-bucket-versioning --bucket <bucketname>`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/get-bucket-versioning.html) returns no `"Status"`, which means `"Disabled"`. In this case, you need to remove the `aws_s3_bucket_versioning` resource.

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
$ cat ./test-fixtures/awsv4upgrade/aws_s3_bucket.tf
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

Apply a filter for `awsv4upgrade`:

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
