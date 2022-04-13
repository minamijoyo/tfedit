#!/bin/bash

set -eo pipefail

script_full_path=$(dirname "$0")

# test simple
bash "$script_full_path/awsv4upgrade.sh" run simple

# test all
repo_root_dir="$(git rev-parse --show-toplevel)"
fixturesdir="$repo_root_dir/test-fixtures/awsv4upgrade/aws_s3_bucket/"

# Exclude simple because it has been already tested above.
# Exclude grant because owner id cannot be set automatically.
fixtures=$(
  find $fixturesdir -type d -mindepth 1 -maxdepth 1 -exec basename {} \; | sort \
  | grep -v -e '^simple$' \
  | grep -v -e '^grant$' \
)

for fixture in ${fixtures}
do
  echo $fixture
  bash "$script_full_path/awsv4upgrade.sh" run $fixture
done
