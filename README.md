# tfedit
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![GitHub release](https://img.shields.io/github/release/minamijoyo/tfedit.svg)](https://github.com/minamijoyo/tfedit/releases/latest)
[![GoDoc](https://godoc.org/github.com/minamijoyo/tfedit/tfedit?status.svg)](https://godoc.org/github.com/minamijoyo/tfedit)

## Features

Easy refactoring Terraform configurations in a scalable way.

- CLI-friendly: Read HCL from stdin, apply filters and write results to stdout, easily pipe and combine other commands.
- Keep comments: Update lots of existing Terraform configurations without losing comments as much as possible.
- Built-in operations:
  - filter awsv4upgrade: Upgrade configurations to AWS provider v4.
- Generate a migration file for state operations: Read a Terraform plan file in JSON format and generate a migration file in [tfmigrate](https://github.com/minamijoyo/tfmigrate) HCL format. Currently, only import actions required by awsv4upgrade are supported.

Although the initial goal of this project is providing a way for bulk refactoring of the `aws_s3_bucket` resource required by breaking changes in AWS provider v4, but the project scope is not limited to specific use-cases. It's by no means intended to be an upgrade tool for all your providers. Instead of covering all you need, it provides reusable building blocks for Terraform refactoring and shows examples for how to compose them in real world use-cases.

## awsv4upgrade
### Overview

In short, given the following Terraform configuration file for the AWS provider v3:

```main.tf
$ cat ./test-fixtures/awsv4upgrade/aws_s3_bucket/simple/main.tf
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  acl    = "private"
}
```

Apply a filter for `awsv4upgrade`:

```main.tf
$ tfedit filter awsv4upgrade -f ./test-fixtures/awsv4upgrade/aws_s3_bucket/simple/main.tf
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}
```

You can see the `acl` argument has been split into an `aws_s3_bucket_acl` resource for the AWS provider v4 compatible.

To resolve the conflict between the configuration and the existing state, you need to import the new resource. As you know, you can run the `terraform import` command directly, but if you prefer to check the upgrade results without updating remote state, use [tfmigrate](https://github.com/minamijoyo/tfmigrate), which allows you to run the `terraform import` command in a declarative way.

Generate a migration file for importing the new resource from a Terraform plan:

```
$ terraform plan -out=tmp.tfplan
$ terraform show -json tmp.tfplan | tfedit migration fromplan -o=tfmigrate_fromplan.hcl
$ cat tfmigrate_fromplan.hcl
migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_acl.example tfedit-test,private",
  ]
}
```

Run the `tfmigrate plan` command to check to see if the `terraform plan` command has no changes after the migration without updating remote state:

```
# tfmigrate plan tfmigrate_fromplan.hcl
(snip.)
YYYY/MM/DD hh:mm:ss [INFO] [migrator] state migrator plan success!
# echo $?
0
```

If looks good, apply it:

```
# tfmigrate apply tfmigrate_fromplan.hcl
(snip.)
YYYY/MM/DD hh:mm:ss [INFO] [migrator] state migrator apply success!
# echo $?
0
```

This is a brief overview of what tfedit is, but an executable example is described later.

If you are not ready for the upgrade, you can pin version constraints in your Terraform configurations with [tfupdate](https://github.com/minamijoyo/tfupdate).

### Implementation status:
For upgrading AWS provider v4, some rules have not been implemented yet. The current implementation status is as follows:

[S3 Bucket Refactor](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#s3-bucket-refactor)

- [x] Arguments of aws_s3_bucket resource
  - [x] acceleration_status
  - [x] acl
  - [x] cors_rule
  - [x] grant
  - [x] lifecycle_rule
  - [x] logging
  - [x] object_lock_configuration rule
  - [x] policy
  - [x] replication_configuration
  - [x] request_payer
  - [x] server_side_encryption_configuration
  - [x] versioning
  - [x] website
- [ ] Meta arguments of resource
  - [x] provider
  - [x] count
  - [x] for_each
  - [ ] dynamic
- [ ] Rename references in an expression to new resource type
- [x] Generate import commands for new split resources

[New Provider Arguments](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/guides/version-4-upgrade#new-provider-arguments)

- [ ] Arguments of provider aws
  - [x] s3_force_path_style

### Known limitations:
- Some arguments were changed not only their names but also valid values. In this case, if a value of the argument is a variable, not literal, it's impossible to automatically rewrite the value of the variable. It potentially could be passed from outside of module or even overwritten at runtime. If it's not literal, you need to change the value of the variable by yourself. The following arguments have this limitation:
  - grant:
    - permissions: A [permissions](https://registry.terraform.io/providers/hashicorp/aws/3.74.3/docs/resources/s3_bucket#permissions) attribute of grant block was a list in v3, but in v4 we need to set each [permission](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl#permission) to each grant block respectively. If the `permissions` attribute is passed as a variable or generated by a function, it cannot be split automatically.
  - lifecycle_rule:
    - enabled = true => status = "Enabled"
    - enabled = false => status = "Disabled"
    - transition:
      - date = "2022-12-31" => date = "2022-12-31T00:00:00Z"
    - expiration:
      - date = "2022-12-31" => date = "2022-12-31T00:00:00Z"
  - object_lock_configuration:
    - object_lock_configuration.object_lock_enabled = "Enabled" => object_lock_enabled = true
  - versioning:
    - enabled = true => status = "Enabled"
    - enabled = false => It also depends on the current status of your bucket. Set `status = "Suspended"` or use `for_each` to avoid creating `aws_s3_bucket_versioning` resource.
    - mfa_delete = true => mfa_delete = "Enabled"
    - mfa_delete = false => mfa_delete = "Disabled"
- Some arguments cannot be converted correctly without knowing the current state of AWS resources. The tfedit never calls the AWS API on your behalf. You have to check it by yourself. The following arguments have this limitation:
  - grant:
    - owner: A [grant](https://registry.terraform.io/providers/hashicorp/aws/3.74.3/docs/resources/s3_bucket#grant) argument of aws_s3_bucket in v3 doesn’t have an owner block, but an access_control_policy argument of aws_s3_bucket_acl in v4 has an [owner](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_acl#access_control_policy) block as required. There is no way to set it automatically without the AWS API call. You need to set the owner by yourself. You can get your AWS canonical user id with [`aws s3api get-bucket-acl --bucket <bucketname>`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/get-bucket-acl.html) or use [aws_canonical_user_id](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/canonical_user_id) data source.

  - lifecycle_rule:
    - filter: When [`aws s3api get-bucket-lifecycle-configuration --bucket <bucketname>`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/get-bucket-lifecycle-configuration.html) returns `"Filter" : {}` without a prefix, you need to set rule.filter as `filter {}`.
    - rule.id: An argument of [`lifecycle_rule.id`](https://registry.terraform.io/providers/hashicorp/aws/3.74.3/docs/resources/s3_bucket#lifecycle_rule) of `aws_s3_bucket` in v3 is optional and computed, but an argument of [`rule.id`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_lifecycle_configuration#rule) of `aws_s3_bucket_lifecycle_configuration` in v4 is required. If the `id` is omitted in the configuration, there is no way to set it automatically without the AWS API call. You need to set it by yourself. You can get the rule id with [`aws s3api get-bucket-lifecycle-configuration --bucket <bucketname>`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/get-bucket-lifecycle-configuration.html).
  - versioning:
    - enabled: Starting from v3.70.0, `enabled = false` for a new bucket doesn't set "Suspended" explicitly. When [`aws s3api get-bucket-versioning --bucket <bucketname>`](https://awscli.amazonaws.com/v2/documentation/api/latest/reference/s3api/get-bucket-versioning.html) returns no `"Status"`, which means `"Disabled"`. In this case, you need to remove the `aws_s3_bucket_versioning` resource.

### Example

We recommend you to play an example in a sandbox environment first, which is safe to run `terraform` and `tfmigrate` command without any credentials. The sandbox environment mocks the AWS API with [localstack](https://github.com/localstack/localstack) and doesn't actually create any resources. So you can safely and easily understand how it works.

Build a sandbox environment with docker-compose and run bash:

```
$ git clone https://github.com/minamijoyo/tfedit
$ cd tfedit/
$ docker-compose build
$ docker-compose run --rm tfedit /bin/bash
```

In the sandbox environment, create and initialize a working directory from test fixtures:

```
# mkdir -p tmp/dir1 && cd tmp/dir1
# terraform init -from-module=../../test-fixtures/awsv4upgrade/aws_s3_bucket/simple/
# cat main.tf
```

This example contains a simple `aws_s3_bucket` resource:

```main.tf
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
  acl    = "private"
}
```

Apply it and create the `aws_s3_bucket` resource with the AWS provider v3.74.3, which is the last version without deprecation warnings:

```
# terraform -v
Terraform v1.1.8
on linux_amd64
+ provider registry.terraform.io/hashicorp/aws v3.74.3

# terraform apply -auto-approve
# terraform state list
aws_s3_bucket.example
```

Then, let's upgrade the AWS provider to the latest v4.x. We recommend upgrading to v4.9.0 or later because before 4.9.0 includes some breaking changes. To update the provider version constraint, of course you can edit the `required_providers` block in the `config.tf` with your text editor, but it's easy to do with [tfupdate](https://github.com/minamijoyo/tfupdate):

```
# tfupdate provider aws -v "~> 4.9" .
# terraform init -upgrade
# terraform -v
Terraform v1.1.8
on linux_amd64
+ provider registry.terraform.io/hashicorp/aws v4.9.0
```

You can see a deprecation warning as follows:

```
# terraform validate
╷
│ Warning: Argument is deprecated
│
│   with aws_s3_bucket.example,
│   on main.tf line 3, in resource "aws_s3_bucket" "example":
│    3:   acl    = "private"
│
│ Use the aws_s3_bucket_acl resource instead
╵
Success! The configuration is valid, but there were some validation warnings as shown above.
```

Now, it's time to upgrade Terraform configuration to the AWS provider v4 compatible with `tfedit`:

```
# tfedit filter awsv4upgrade -u -f main.tf
# cat main.tf
```

You can see the `acl` argument has been split into an `aws_s3_bucket_acl` resource:

```
resource "aws_s3_bucket" "example" {
  bucket = "tfedit-test"
}

resource "aws_s3_bucket_acl" "example" {
  bucket = aws_s3_bucket.example.id
  acl    = "private"
}
```

You can also see that the deprecation warning has been resolved:

```
# terraform validate
Success! The configuration is valid.
```

At this point, if you run the `terraform plan` command, you can see that a new `aws_s3_bucket_acl` resource will be created:

```
# terraform plan
(snip.)
Plan: 1 to add, 0 to change, 0 to destroy.
```

To resolve the conflict between the configuration and the existing state, you need to import the new resource. As you know, you can run the `terraform import` command directly, but if you prefer to check the upgrade results without updating remote state, use [tfmigrate](https://github.com/minamijoyo/tfmigrate), which allows you to run the `terraform import` command in a declarative way.

Generate a migration file for importing the new resource from a Terraform plan:

```
$ terraform plan -out=tmp.tfplan
$ terraform show -json tmp.tfplan | tfedit migration fromplan -o=tfmigrate_fromplan.hcl
$ cat tfmigrate_fromplan.hcl
migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_acl.example tfedit-test,private",
  ]
}
```

Run the `tfmigrate plan` command to check to see if the `terraform plan` command has no changes after the migration without updating remote state:

```
# tfmigrate plan tfmigrate_fromplan.hcl
(snip.)
YYYY/MM/DD hh:mm:ss [INFO] [migrator] state migrator plan success!
# echo $?
0
```

If looks good, apply it:

```
# tfmigrate apply tfmigrate_fromplan.hcl
(snip.)
YYYY/MM/DD hh:mm:ss [INFO] [migrator] state migrator apply success!
# echo $?
0
```

The `tfmigrate apply` command computes a new state and pushes it to remote state.
It will fail if the `terraform plan` command detects any diffs with the new state.

Finally, You can confirm the latest remote state has no changes with the `terraform plan` command in v4:

```
# terraform plan
(snip.)
No changes. Infrastructure is up-to-date.

# terraform state list
aws_s3_bucket.example
aws_s3_bucket_acl.example
```

To clean up the sandbox environment:

```
# terraform destroy -auto-approve
# cd ../../
# rm -rf tmp/dir1
# exit
$ docker-compose down
```

Tips: If you see something was wrong, you can run the `awslocal` command, which is configured to call AWS APIs to the localstack endpoint:

```
$ docker exec -it tfedit_localstack_1 awslocal s3api list-buckets
```

## Install

### Homebrew

If you are macOS user:

```
$ brew install minamijoyo/tfedit/tfedit
```

### Download

Download the latest compiled binaries and put it anywhere in your executable path.

https://github.com/minamijoyo/tfedit/releases

### Source

If you have Go 1.18+ development environment:

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
  migration   Generate a migration file for state operations
  version     Print version

Flags:
  -h, --help   help for tfedit

Use "tfedit [command] --help" for more information about a command.
```

```
$ tfedit filter --help
Apply a built-in filter

Usage:
  tfedit filter [flags]
  tfedit filter [command]

Available Commands:
  awsv4upgrade Apply a built-in filter for awsv4upgrade

Flags:
  -h, --help   help for filter

Global Flags:
  -f, --file string   A path of input file (default "-")
  -u, --update        Update files in-place

Use "tfedit filter [command] --help" for more information about a command.
```

```
$ tfedit filter awsv4upgrade --help
Apply a built-in filter for awsv4upgrade

Upgrade configurations to AWS provider v4.

Usage:
  tfedit filter awsv4upgrade [flags]

Flags:
  -h, --help   help for awsv4upgrade

Global Flags:
  -f, --file string   A path of input file (default "-")
  -u, --update        Update files in-place
```

By default, the input is read from stdin, and the output is written to stdout.
You can also read a file with `-f` flag, and update the file in-place with `-u` flag.

```
$ tfedit migration --help
Generate a migration file for state operations

Usage:
  tfedit migration [flags]
  tfedit migration [command]

Available Commands:
  fromplan    Generate a migration file from Terraform JSON plan file

Flags:
  -h, --help   help for migration

Use "tfedit migration [command] --help" for more information about a command.
```

```
$ tfedit migration fromplan --help
Generate a migration file from Terraform JSON plan file

Read a Terraform plan file in JSON format and
generate a migration file in tfmigrate HCL format.
Currently, only import actions required by awsv4upgrade are supported.

Usage:
  tfedit migration fromplan [flags]

Flags:
  -d, --dir string    Set a dir attribute in a migration file
  -f, --file string   A path to input Terraform JSON plan file (default "-")
  -h, --help          help for fromplan
  -o, --out string    Write a migration file to a given path (default "-")
```

By default, the input is read from stdin, and the output is written to stdout.
You can also read a file with `-f` flag, and write a file with `-o` flag.

## License

MIT
