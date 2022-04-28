package migration

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateFromPlan(t *testing.T) {
	cases := []struct {
		desc     string
		planFile string
		dir      string
		ok       bool
		want     string
	}{
		{
			desc:     "import simple",
			planFile: "test-fixtures/import_simple.tfplan.json",
			dir:      "",
			ok:       true,
			want: `migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_acl.example tfedit-test,private",
  ]
}
`,
		},
		{
			desc:     "import simple with dir",
			planFile: "test-fixtures/import_simple.tfplan.json",
			dir:      "foo",
			ok:       true,
			want: `migration "state" "fromplan" {
  dir = "foo"
  actions = [
    "import aws_s3_bucket_acl.example tfedit-test,private",
  ]
}
`,
		},
		{
			desc:     "import full",
			planFile: "test-fixtures/import_full.tfplan.json",
			dir:      "",
			ok:       true,
			want: `migration "state" "fromplan" {
  actions = [
    "import aws_s3_bucket_accelerate_configuration.example tfedit-test",
    "import aws_s3_bucket_acl.example tfedit-test,private",
    "import aws_s3_bucket_acl.log tfedit-log,log-delivery-write",
    "import aws_s3_bucket_cors_configuration.example tfedit-test",
    "import aws_s3_bucket_lifecycle_configuration.example tfedit-test",
    "import aws_s3_bucket_logging.example tfedit-test",
    "import aws_s3_bucket_object_lock_configuration.example tfedit-test",
    "import aws_s3_bucket_policy.example tfedit-test",
    "import aws_s3_bucket_replication_configuration.example tfedit-test",
    "import aws_s3_bucket_request_payment_configuration.example tfedit-test",
    "import aws_s3_bucket_server_side_encryption_configuration.example tfedit-test",
    "import aws_s3_bucket_versioning.example tfedit-test",
    "import aws_s3_bucket_website_configuration.example tfedit-test",
  ]
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			planJSON, err := os.ReadFile(tc.planFile)
			if err != nil {
				t.Fatalf("failed to read file: %s", err)
			}

			output, err := GenerateFromPlan(planJSON, tc.dir)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, got: %s", got)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Fatalf("got:\n%s\nwant:\n%s\ndiff:\n%s", got, tc.want, diff)
			}
		})
	}
}
