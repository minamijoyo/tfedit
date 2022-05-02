#!/bin/bash

set -eo pipefail

usage()
{
  cat << EOF
  Usage: `basename $0` <command> <fixture>

  Arguments:
    command: A name of step tu run. Valid values are:
             run | setup | upgrade | filter | generate | migrate | cleanup
    fixture: A name of fixture in test-fixtures/awsv4upgrade/aws_s3_bucket/
EOF
}

setup()
{
  terraform init -input=false -no-color -from-module="$FIXTUREDIR"
  terraform -v

  terraform apply -input=false -no-color -auto-approve
  terraform state list
}

upgrade()
{
  tfupdate provider aws -v "~> 4.9" .
  terraform init -input=false -no-color -upgrade
  terraform -v

  # fix path style for sandbox only
  hcledit attribute rm -u -f config.tf provider.aws.s3_force_path_style
  hcledit attribute append -u -f config.tf provider.aws.s3_use_path_style true
}

filter()
{
  terraform validate -json -no-color
  before_count=$(terraform validate -json -no-color | jq '[.error_count, .warning_count] | add')
  if [[ $before_count -eq 0 ]]; then
    echo "expected to an error before filter"
    exit 1
  fi

  cat main.tf

  find . -type f -name '*.tf' -print0 | xargs -0 -I {} tfedit filter awsv4upgrade -u -f {}

  cat main.tf

  terraform validate -json -no-color
  after_count=$(terraform validate -json -no-color | jq '[.error_count, .warning_count] | add')
  if [[ $after_count -ne 0 ]]; then
    echo "expected to no error after filter"
    exit 1
  fi
}

generate()
{
  terraform plan -input=false -no-color -out=tmp.tfplan
  terraform show -json tmp.tfplan | tfedit migration fromplan -o=tfmigrate_fromplan.hcl
  cat tfmigrate_fromplan.hcl
  diff -u tfmigrate_want.hcl tfmigrate_fromplan.hcl
  rm -f tmp.tfplan
}

migrate()
{
  tfmigrate apply tfmigrate_fromplan.hcl
  terraform plan -input=false -no-color -detailed-exitcode
  terraform state list
}

cleanup()
{
  terraform destroy -input=false -no-color -auto-approve
  find ./ -mindepth 1 -delete
}

run()
{
  setup
  upgrade
  filter
  generate
  migrate
  cleanup
}

# main
if [[ $# -ne 2 ]]; then
  usage
  exit 1
fi

set -x

COMMAND=$1
FIXTURE=$2

REPO_ROOT_DIR="$(git rev-parse --show-toplevel)"
WORKDIR="$REPO_ROOT_DIR/tmp/testacc/awsv4upgrade/aws_s3_bucket/$FIXTURE"
FIXTUREDIR="$REPO_ROOT_DIR/test-fixtures/awsv4upgrade/aws_s3_bucket/$FIXTURE/"
mkdir -p "$WORKDIR"
pushd "$WORKDIR"

case "$COMMAND" in
  run | setup | upgrade | filter | generate | migrate | cleanup )
    "$COMMAND"
    RET=$?
    ;;
  *)
    usage
    RET=1
    ;;
esac

popd
exit $RET
